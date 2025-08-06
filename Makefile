VERSION := 1.1.0

DIST_DIR := bin
WEB_DIR := ${DIST_DIR}/web
WEB_WASM := $(WEB_DIR)/matchemojis.wasm
WEB_WASM_TMP := $(WEB_DIR)/matchemojis_tmp.wasm
DESKTOP_DIR:=${DIST_DIR}/desktop
MODULE := github.com/programatta/pairs

.PHONY: build build-web run run-web clean

build:
	mkdir -p ${DESKTOP_DIR}
	go build -ldflags "-X '$(MODULE)/internal.Version=$(VERSION)'" -o ${DESKTOP_DIR}/matchemojis main.go

build-win:
	mkdir -p ${DESKTOP_DIR}
	env GOOS=windows GOARCH=amd64 go build -ldflags "-X '$(MODULE)/internal.Version=$(VERSION)'" -o ${DESKTOP_DIR}/matchemojis.exe main.go

# Requiere de un OSX para realizar compilación nativa con bindings de C
# build-mac:
# 	env GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o ${DIST_DIR}/matchemojis-mac main.go

# build-mac-arm:
# 	env GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.Version=$(VERSION)" -o ${DIST_DIR}/matchemojis-macarm main.go

build-web:
	mkdir -p ${WEB_DIR}
	env GOOS=js GOARCH=wasm go build -ldflags="-s -w -X '$(MODULE)/internal.Version=$(VERSION)'" -buildvcs=false -o ${WEB_WASM_TMP} ${MODULE}
	wasm-opt -Oz --enable-bulk-memory --strip-debug --strip-dwarf --strip-producers ${WEB_WASM_TMP} -o ${WEB_WASM}
	rm ${WEB_WASM_TMP}
	cp $$(go env GOROOT)/lib/wasm/wasm_exec.js ${WEB_DIR}
	cp bin/wasm-template/helplib.js ${WEB_DIR}
	printf '%s\n' \
	'<!DOCTYPE html>' \
	'<html>' \
	'  <head>' \
	'    <meta charset="UTF-8">' \
	'    <title>Match Emojis - Ebiten</title>' \
	'  </head>' \
	'  <body>' \
	'    <script src="wasm_exec.js"></script>' \
	'    <script src="helplib.js"></script>' \
	'    <script>' \
	'      const go = new Go();' \
	'      polyfillIinstantiateStreamingSupport();' \
	'      WebAssembly.instantiateStreaming(fetch("matchemojis.wasm"), go.importObject).then(result => {' \
	'        go.run(result.instance);' \
	'        activeTouchesToMouseInDevice();' \
	'      });' \
	'    </script>' \
	'  </body>' \
	'</html>' \
	> ${WEB_DIR}/index.html

build-all: build build-win build-web

run:
	go run main.go

run-web:
	go run github.com/hajimehoshi/wasmserve@latest .

clean:
	rm -rf ${DESKTOP_DIR}
	rm -rf ${WEB_DIR}
