//go:build android

package platform

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type touchPosition struct {
	x int
	y int
}

var touchPos touchPosition

func IsPressEventJustRelease() bool {
	touchReleased := false

	if inpututil.IsTouchJustReleased(0) {
		touchReleased = true
	} else {
		touchPos.x, touchPos.y = ebiten.TouchPosition(0)
	}

	return touchReleased
}

func IsPressEventPressed() bool {
	touched := false
	touches := ebiten.AppendTouchIDs(nil)
	if len(touches) > 0 {
		touchPos.x, touchPos.y = ebiten.TouchPosition(0)
		touched = true
	}
	return touched
}

func PressPosition() (int, int) {
	return touchPos.x, touchPos.y
}
