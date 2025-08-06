package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/programatta/pairs/internal/config"
	"github.com/programatta/pairs/internal/language"
	"github.com/programatta/pairs/internal/utils"
)

type PresentationOverlay struct {
	level             uint
	target            uint
	matchsDescription string
	presentationTime  float32
	remove            bool
}

func NewPresentationOverlay(level, target uint, matchsDescription string) *PresentationOverlay {
	return &PresentationOverlay{
		level:             level,
		target:            target,
		matchsDescription: matchsDescription,
	}
}

//-----------------------------------------------------------------------------
// Implements Overlay Interface
//-----------------------------------------------------------------------------

func (po *PresentationOverlay) Update() {
	po.presentationTime += config.Dt
	if po.presentationTime > presentationDelay {
		po.presentationTime = 0
		po.remove = true
	}
}

func (po *PresentationOverlay) Draw(screen *ebiten.Image, textFace *text.GoTextFace) {
	targWidth := screen.Bounds().Dx()
	targHeight := screen.Bounds().Dy()
	targImg := ebiten.NewImage(targWidth, targHeight)
	targImg.Fill(color.RGBA{0x49, 0x50, 0x57, 205})

	uiLevelText := fmt.Sprintf("%s: %d", language.Value.Level, po.level)
	widthText, _ := text.Measure(uiLevelText, textFace, 0)
	opText := &text.DrawOptions{}
	opText.GeoM.Translate(float64(config.WindowWidth)/2-widthText/2, float64(config.WindowHeight)/2-50)
	opText.ColorScale.ScaleWithColor(color.White)
	text.Draw(targImg, uiLevelText, textFace, opText)

	uiPairsText := fmt.Sprintf("%s: %d %s", language.Value.Goal, po.target, po.matchsDescription)
	widthText, _ = text.Measure(uiPairsText, textFace, 0)
	opText = &text.DrawOptions{}
	opText.GeoM.Translate(float64(config.WindowWidth)/2-widthText/2, float64(config.WindowHeight)/2)
	opText.ColorScale.ScaleWithColor(color.White)
	text.Draw(targImg, uiPairsText, textFace, opText)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	op.ColorScale.ScaleAlpha(1.0 - float32(utils.EaseInQuint(float64(po.presentationTime/presentationDelay))))

	screen.DrawImage(targImg, op)
}

func (po *PresentationOverlay) CanRemove() bool {
	return po.remove
}

const presentationDelay float32 = 2.5
