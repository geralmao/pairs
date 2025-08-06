//go:build js && wasm

package images

import "embed"

//go:embed emojis_lite/*.png
var EmojisDataFS embed.FS

const EmojisDirName string = "emojis_lite"
