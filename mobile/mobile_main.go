//go:build android

package mobile

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/mobile"
	"github.com/programatta/pairs/internal"
	"github.com/programatta/pairs/internal/platform"
)

func init() {
	ebiten.SetFullscreen(true)
	mobile.SetGame(internal.NewGame())
}

// At least one exported function is required by gomobile.
func Dummy() {}

// Funcion que Java llamara para pasar la ruta de la aplicacion para ser
// utilizada por golang y almacenar los datos de juego.
func SetAndroidDataPath(dataPath string) {
	platform.SetUserDataPath(dataPath)
}

// Funcion que Java llamara para pasar el idioma del dispositivo para ser
// utilizado por golang.
func SetAndroidLanguage(language string) {
	platform.SetLanguage(language)
}
