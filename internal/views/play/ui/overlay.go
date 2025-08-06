package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Overlay interface {
	CanRemove() bool
	Update()
	Draw(screen *ebiten.Image, textFace *text.GoTextFace)
}
