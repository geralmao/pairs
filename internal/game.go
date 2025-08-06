package internal

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/programatta/pairs/internal/config"
	"github.com/programatta/pairs/internal/language"
	"github.com/programatta/pairs/internal/sounds"
	"github.com/programatta/pairs/internal/utils"
	"github.com/programatta/pairs/internal/views"
	"github.com/programatta/pairs/internal/views/loader"
	"github.com/programatta/pairs/internal/views/menu"
	"github.com/programatta/pairs/internal/views/play"
	"github.com/programatta/pairs/internal/views/settings"
)

// Capturamos la versión desde -ldflags "-X main.Version=$(VERSION)" desde el makefile.
var Version = "dev"

type gameTaskId = int

const (
	loaderTask gameTaskId = iota
	runTask
)

type Game struct {
	appViews      map[views.ViewId]views.Viewer
	currentView   views.Viewer
	currentViewId views.ViewId
	context       *config.GameContext
	executeTask   map[gameTaskId]func()
	currentTaskId gameTaskId
	offscreen     *ebiten.Image
}

func NewGame() *Game {
	fmt.Printf("\nMatch Emojis v%s\n\n", Version)

	game := &Game{}
	game.context = &config.GameContext{
		Version: Version,
	}

	textFace := utils.LoadEmbeddedFont(32)
	soundController := sounds.NewSoundController(game.context)

	language.Init()

	game.appViews = make(map[views.ViewId]views.Viewer)
	game.appViews[views.Loader] = loader.NewLoaderView()
	// game.appViews[views.Menu] = menu.NewMenuView(textFace, soundController)
	// game.appViews[views.Play] = play.NewPlayView(textFace, soundController)
	// game.appViews[views.Settings] = settings.NewSettingsView(textFace, soundController)

	//Primera vista en aparecer.
	game.currentViewId = views.Loader
	game.currentView = game.appViews[game.currentViewId]

	game.currentTaskId = loaderTask
	game.executeTask = make(map[gameTaskId]func())
	game.executeTask[loaderTask] = func() {
		game.doLoader()

		game.appViews[views.Menu] = menu.NewMenuView(textFace, soundController)
		game.appViews[views.Play] = play.NewPlayView(textFace, soundController)
		game.appViews[views.Settings] = settings.NewSettingsView(textFace, soundController)
	}
	game.executeTask[runTask] = game.doRun
	game.offscreen = ebiten.NewImage(config.WindowWidth, config.WindowHeight)
	return game
}

// ----------------------------------------------------------------------------
// Implementa Ebiten Game Interface
// ----------------------------------------------------------------------------

// Update realiza el cambio de estado si es necesario y permite procesar
// eventos y actualizar su lógica.
func (g *Game) Update() error {
	g.executeTask[g.currentTaskId]()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(config.BackgroundColorApplication)

	g.offscreen.Clear()
	g.offscreen.Fill(config.BackgroundColorApplication)
	g.currentView.Draw(g.offscreen)

	// Dibuja el lienzo lógico escalado y centrado en la pantalla real
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.context.Scale, g.context.Scale)
	op.GeoM.Translate(g.context.OffsetX, g.context.OffsetY)
	op.Filter = ebiten.FilterLinear

	screen.DrawImage(g.offscreen, op)
}

// Layout determina el tamaño del canvas
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if g.context.Scale == 0 {
		g.context.OutsideWidth = outsideWidth
		g.context.OutsideHeight = outsideHeight

		// Calcula escala y offset para centrar
		sw := g.offscreen.Bounds().Dx()
		sh := g.offscreen.Bounds().Dy()
		dw := outsideWidth
		dh := outsideHeight

		scaleX := float64(dw) / float64(sw)
		scaleY := float64(dh) / float64(sh)
		g.context.Scale = math.Min(scaleX, scaleY)

		g.context.OffsetX = (float64(dw) - float64(sw)*g.context.Scale) / 2
		g.context.OffsetY = (float64(dh) - float64(sh)*g.context.Scale) / 2
	}
	return outsideWidth, outsideHeight
}

func (g *Game) doLoader() {
	g.currentView.Start(g.context)
	g.currentTaskId = runTask
}

func (g *Game) doRun() {
	nextViewId := g.currentView.NextView()
	if nextViewId != g.currentViewId {
		g.currentView = views.NewTransitionView(g.appViews[g.currentViewId], g.appViews[nextViewId], func() {
			g.currentView = g.appViews[g.currentViewId]
		})
		g.currentView.Start(g.context)
		g.currentViewId = nextViewId
	}

	g.currentView.ProcessEvents()
	g.currentView.Update()
}
