package views

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/programatta/pairs/internal/config"
	"github.com/programatta/pairs/internal/utils"
)

type TransitionView struct {
	from                 Viewer
	to                   Viewer
	onFinishedTransition func()
	transitionTime       float32
}

func NewTransitionView(from, to Viewer, onFinishedTransition func()) *TransitionView {
	transitionView := &TransitionView{
		from:                 from,
		to:                   to,
		onFinishedTransition: onFinishedTransition,
	}

	return transitionView
}

// ----------------------------------------------------------------------------
// Implements Viewer Interface
// ----------------------------------------------------------------------------

func (tv *TransitionView) Start(context *config.GameContext) {
	tv.to.Start(context)
}

func (tv *TransitionView) ProcessEvents() {}

func (tv *TransitionView) Update() {
	if tv.transitionTime == transitionDelay {
		if tv.onFinishedTransition != nil {
			tv.onFinishedTransition()
			tv.onFinishedTransition = nil
			tv.from = nil
			tv.to = nil
		}
	} else {
		tv.transitionTime += config.Dt
		if tv.transitionTime >= transitionDelay {
			tv.transitionTime = transitionDelay
		}
		tv.to.Update()
	}
}

func (tv *TransitionView) Draw(screen *ebiten.Image) {
	alpha := float32(utils.EaseInOutCubic(float64(tv.transitionTime / transitionDelay)))
	invAlpha := 1.0 - alpha

	// Dibujamos la vista actual (from) con alpha decreciente
	fromImage := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	tv.from.Draw(fromImage)
	opFrom := &ebiten.DrawImageOptions{}
	opFrom.ColorScale.ScaleAlpha(float32(invAlpha))
	screen.DrawImage(fromImage, opFrom)

	// Dibujamos la vista nueva (to) con alpha creciente
	toImage := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	tv.to.Draw(toImage)
	opTo := &ebiten.DrawImageOptions{}
	opTo.ColorScale.ScaleAlpha(alpha)

	screen.DrawImage(toImage, opTo)
}

func (tv *TransitionView) NextView() ViewId {
	return tv.to.NextView()
}

const transitionDelay float32 = 0.75 // 0.52 //primera: 0.92
