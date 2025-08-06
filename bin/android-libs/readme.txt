*** NO SE AÑADEN LOS FICHEROS POR SER BINARIOS PESADOS ***

En este directorio deberá aparecer los siguientes archivos:
- pairlib-sources.jar
- pairlib.aar

Estos ficheros se generan al usar la imagen docker que contiene el sdk de java y android:
docker run --rm -v "$PWD":/code -w /code go-android-builder ebitenmobile bind -target android -javapkg com.programatta.games.pairs.corelib -o bin/android-libs/pairlib.aar github.com/programatta/pairs/mobile

El fichero pairlib.aar lo copiaremos en bin/android/app/libs


*** FILES ARE NOT ADDED AS THEY ARE HEAVY BINARIES ***

In this directory the following files should be listed:
- pairlib-sources.jar
- pairlib.aar

These files are generated when using the docker image containing the java and android sdk:
docker run --rm -v "$PWD":/code -w /code go-android-builder ebitenmobile bind -target android -javapkg com.programatta.games.pairs.corelib -o bin/android-libs/pairlib.aar github.com/programatta/pairs/mobile

Copy the pairlib.aar file to bin/android/app/libs

Translated with DeepL.com (free version)