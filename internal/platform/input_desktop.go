//go:build !android

package platform

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func IsPressEventJustRelease() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}

func IsPressEventPressed() bool {
	return ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}

func PressPosition() (int, int) {
	return ebiten.CursorPosition()
}
