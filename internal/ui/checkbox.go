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

type Checkbox struct {
	posX             float64
	posY             float64
	check            *ebiten.Image
	checkActive      *ebiten.Image
	bgColor          color.Color
	checkedColor     color.Color
	checked          bool
	callbackNotifier func(checked bool)
	context          *config.GameContext
}

func NewCheckbox(posX, posY, width, height float64, checked bool, bgColor color.Color, checkedColor color.Color) *Checkbox {
	checkbox := &Checkbox{
		posX:         posX,
		posY:         posY,
		checked:      checked,
		bgColor:      bgColor,
		checkedColor: checkedColor,
	}
	checkbox.check = ebiten.NewImage(int(width), int(height))

	widthM := width * 80 / 100
	heightM := height * 80 / 100

	checkbox.checkActive = ebiten.NewImage(int(widthM), int(heightM))
	return checkbox
}

func (cb *Checkbox) SetContext(context *config.GameContext) {
	cb.context = context
}

func (cb *Checkbox) SetValue(value bool) {
	cb.checked = value
}

func (cb *Checkbox) Update() {
	cb.processEvents()
}

func (cb *Checkbox) Draw(screen *ebiten.Image) {
	cb.check.Fill(cb.bgColor)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(cb.posX, cb.posY)

	if cb.checked {
		cb.checkActive.Fill(cb.checkedColor)
		opAct := &ebiten.DrawImageOptions{}
		opAct.GeoM.Translate(float64(cb.check.Bounds().Dx())/2-float64(cb.checkActive.Bounds().Dx())/2, float64(cb.check.Bounds().Dy())/2-float64(cb.checkActive.Bounds().Dy())/2)
		cb.check.DrawImage(cb.checkActive, opAct)
	}

	screen.DrawImage(cb.check, op)
}

func (cb *Checkbox) OnClick(onCallback func(checked bool)) {
	cb.callbackNotifier = onCallback
}

// ----------------------------------------------------------------------------
// Implementa Collider Interface
// ----------------------------------------------------------------------------

func (cb *Checkbox) Rect() (float64, float64, float64, float64) {
	return cb.posX, cb.posY, float64(cb.check.Bounds().Dx()), float64(cb.check.Bounds().Dy())
}

func (cb *Checkbox) processEvents() {
	x, y := platform.PressPosition()
	x, y = utils.GetPositionInGameCoords(x, y, cb.context)

	inside := collider.CheckPointInsideRect(float64(x), float64(y), cb)

	if platform.IsPressEventJustRelease() && inside {
		cb.checked = !cb.checked
		if cb.callbackNotifier != nil {
			cb.callbackNotifier(cb.checked)
		}
	}
}

// -----------------------------------------------------------------------------
// Checkbox with text
// -----------------------------------------------------------------------------
type CheckboxText struct {
	Checkbox
	title       string
	textFace    *text.GoTextFace
	isTextRight bool
	titlePosX   float64
	titlePosY   float64
}

func NewCheckboxWithText(posX, posY, width, height float64, checked bool, bgColor color.Color, checkedColor color.Color, title string, textFace *text.GoTextFace, isTextRight bool) *CheckboxText {

	textX := posX
	textY := posY

	widthText, heightText := text.Measure(title, textFace, 0)
	if isTextRight {
		textX = posX + width + 10.0
		textY = posY + height/2 - heightText/3
	} else {
		textX = posX
		textY = posY + height/2 - heightText/3
		posX = textX + widthText + 10.0
	}

	checkbox := &CheckboxText{
		Checkbox:    *NewCheckbox(posX, posY, width, height, checked, bgColor, checkedColor),
		title:       title,
		textFace:    textFace,
		isTextRight: isTextRight,
		titlePosX:   textX,
		titlePosY:   textY,
	}
	return checkbox
}

func (cbt *CheckboxText) Draw(screen *ebiten.Image) {
	cbt.Checkbox.Draw(screen)

	opText := &text.DrawOptions{}
	opText.GeoM.Translate(cbt.titlePosX, cbt.titlePosY)
	opText.ColorScale.ScaleWithColor(color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF})
	text.Draw(screen, cbt.title, cbt.textFace, opText)
}
