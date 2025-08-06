package views

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/programatta/pairs/internal/config"
)

type ViewId int

const (
	Loader ViewId = iota
	Menu
	Play
	Settings
)

// View define el comportamiento de un estado de juego
type Viewer interface {
	Start(context *config.GameContext)
	ProcessEvents()
	Update()
	Draw(screen *ebiten.Image)
	NextView() ViewId
}
