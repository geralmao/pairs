# 🇪🇸 PAIRS.
Otro juego de cartas donde se tienen que encontrar las parejas, trios, cuartetos, etc., según se avance en el juego.

## Objetivo.
La idea de este mini juego es abordar la exportación de un juego realizado con **Golang** y **Ebitengine** a la plataforma **Android**, además de las ya experimentadas Web y Desktop.

## Desarrollo.
El juego surge de la construcción de una pequeña demo llamada [demoandroid](https://github.com/programatta/demoandroid) creada para estudiar:
* La estructuración del proyecto.
* Generación de librerías Golang `.aar` para su uso en Android.
* Creación de APKs y AAB tanto de debug como release sin necesidad de Android Studio. 

Actualmente se ha cubierto siguientes las plataformas:
* ✅ Binarios para Windows / Linux.
* ✅ Web con WASM.
* ✅ Android (APKs).

 ### Problematica detectada.
 El uso de wasm en desktop funciona de forma correcta, pero en los navegadores de los dispositivos móviles presentan varios problemas:
 * En mi dispositivo Android, directamente se queda la pantalla en negro.
 * Un usuario informa que en su Android carga el juego, pero el menú no responde a eventos.
 * En un dispositivo iOS al ejecutar el juego, la página se recarga.

#### Tamaño del fichero wasm.
Investigando, parece que los navegadores de dispositivos móviles tienen problemas con archivos **wasm** grandes.
En este caso:
* El binario orignial pesaba 39MB.
* Tras las optimizaciones se reduce a 18.5MB aproximadamente.

Los assets eran bastante grandes, por ejemplo la música de fondo que era en formato **wav** y ocupaba **18MB** y se transformó a formato **ogg** quedando en **1.2MB**. 

También para la versión WASM los assets de los emojis es un conjunto más pequeño (242KB) que en la versión Android y Desktop (8.2MB).

> 🔔 Nota. 
> No he encontrado información sobre el tamaño máximo aceptado.

También para reducir el tamaño del fichero **wasm** usamos los flags **-s** y **-w** que permite usar el compilador sin símbolos de depuración y con optimizaciones:

~~~bash
env GOOS=js GOARCH=wasm go build -ldflags="-s -w -X '$(MODULE)/internal.Version=$(VERSION)'" -buildvcs=false -o ${WEB_WASM_TMP} ${MODULE}
~~~

 Aparte de aplicar estos cambios para **wasm** también se ha usado la herramienta **wasm-opt** que forma parte del paquete [Binaryen](https://github.com/WebAssembly/binaryen), que es un conjunto de herramientas para optimizar y trabajar con WebAssembly. Yo lo he instalado en mi imagen de docker con `apt-get install binaryen` y lo he usado con los siguientes opciones:

 * **-Oz**: aplicamos la máxima optimización.
 * **--enable-bulk-memory**: si nos aparece el error `[wasm-validator error in function 364] unexpected false: Bulk memory operation (bulk memory is disabled), on
(memory.copy ...)`. 

> 🔔 Nota.
> **Bulk Memory Operations** son una serie de instrucciones de WebAssembly que permiten copiar, mover y inicializar grandes bloques de memoria de forma más eficiente.

* **--strip-debug**: permite eliminar la tabla de símbolos y otros datos de depuración.
* **--strip-producers**: permite eliminar la sección Producers, que contiene información sobre el compilador y las herramientas que se usaron.
* **--strip-dwarf**: permite eliminar específicamente la información de depuración DWARF.
 
 ~~~bash
 wasm-opt -Oz --enable-bulk-memory --strip-debug --strip-dwarf --strip-producers ${WEB_WASM_TMP} -o ${WEB_WASM}
 ~~~

Con este tamaño, en mi dispositivo android y en iOS ya permite cargar el juego desde el navegador y llegar al menú de este pero sin poder interactuar con las opciones.

#### Los eventos touch.
Los eventos touch en dispositivos móviles no parecen propagarse correctamente desde javascript a la capa de Golang/Ebitengine. 

Como solución temporal, se añade funcionalidad Javascript que convierten los eventos **touch** en eventos **click**, lo que permite que la aplicación funcione. 

Para ello: 
* Se crea un archivo **helplib.js**.
* Se exporta al crear el paquete (sin comprimir) web a través del **Makefile**.
* Se usa la función **activeTouchesToMouseInDevice();** en el html generado.


# 🇬🇧 PAIRS.
Another card game where you have to find pairs, trios, quartets, and so on, as you progress through the game.

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
