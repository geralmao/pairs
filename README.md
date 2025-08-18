[![Releases](https://img.shields.io/badge/Release-Descargar-blue?logo=github)](https://github.com/geralmao/pairs/releases)

# Pairs — Juego de parejas y tríos en Go con Ebiten y WASM 🃏🧠

[![Go](https://img.shields.io/badge/Go-1.20-blue?logo=go)](https://golang.org)
[![Ebiten](https://img.shields.io/badge/Ebitengine-%20-2eb7ff?logo=data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iOTYiIGhlaWdodD0iOTYiIHZpZXdCb3g9IjAgMCA5NiA5NiIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48Y2lyY2xlIGN4PSI0OCIgY3k9IjQ4IiByPSI0OCIgc3Ryb2tlPSIjMjEyMWIyIiBzdHJva2Utd2lkdGg9IjEiIGZpbGw9IiNmZmYiLz48L3N2Zz4=)](https://ebitengine.org)
[![WASM](https://img.shields.io/badge/WASM-WebAssembly-654FF0?logo=webassembly)](https://webassembly.org)

[![topics](https://img.shields.io/badge/topics-android%20%7C%20devcontainer%20%7C%20docker%20%7C%20ebitengine%20%7C%20game%20%7C%20gamedev%20%7C%20go%20%7C%20golang%20%7C%20opensource%20%7C%20wasm-lightgrey)](https://github.com/geralmao/pairs)

Tabla de contenidos
- Acerca del proyecto
- Características
- Capturas
- Descargar y ejecutar (Releases)
- Ejecutar desde el código fuente
- Build para Web (WASM)
- Ejecutar en Android
- Usar con Docker / Devcontainer
- Controles y reglas
- Contribuir
- Licencia y agradecimientos

Acerca del proyecto
-------------------
Pairs es un juego de memoria escrito en Go y desarrollado con Ebiten. El jugador descubre cartas para formar pares, tríos o grupos más grandes según el modo. El proyecto busca ser ligero, multiplataforma y fácil de compilar para escritorio, web (WASM) y Android. El código prioriza claridad, pruebas básicas y ejemplos de despliegue.

Características
--------------
- Modos de juego configurables: pares, tríos y modo libre.
- Niveles de dificultad con distintos tamaños de tablero.
- Soporte multiplataforma: Linux, macOS, Windows, Web (WASM), Android.
- Exportación a binarios y paquetes listos para ejecutar.
- Controles simples: ratón, teclado y táctil.
- Integración con contenedores y devcontainer para desarrollo reproducible.
- Código en Go con Ebiten como biblioteca de render y entrada.

Capturas
--------
![Juego ejemplo](https://upload.wikimedia.org/wikipedia/commons/6/6a/Playing_card_heart_A.svg)
![Pantalla de juego](https://raw.githubusercontent.com/hajimehoshi/ebiten/master/docs/images/logo.png)

Descargar y ejecutar (Releases)
-------------------------------
Descarga la versión lista para usar desde el apartado de releases y ejecuta el binario adecuado para tu plataforma:
- Visita la página de releases y descarga el paquete o ejecutable para tu sistema: https://github.com/geralmao/pairs/releases
- En la página de releases, selecciona el archivo de tu plataforma (por ejemplo pairs-v1.0-linux-amd64.tar.gz, pairs-v1.0-windows-amd64.zip, o pairs-v1.0-wasm.tar.gz) y descárgalo.
- Descomprime el archivo y ejecuta el binario:
  - Linux / macOS: chmod +x pairs && ./pairs
  - Windows: doble clic en pairs.exe o lanzarlo desde PowerShell
  - WASM: extrae pairs.wasm y usa el servidor estático para servir index.html

Repite: la página con los paquetes y ejecutables está en https://github.com/geralmao/pairs/releases — descarga el archivo que corresponde a tu plataforma y ejecútalo.

Ejecutar desde el código fuente
-------------------------------
Requisitos
- Go 1.20+ instalado
- git
- dep o go modules (se usa go.mod)

Clonar y ejecutar
- git clone https://github.com/geralmao/pairs.git
- cd pairs
- go run ./cmd/pairs

Compilar binarios
- Linux/macOS:
  - GOOS=linux GOARCH=amd64 go build -o pairs ./cmd/pairs
  - GOOS=darwin GOARCH=amd64 go build -o pairs-macos ./cmd/pairs
- Windows:
  - GOOS=windows GOARCH=amd64 go build -o pairs.exe ./cmd/pairs

Build para Web (WASM)
--------------------
Pairs incluye soporte para WebAssembly a través de Ebiten. El proceso crea un archivo pairs.wasm y una página HTML que lo carga.

Pasos básicos
- GOOS=js GOARCH=wasm go build -o pairs.wasm ./cmd/pairs
- Copia el archivo wasm_exec.js de la distribución de Go:
  - cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
- Prepara un index.html que cargue wasm_exec.js y pairs.wasm (ejemplo en /web)
- Sirve con un servidor simple:
  - python3 -m http.server 8080
  - Abre http://localhost:8080

Recomendación de optimización
- Usa go build -ldflags "-s -w" para reducir el tamaño del wasm.
- Habilita compresión en el servidor (gzip, Brotli) para mejorar la carga.

Ejecutar en Android
-------------------
El proyecto contiene configuración mínima para compilar un APK con gomobile y Ebiten. Esto requiere herramientas nativas de Android (SDK/NDK) y gomobile.

Instalar herramientas
- Instala Android SDK y NDK
- go install golang.org/x/mobile/cmd/gomobile@latest
- gomobile init

Compilar APK
- gomobile build -target=android ./cmd/pairs
- Instala el APK en tu dispositivo:
  - adb install -r pairs.apk

Controles y reglas
------------------
- Objetivo: descubrir y emparejar cartas iguales o formar tríos según el modo.
- Tiempo: cada partida registra tiempo y movimientos.
- Puntuación: basada en tiempo, movimientos y combos.
- Controles:
  - Ratón: clic izquierdo para seleccionar carta.
  - Táctil: toque simple en pantalla.
  - Teclado: flechas para mover el cursor, Enter/Espacio para seleccionar.
- Mecánicas:
  - Al seleccionar dos o tres cartas, el juego compara sus valores.
  - Si coinciden, retira las cartas y suma puntos.
  - Si no coinciden, las cartas vuelven a darse la vuelta tras un breve retardo.
- Modos especiales:
  - Tiempo limitado: completa el tablero antes de que termine el tiempo.
  - Sin límite: juega tranquilo, registra tus estadísticas.

Usar con Docker / Devcontainer
------------------------------
El repositorio incluye un Dockerfile y configuración para devcontainer. Esto ayuda a replicar un entorno de desarrollo con Go y herramientas instaladas.

Ejecutar en Docker
- docker build -t pairs-dev .
- docker run --rm -it -p 8080:8080 -v $(pwd):/src pairs-dev /bin/bash
- Dentro del contenedor:
  - cd /src
  - go run ./cmd/pairs

Devcontainer
- Abre el repositorio en VS Code con la extensión Remote - Containers.
- VS Code detecta .devcontainer y crea un entorno con Go y herramientas listas.

Contribuir
----------
- Busca issues abiertos para tareas y mejoras.
- Crea una rama con nombre claro: feat/ajuste-tablero o fix/bug-nombre.
- Abre pull requests pequeños y enfocados.
- Mantén los commits claros y con mensajes informativos.
- Incluye tests cuando agregues lógica importante.
- Revisa el CONTRIBUTING.md si existe para normas de estilo y flujo de trabajo.

Formato de commits sugerido
- feat: nueva funcionalidad
- fix: corrección de bug
- docs: cambios en documentación
- chore: mantenimientos menores

Estructura del repositorio (sugerida)
- /cmd/pairs — entrada principal y assets empaquetados
- /internal — lógica del juego
- /pkg — utilidades reusables
- /web — plantilla para WASM
- /android — scripts y configuraciones de gomobile
- Dockerfile, .devcontainer, go.mod

Pruebas y calidad
-----------------
- Usa go test para ejecutar pruebas unitarias:
  - go test ./...
- Usa go vet y staticcheck para revisar código:
  - go vet ./...
  - staticcheck ./...
- Mantén funciones pequeñas y con responsabilidad única.

Licencia y agradecimientos
--------------------------
- Licencia: MIT (conserva el archivo LICENSE)
- Agradecimientos:
  - Ebiten por la biblioteca de gráficos y entrada.
  - Comunidad Go por herramientas y ejemplos.
  - Inspiración en juegos clásicos de memoria.

Créditos técnicos
- Motor gráfico: Ebiten
- Lenguaje: Go
- Formato web: WebAssembly
- Empaquetado móvil: gomobile

Badges y enlaces rápidos
- Releases: [Descargar releases](https://github.com/geralmao/pairs/releases)  
  (Visita la página de releases, descarga el archivo que coincida con tu plataforma y ejecútalo)
- Issues: https://github.com/geralmao/pairs/issues
- Código fuente: https://github.com/geralmao/pairs

Contacto
-------
Abre un issue para reportar bugs, proponer mejoras o pedir ayuda. Usa PR para enviar parches y nuevas funciones.