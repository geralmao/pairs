package loader

import (
	"bytes"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/programatta/pairs/internal/assets/images"
	"github.com/programatta/pairs/internal/assets/lang"
	"github.com/programatta/pairs/internal/config"
	"github.com/programatta/pairs/internal/language"
	"github.com/programatta/pairs/internal/platform"
	"github.com/programatta/pairs/internal/storage"
	"github.com/programatta/pairs/internal/views"
)

type LoaderView struct {
	nextViewId views.ViewId
	image      *ebiten.Image
	time       float32
}

func NewLoaderView() *LoaderView {

	imgTmp, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(images.LoaderData))
	if err != nil {
		panic(err)
	}

	return &LoaderView{
		nextViewId: views.Loader,
		image:      imgTmp,
	}
}

// ----------------------------------------------------------------------------
// Implements Viewer Interface
// ----------------------------------------------------------------------------

func (lv *LoaderView) Start(context *config.GameContext) {
	gameData := loadGameData()

	context.IsFxActive = gameData.FxActive
	context.IsSoundActive = gameData.MusicActive
	context.Language = gameData.Language
	if context.Language == "" {
		fmt.Println("Language not found... getting current from system...")
		language := getCurrentLanguage()
		fmt.Printf("\nLanguage found [%s]", language)

		context.Language = language
		storage.SaveGameData(&storage.GameData{
			FxActive:    context.IsFxActive,
			MusicActive: context.IsSoundActive,
			Language:    context.Language,
		})
	} else {
		fmt.Printf("\nLanguage found in data [%s]\n", context.Language)
	}
	language.LoadById(context.Language)
}

func (lv *LoaderView) ProcessEvents() {}

func (lv *LoaderView) Update() {
	lv.time += config.Dt
	if lv.time > splashDelay {
		lv.nextViewId = views.Menu
	}
}

func (lv *LoaderView) Draw(screen *ebiten.Image) {
	screen.Fill(config.BackgroundColorApplication)
	screen.DrawImage(lv.image, &ebiten.DrawImageOptions{})
}

func (lv *LoaderView) NextView() views.ViewId {
	return lv.nextViewId
}

// -----------------------------------------------------------------------------
// Helper functions
// -----------------------------------------------------------------------------

func loadGameData() *storage.GameData {
	gameData, errGameData := storage.LoadGameData()
	if errGameData != nil {
		fmt.Printf("\nError load game data: [%v]\n\n", errGameData)
		fmt.Printf("\nCreating a new game data...\n")
		gameData = &storage.GameData{
			FxActive:    true,
			MusicActive: true,
		}

		errGameData = storage.SaveGameData(gameData)
		if errGameData != nil {
			fmt.Printf("\nError save game data: [%v]\n\n", errGameData)
		}
	}

	return gameData
}

func getCurrentLanguage() string {
	langIdDetected := platform.GetSystemLanguage()
	if !lang.IsLangIdSupported(langIdDetected) {
		langIdDetected = lang.LangDefault
	}
	return langIdDetected
}

const splashDelay float32 = 0.75
