package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/programatta/pairs/internal/config"
	"github.com/programatta/pairs/internal/utils"
)

type PointsFeedback struct {
	points   string
	time     float32
	progress float32
	remove   bool
}

func NewPointsFeedback(points uint) *PointsFeedback {
	spoints := fmt.Sprintf("+%d Pts", points)
	pointsFeedback := &PointsFeedback{
		points: spoints,
	}

	return pointsFeedback
}

//-----------------------------------------------------------------------------
// Implements Overlay Interface
//-----------------------------------------------------------------------------

func (pfb *PointsFeedback) Update() {
	pfb.time += config.Dt
	pfb.progress = pfb.time / feedbackMatchDelay
	if pfb.time > feedbackMatchDelay {
		pfb.remove = true
		pfb.progress = 0
		pfb.time = 0
	}
}

func (pfb *PointsFeedback) Draw(screen *ebiten.Image, textFace *text.GoTextFace) {
	scale := 1.0 + 8.5*utils.EaseOutQuad(float64(pfb.progress))
	screenCenterX := float64(config.WindowWidth) / 2
	screenCenterY := float64(config.WindowHeight) / 2

	widthText, heightText := text.Measure(pfb.points, textFace, 0)
	halfWText := (widthText / 2)
	halfHText := (heightText / 2)

	op := &text.DrawOptions{}
	//! Las transformaciones se hacen en orden inverso a lo que queremos, es decir si queremos:
	//! - traslación al centro, escalar y translación a su mitad.
	op.GeoM.Translate(-halfWText, -halfHText)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(screenCenterX, screenCenterY)

	op.ColorScale.ScaleWithColor(color.White)
	op.ColorScale.ScaleAlpha(1.0 - float32(utils.EaseOutQuad(float64(pfb.progress))))
	text.Draw(screen, pfb.points, textFace, op)
}

func (pfb *PointsFeedback) CanRemove() bool {
	return pfb.remove
}

const feedbackMatchDelay float32 = 1.4
