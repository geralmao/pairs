package images

import (
	_ "embed"
	"fmt"
)

//go:embed loader.jpg
var LoaderData []byte

//go:embed hourglass.png
var HourGlassData []byte

func GetEmojisNamesFromFS() []string {
	var emojisNames []string

	dirEntries, errDirEntries := EmojisDataFS.ReadDir(EmojisDirName)
	if errDirEntries != nil {
		panic(errDirEntries)
	}

	for _, dirEntry := range dirEntries {
		emojisNames = append(emojisNames, fmt.Sprintf("%s/%s", EmojisDirName, dirEntry.Name()))
	}

	return emojisNames
}
