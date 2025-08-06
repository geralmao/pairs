//go:build !js && !wasm

package images

import "embed"

//go:embed emojis/*.png
var EmojisDataFS embed.FS

const EmojisDirName string = "emojis"
