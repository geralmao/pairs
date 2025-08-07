# 🇬🇧 MATCH EMOJIS.
Another card game where you have to find pairs, trios, quartets, and so on, as you progress through the game.

## Sources and licenses of the assets used.
* **Emojis**: Downloaded from [OpenMoji](https://openmoji.org) (license: [Creative Commons Attribution-ShareAlike 4.0 International (CC BY-SA 4.0)](https://creativecommons.org/licenses/by-sa/4.0/)) project designed by [HfG Schwäbisch Gmünd](https://www.hfg-gmuend.de/).
* **Font**: [Luckiest Guy](https://fonts.google.com/specimen/Luckiest+Guy) from Google Fonts (license: [Apache 2.0](https://fonts.google.com/specimen/Luckiest+Guy/license)) designed by [Astigmatic](https://fonts.google.com/?query=Astigmatic). 
* **Sound effects**:
  * [30 Free Game Sound FX Pack](https://www.gamedevmarket.net/asset/30-free-game-sound-fx-pack) by [GameBurp](https://www.gamedevmarket.net/member/gameburp).
  * [50 Free Game Sounds Pack](https://www.gamedevmarket.net/asset/50-free-game-sounds-pack-not-a-placeholder) by [PlaceHolderAssets](https://www.gamedevmarket.net/member/placeholderassets). 
* **Music**: "Fligh Home", included in [“lofi world vol 1”](https://www.gamedevmarket.net/asset/lofi-world-volume-1-7-free-lofi-tracks) composed by [kummel](https://www.gamedevmarket.net/member/kummel).
* **Hourglass sprite**: [game-hourglass-pixelated](https://www.vecteezy.com/png/54978930-game-hourglass-pixelated) created by [Idalba Granada](https://www.vecteezy.com/members/studiogstock). 

## Objective.
The idea behind this mini-game is to explore the process of exporting a game made with **Golang** and **Ebitengine** to the **Android**platform, in addition to the already tested Web and Desktop platforms.

## Development.
The game originated from a small demo project called  [demoandroid](https://github.com/programatta/demoandroid), created to study:
* Project structuring.
* Generation of Golang `.aar` libraries for use in Android.
* Creation of APKs and AABs for both debug and release, without the need for Android Studio.

Currently, the following platforms are supported:
* ✅ Binaries for Windows / Linux.
* ✅ Web with WASM.
* ✅ Android (APKs).

### Identified Issues.
Using WASM works correctly on desktop browsers, but several problems appear on mobile browsers:
* On my Android device, the screen remains black.
* A user reported that the game loads on their Android device, but the menu does not respond to any events.
* On an iOS device, the page reloads when the game is launched.

#### WASM File Size.
After some investigation, it seems mobile browsers struggle with large **wasm** files.

In this case:
* The original binary was **39MB**.
* After optimization, the size was reduced to around **18.5MB**.

The assets were quite large, for example, the background music in **wav** format was **18MB**, but after converting to **ogg**, it was reduced to **1.2MB**. 

Also, for the WASM version, the emoji assets were significantly reduced (242KB) compared to the Android and Desktop versions (8.2MB).

> 🔔 Note. 
> I couldn’t find any official information about maximum accepted WASM sizes.

To reduce the size of the **wasm** file, the following Go build flags **-s** y **-w** were used to remove debug symbols and apply optimizations:

~~~bash
env GOOS=js GOARCH=wasm go build -ldflags="-s -w -X '$(MODULE)/internal.Version=$(VERSION)'" -buildvcs=false -o ${WEB_WASM_TMP} ${MODULE}
~~~

 Additionally, the tool **wasm-opt** from the [Binaryen](https://github.com/WebAssembly/binaryen) package was used, a set of tools to optimize and work with WebAssembly. I installed it in my Docker image using `apt-get install binaryen` and used it with the following options:

 * **-Oz**: Apply maximum size optimization.
 * **--enable-bulk-memory**: Needed if you get the error `[wasm-validator error in function 364] unexpected false: Bulk memory operation (bulk memory is disabled), on
(memory.copy ...)`. 

> 🔔 Note.
> **Bulk Memory Operations** are a set of WebAssembly instructions that allow copying, moving, and initializing large memory blocks more efficiently.

* **--strip-debug**: Removes debug symbol tables and related data.
* **--strip-producers**: Removes metadata about the compiler/tools used.
* **--strip-dwarf**: Removes DWARF debug information.
 
 ~~~bash
 wasm-opt -Oz --enable-bulk-memory --strip-debug --strip-dwarf --strip-producers ${WEB_WASM_TMP} -o ${WEB_WASM}
 ~~~

With this reduced size, the game now loads from the browser on both Android and iOS, and reaches the main menu, although it still doesn’t respond to input.

#### Touch Events.
Touch events on mobile devices don’t seem to propagate correctly from JavaScript to the Golang/Ebitengine layer.

A temporary solution was added using JavaScript to convert **touch** events into **click** events, which allows the game to respond. 

To achieve this: 
* A file named **helplib.js** was created.
* It is included in the uncompressed web build via the **Makefile**.
* The function **activeTouchesToMouseInDevice();** is called in the generated HTML.
