package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/programatta/pairs/internal/collider"
	"github.com/programatta/pairs/internal/config"
	"github.com/programatta/pairs/internal/platform"
	"github.com/programatta/pairs/internal/utils"
)

type buttonState = int

const (
	normal buttonState = iota
	onOver
	onPress
)

type Button struct {
	posX             float64
	posY             float64
	title            string
	button           *ebiten.Image
	color            color.Color
	textColor        color.Color
	textFace         *text.GoTextFace
	state            buttonState
	callbackNotifier func()
	context          *config.GameContext
}

func NewButton(posX, posY, width, height float64, title string, textFace *text.GoTextFace) *Button {
	playButton := &Button{posX: posX, posY: posY, title: title}
	playButton.button = ebiten.NewImage(int(width), int(height))
	playButton.color = color.NRGBA{0x7f, 0xcb, 0xbb, 0xff}
	playButton.textColor = color.NRGBA{0xff, 0xff, 0xff, 0xff}
	playButton.textFace = textFace
	playButton.state = normal
	return playButton
}

func (pb *Button) SetContext(context *config.GameContext) {
	pb.context = context
}

func (pb *Button) Update() {
	pb.processEvents()

	switch pb.state {
	case onOver:
		pb.color = color.NRGBA{0xb9, 0xfb, 0xc0, 0xff}
	case onPress:
		pb.color = color.NRGBA{0x90, 0xdb, 0xf4, 0xff}
	default:
		pb.color = color.NRGBA{0x7f, 0xcb, 0xbb, 0xff}
	}
}

func (pb *Button) Draw(screen *ebiten.Image) {
	pb.button.Fill(pb.color)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pb.posX, pb.posY)

	widthText, heightText := text.Measure(pb.title, pb.textFace, 0)

	opText := &text.DrawOptions{}
	opText.GeoM.Translate(float64(pb.button.Bounds().Dx())/2-widthText/2, float64(pb.button.Bounds().Dy())/2-heightText/3)
	op.ColorScale.ScaleWithColor(pb.textColor)
	text.Draw(pb.button, pb.title, pb.textFace, opText)

	screen.DrawImage(pb.button, op)
}

func (pb *Button) OnClick(onCallback func()) {
	pb.callbackNotifier = onCallback
}

// ----------------------------------------------------------------------------
// Implementa Collider Interface
// ----------------------------------------------------------------------------

func (pb *Button) Rect() (float64, float64, float64, float64) {
	return pb.posX, pb.posY, float64(pb.button.Bounds().Dx()), float64(pb.button.Bounds().Dy())
}

func (pb *Button) processEvents() {
	pb.state = normal

	//Si el cursor está sobre el botón damos feedback cambiando el color.
	x, y := platform.PressPosition()
	x, y = utils.GetPositionInGameCoords(x, y, pb.context)

	inside := collider.CheckPointInsideRect(float64(x), float64(y), pb)

	if inside {
		pb.state = onOver
		if platform.IsPressEventPressed() {
			pb.state = onPress
		}
	}

	if platform.IsPressEventJustRelease() && inside {
		if pb.callbackNotifier != nil {
			pb.callbackNotifier()
		}
	}
}
