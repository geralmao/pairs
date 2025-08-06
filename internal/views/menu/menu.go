package menu

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/programatta/pairs/internal/config"
	"github.com/programatta/pairs/internal/language"
	"github.com/programatta/pairs/internal/platform"
	"github.com/programatta/pairs/internal/sounds"
	"github.com/programatta/pairs/internal/ui"
	"github.com/programatta/pairs/internal/views"
)

type MenuView struct {
	textFace    *text.GoTextFace
	soundCtrl   *sounds.SoundController
	nextViewId  views.ViewId
	playBtn     *ui.Button
	settingsBtn *ui.Button
	exitBtn     *ui.Button
	version     string
}

func NewMenuView(textFace *text.GoTextFace, soundCtrl *sounds.SoundController) *MenuView {
	menuView := &MenuView{
		textFace:   textFace,
		soundCtrl:  soundCtrl,
		nextViewId: views.Menu,
	}

	menuView.playBtn = ui.NewButton(
		float64(config.WindowWidth)/2-115,
		float64(config.WindowHeight)/2-130,
		230, 70, language.Value.Play,
		menuView.textFace,
	)

	menuView.settingsBtn = ui.NewButton(
		float64(config.WindowWidth)/2-115,
		float64(config.WindowHeight)/2,
		230, 70, language.Value.Settings,
		menuView.textFace,
	)

	menuView.exitBtn = ui.NewButton(
		float64(config.WindowWidth)/2-115,
		float64(config.WindowHeight)/2+130,
		230, 70, language.Value.Exit,
		menuView.textFace,
	)

	return menuView
}

// ----------------------------------------------------------------------------
// Implements Viewer Interface
// ----------------------------------------------------------------------------

func (mv *MenuView) Start(context *config.GameContext) {
	mv.nextViewId = views.Menu
	mv.playBtn.SetContext(context)
	mv.settingsBtn.SetContext(context)
	mv.exitBtn.SetContext(context)
	mv.version = context.Version
}

func (mv *MenuView) ProcessEvents() {
	mv.playBtn.OnClick(func() {
		mv.nextViewId = views.Play
		mv.soundCtrl.PlayFx(sounds.ClickButton)
	})
	mv.settingsBtn.OnClick(func() {
		mv.nextViewId = views.Settings
		mv.soundCtrl.PlayFx(sounds.ClickButton)
	})
	mv.exitBtn.OnClick(func() {
		mv.soundCtrl.PlayFx(sounds.ClickButton)
		go func() {
			time.Sleep(time.Millisecond * 250)
			platform.ExitGame()
		}()
	})
}

func (mv *MenuView) Update() {
	mv.playBtn.Update()
	mv.settingsBtn.Update()
	mv.exitBtn.Update()
}

func (mv *MenuView) Draw(screen *ebiten.Image) {
	uiMenuHeaderText := language.Value.Menu
	widthText, _ := text.Measure(uiMenuHeaderText, mv.textFace, 0)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(config.WindowWidth)/2-widthText/2, 100)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, uiMenuHeaderText, mv.textFace, op)

	mv.playBtn.Draw(screen)
	mv.settingsBtn.Draw(screen)
	mv.exitBtn.Draw(screen)
	mv.drawVersion(screen)
}

func (mv *MenuView) NextView() views.ViewId {
	return mv.nextViewId
}

func (mv *MenuView) drawVersion(screen *ebiten.Image) {
	version := fmt.Sprintf("v%s ", mv.version)

	size := mv.textFace.Size
	mv.textFace.Size = 18

	widthVersionText, heightVersionText := text.Measure(version, mv.textFace, 0)

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(config.WindowWidth)-widthVersionText, float64(config.WindowHeight)-heightVersionText)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, version, mv.textFace, op)

	mv.textFace.Size = size
}
