package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/programatta/pairs/internal/config"
	"github.com/programatta/pairs/internal/language"
	"github.com/programatta/pairs/internal/utils"
	"github.com/programatta/pairs/internal/views/play/common"
)

type finishLevelState int

const (
	show finishLevelState = iota
	bonus
)

type FinishLevelOverlay struct {
	timeLeftShow               uint
	score                      uint
	notifier                   common.Notifier
	feedbackLevelCompletedTime float32
	feedbackCountTime          float32
	remove                     bool
	state                      finishLevelState
}

func NewFinishLevelOverlay(timeLeftShow float32, score uint, notifier common.Notifier) *FinishLevelOverlay {
	return &FinishLevelOverlay{
		timeLeftShow: uint(timeLeftShow),
		score:        score,
		notifier:     notifier,
		state:        show,
	}
}

//-----------------------------------------------------------------------------
// Implements Overlay Interface
//-----------------------------------------------------------------------------

func (flo *FinishLevelOverlay) Update() {
	switch flo.state {
	case show:
		flo.feedbackCountTime += config.Dt
		if flo.feedbackCountTime > feedbackFinishShowDelay {
			flo.feedbackCountTime = 0
			flo.state = bonus
		}
	case bonus:
		if flo.timeLeftShow > 0 {
			flo.feedbackCountTime += config.Dt
			if flo.feedbackCountTime > feedbackCountTimeDelay {
				flo.feedbackCountTime = 0
				flo.score, flo.timeLeftShow = flo.notifier.OnUIOverlayScoreBonus()
			}
		} else {
			flo.timeLeftShow = 0
			flo.feedbackLevelCompletedTime += config.Dt
			if flo.feedbackLevelCompletedTime > feedbackLevelCompletedDelay {
				flo.feedbackLevelCompletedTime = 0
				flo.remove = true
			}
		}
	}
}

func (flo *FinishLevelOverlay) Draw(screen *ebiten.Image, textFace *text.GoTextFace) {
	targWidth := screen.Bounds().Dx()
	targHeight := screen.Bounds().Dy()
	targImg := ebiten.NewImage(targWidth, targHeight)
	targImg.Fill(color.RGBA{0x49, 0x50, 0x57, 205})

	uiLevelText := language.Value.LevelCompleted
	widthText, _ := text.Measure(uiLevelText, textFace, 0)
	opText := &text.DrawOptions{}
	opText.GeoM.Translate(float64(config.WindowWidth)/2-widthText/2, float64(config.WindowHeight)/2-50)
	opText.ColorScale.ScaleWithColor(color.White)
	text.Draw(targImg, uiLevelText, textFace, opText)

	uiBonusTimeText := fmt.Sprintf("%s: %d", language.Value.TimeBonus, uint(flo.timeLeftShow))
	widthText, _ = text.Measure(uiBonusTimeText, textFace, 0)
	opText = &text.DrawOptions{}
	opText.GeoM.Translate(float64(config.WindowWidth)/2-widthText/2, float64(config.WindowHeight)/2)
	opText.ColorScale.ScaleWithColor(color.White)
	text.Draw(targImg, uiBonusTimeText, textFace, opText)

	uiBonusPointsText := fmt.Sprintf("%s: %d", language.Value.Score, flo.score)
	widthText, _ = text.Measure(uiBonusPointsText, textFace, 0)
	opText = &text.DrawOptions{}
	opText.GeoM.Translate(float64(config.WindowWidth)/2-widthText/2, float64(config.WindowHeight)/2+50)
	opText.ColorScale.ScaleWithColor(color.White)
	text.Draw(targImg, uiBonusPointsText, textFace, opText)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	op.ColorScale.ScaleAlpha(1.0 - float32(utils.EaseInQuint(float64(flo.feedbackLevelCompletedTime/feedbackLevelCompletedDelay))))

	screen.DrawImage(targImg, op)
}

func (flo *FinishLevelOverlay) CanRemove() bool {
	return flo.remove
}

const feedbackFinishShowDelay float32 = 0.75
const feedbackCountTimeDelay float32 = 0.05
const feedbackLevelCompletedDelay float32 = 3.0
