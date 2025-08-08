package sounds

import (
	"embed"
	"fmt"
)

//go:embed fx/*.ogg
var SoundFxFS embed.FS

//go:embed music/*.ogg
var BackgroundMusicFS embed.FS

func GetSoundFxNamesFromFS() []string {
	var fxsNames []string

	dirEntries, errDirEntries := SoundFxFS.ReadDir("fx")
	if errDirEntries != nil {
		panic(errDirEntries)
	}

	for _, dirEntry := range dirEntries {
		fxsNames = append(fxsNames, fmt.Sprintf("fx/%s", dirEntry.Name()))
	}

	return fxsNames
}
