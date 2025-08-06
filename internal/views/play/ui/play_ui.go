package ui

import (
	"bytes"
	"fmt"
	"image/color"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/programatta/pairs/internal/assets/images"
	"github.com/programatta/pairs/internal/config"
	"github.com/programatta/pairs/internal/language"
	"github.com/programatta/pairs/internal/views/play/common"
)

type PlayUI struct {
	textFace       *text.GoTextFace
	notifier       common.Notifier
	clock          *ebiten.Image
	pointsFeedbaks []Overlay
	overlays       []Overlay
}

func NewPlayUI(textFace *text.GoTextFace, notifier common.Notifier) *PlayUI {
	imgTmp, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(images.HourGlassData))
	if err != nil {
		panic(err)
	}

	return &PlayUI{
		textFace: textFace,
		notifier: notifier,
		clock:    imgTmp,
	}
}

func (pui *PlayUI) AddFeedbackNewPoints(points uint) {
	pui.pointsFeedbaks = append(pui.pointsFeedbaks, NewPointsFeedback(points))
}

func (pui *PlayUI) AddOverlayPresentation(level uint, target int, matchsDescription string) {
	pui.overlays = append(pui.overlays, NewPresentationOverlay(level, uint(target), matchsDescription))
}

func (pui *PlayUI) AddOverlayFinishLevel(timeLeft float32, score uint) {
	pui.overlays = append(pui.overlays, NewFinishLevelOverlay(timeLeft, score, pui.notifier))
}

func (pui *PlayUI) DrawHeader(screen *ebiten.Image, timeLeft float32, score uint) {
	opclock := &ebiten.DrawImageOptions{}
	opclock.GeoM.Translate(float64(config.OffsetX), 8)
	screen.DrawImage(pui.clock, opclock)

	uiTimeText := fmt.Sprintf(": %d", uint(timeLeft))
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(config.OffsetX)+32, 12)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, uiTimeText, pui.textFace, op)

	uiScoreText := fmt.Sprintf("%s: %06d", language.Value.Score, score)
	widthText, _ := text.Measure(uiScoreText, pui.textFace, 0)
	op = &text.DrawOptions{}
	op.GeoM.Translate(float64(config.WindowWidth)-widthText-float64(config.OffsetX), 12)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, uiScoreText, pui.textFace, op)
}

func (pui *PlayUI) UpdateOverlays() {
	for _, pointFeedback := range pui.pointsFeedbaks {
		pointFeedback.Update()
	}

	if len(pui.pointsFeedbaks) == 0 {
		for _, overlay := range pui.overlays {
			overlay.Update()
		}
	}

	pui.postUpdate()
}

func (pui *PlayUI) DrawOverlays(screen *ebiten.Image) {
	for _, pointFeedback := range pui.pointsFeedbaks {
		pointFeedback.Draw(screen, pui.textFace)
	}

	if len(pui.pointsFeedbaks) == 0 {
		for _, overlay := range pui.overlays {
			overlay.Draw(screen, pui.textFace)
		}
	}
}

func (pui *PlayUI) postUpdate() {
	pui.pointsFeedbaks = slices.DeleteFunc(pui.pointsFeedbaks, func(pointFeedback Overlay) bool {
		return pointFeedback.CanRemove()
	})

	if len(pui.pointsFeedbaks) == 0 {
		pui.overlays = slices.DeleteFunc(pui.overlays, func(overlay Overlay) bool {
			canRemove := overlay.CanRemove()
			if canRemove {
				pui.notifier.OnUIOverlayFinished()
			}
			overlay = nil
			return canRemove
		})
	}
}
