package settings

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/programatta/pairs/internal/config"
	"github.com/programatta/pairs/internal/language"
	"github.com/programatta/pairs/internal/sounds"
	"github.com/programatta/pairs/internal/storage"
	"github.com/programatta/pairs/internal/ui"
	"github.com/programatta/pairs/internal/views"
)

type SettingsView struct {
	textFace       *text.GoTextFace
	soundCtrl      *sounds.SoundController
	nextViewId     views.ViewId
	context        *config.GameContext
	contextCopy    config.GameContext
	muteFxChbox    *ui.CheckboxText
	muteSoundChbox *ui.CheckboxText
	acceptBtn      *ui.Button
	cancelBtn      *ui.Button
}

func NewSettingsView(textFace *text.GoTextFace, soundCtrl *sounds.SoundController) *SettingsView {
	settingsView := &SettingsView{
		textFace:   textFace,
		soundCtrl:  soundCtrl,
		nextViewId: views.Settings,
	}

	muteTxt := fmt.Sprintf("%s:", language.Value.Mute)
	wMuteTxt, _ := text.Measure(muteTxt, textFace, 0)

	settingsView.muteFxChbox = ui.NewCheckboxWithText(
		float64(config.WindowWidth)/2-(wMuteTxt+30)/2,
		286,
		30, 30,
		false,
		color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF},
		color.NRGBA{0x7F, 0xCB, 0xBB, 0xFF},
		muteTxt,
		textFace,
		false,
	)

	settingsView.muteSoundChbox = ui.NewCheckboxWithText(
		float64(config.WindowWidth)/2-(wMuteTxt+30)/2,
		386,
		30, 30,
		false,
		color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF},
		color.NRGBA{0x7F, 0xCB, 0xBB, 0xFF},
		muteTxt,
		textFace,
		false,
	)

	settingsView.acceptBtn = ui.NewButton(
		float64(config.WindowWidth)/2-200,
		float64(config.WindowHeight)-110,
		190, 70, language.Value.Accept,
		settingsView.textFace,
	)

	settingsView.cancelBtn = ui.NewButton(
		float64(config.WindowWidth)/2+10,
		float64(config.WindowHeight)-110,
		190, 70, language.Value.Cancel,
		settingsView.textFace,
	)

	return settingsView
}

// ----------------------------------------------------------------------------
// Implements Viewer Interface
// ----------------------------------------------------------------------------

func (sv *SettingsView) Start(context *config.GameContext) {
	sv.context = context
	sv.contextCopy = *context
	sv.nextViewId = views.Settings

	sv.muteFxChbox.SetValue(!sv.context.IsFxActive)
	sv.muteFxChbox.SetContext(sv.context)
	sv.muteSoundChbox.SetValue(!sv.context.IsSoundActive)
	sv.muteSoundChbox.SetContext(sv.context)

	sv.acceptBtn.SetContext(sv.context)
	sv.cancelBtn.SetContext(sv.context)
}

func (sv *SettingsView) ProcessEvents() {
	sv.muteFxChbox.OnClick(func(isChecked bool) {
		sv.soundCtrl.PlayFx(sounds.ClickButton)
		sv.context.IsFxActive = !isChecked
	})

	sv.muteSoundChbox.OnClick(func(isChecked bool) {
		sv.soundCtrl.PlayFx(sounds.ClickButton)
		sv.context.IsSoundActive = !isChecked
	})

	sv.acceptBtn.OnClick(func() {
		sv.nextViewId = views.Menu
		sv.soundCtrl.PlayFx(sounds.ClickButton)

		storage.SaveGameData(&storage.GameData{
			FxActive:    sv.context.IsFxActive,
			MusicActive: sv.context.IsSoundActive,
			Language:    sv.context.Language,
		})
	})

	sv.cancelBtn.OnClick(func() {
		sv.nextViewId = views.Menu
		sv.soundCtrl.PlayFx(sounds.ClickButton)
		*sv.context = sv.contextCopy
	})
}

func (sv *SettingsView) Update() {
	sv.muteFxChbox.Update()
	sv.muteSoundChbox.Update()
	sv.acceptBtn.Update()
	sv.cancelBtn.Update()
}

func (sv *SettingsView) Draw(screen *ebiten.Image) {
	uiSettingsHeaderText := language.Value.Settings
	widthText, _ := text.Measure(uiSettingsHeaderText, sv.textFace, 0)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(config.WindowWidth)/2-widthText/2, 100)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, uiSettingsHeaderText, sv.textFace, op)

	//Fx
	uiFxText := fmt.Sprintf("%s:", language.Value.Sounds)
	widthText, _ = text.Measure(uiFxText, sv.textFace, 0)
	op = &text.DrawOptions{}
	op.GeoM.Translate(float64(config.WindowWidth)/2-widthText/2, 250)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, uiFxText, sv.textFace, op)

	//Fx:Mute
	sv.muteFxChbox.Draw(screen)

	//Music
	uiMusicText := fmt.Sprintf("%s:", language.Value.Music)
	widthText, _ = text.Measure(uiMusicText, sv.textFace, 0)
	op = &text.DrawOptions{}
	op.GeoM.Translate(float64(config.WindowWidth)/2-widthText/2, 350)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, uiMusicText, sv.textFace, op)

	//Music:Mute
	sv.muteSoundChbox.Draw(screen)

	//Music:Volume
	uiMusicVolumeText := fmt.Sprintf("%s:", language.Value.Volume)
	widthText, _ = text.Measure(uiMusicVolumeText, sv.textFace, 0)
	op = &text.DrawOptions{}
	op.GeoM.Translate(float64(config.WindowWidth)/2-widthText/2, 430)
	op.ColorScale.ScaleWithColor(color.NRGBA{0xA7, 0xA8, 0xA7, 0xFF})
	text.Draw(screen, uiMusicVolumeText, sv.textFace, op)

	sv.acceptBtn.Draw(screen)
	sv.cancelBtn.Draw(screen)
}

func (sv *SettingsView) NextView() views.ViewId {
	return sv.nextViewId
}
