package play

import (
	"fmt"
	"image/color"
	"math/rand/v2"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/programatta/pairs/internal/assets/images"
	"github.com/programatta/pairs/internal/collider"
	"github.com/programatta/pairs/internal/config"
	"github.com/programatta/pairs/internal/language"
	"github.com/programatta/pairs/internal/platform"
	"github.com/programatta/pairs/internal/sounds"
	"github.com/programatta/pairs/internal/ui"
	"github.com/programatta/pairs/internal/utils"
	"github.com/programatta/pairs/internal/views"
	"github.com/programatta/pairs/internal/views/play/board"
	"github.com/programatta/pairs/internal/views/play/common"
	pui "github.com/programatta/pairs/internal/views/play/ui"
)

//-----------------------------------------------------------------------------
// Estados de juego
//-----------------------------------------------------------------------------

type gameState uint

const (
	play gameState = iota
	gameover
	wingame
	transition
)

//-----------------------------------------------------------------------------
// Subestados de juego (PLAY)
//-----------------------------------------------------------------------------

type playInnerState uint

const (
	targetPresentation playInnerState = iota
	selecting
	waiting
	checking
	unflipping
	newgame
	feedbackFinishedLevel
)

type PlayView struct {
	textFace            *text.GoTextFace
	soundCtrl           *sounds.SoundController
	nextViewId          views.ViewId
	lastGameState       gameState
	gameState           gameState
	playInnerState      playInnerState
	levels              []*Level
	timeLeft            float32
	playBtn             *ui.Button
	backBtn             *ui.Button
	playUI              *pui.PlayUI
	score               uint
	cards               []common.ICard
	cardsFlipped        []common.ICard
	waitingTime         float32
	faceUpCardCount     int
	currentLevel        uint
	context             *config.GameContext
	transitionFromDraw  func(*ebiten.Image)
	transitionToDraw    func(*ebiten.Image)
	transitionTime      float32
	wingameMessagePos   int
	wingameMessageWidth float64
	fromImg             *ebiten.Image
	toImg               *ebiten.Image
}

func NewPlayView(textFace *text.GoTextFace, soundCtrl *sounds.SoundController) *PlayView {
	playView := &PlayView{
		textFace:       textFace,
		soundCtrl:      soundCtrl,
		nextViewId:     views.Play,
		lastGameState:  gameover,
		gameState:      play,
		playInnerState: newgame,
	}

	playView.levels = loadLevels()

	playView.timeLeft = playView.levels[playView.currentLevel].Time()
	playView.playBtn = ui.NewButton(
		float64(config.WindowWidth)/2-135,
		float64(config.WindowHeight)/2+130,
		270, 70, language.Value.PressToPlay,
		playView.textFace,
	)
	playView.backBtn = ui.NewButton(
		float64(config.WindowWidth)/2-135,
		float64(config.WindowHeight)/2+220,
		270, 70, language.Value.BackMenu,
		playView.textFace,
	)
	playView.playUI = pui.NewPlayUI(playView.textFace, playView)

	return playView
}

// ----------------------------------------------------------------------------
// Implements Viewer Interface
// ----------------------------------------------------------------------------

func (pv *PlayView) Start(context *config.GameContext) {
	pv.nextViewId = views.Play
	pv.gameState = play
	pv.lastGameState = gameover
	pv.playInnerState = newgame
	pv.context = context
	pv.playBtn.SetContext(pv.context)
	pv.backBtn.SetContext(pv.context)
	pv.wingameMessagePos = config.WindowWidth + 5
	pv.wingameMessageWidth, _ = text.Measure(language.Value.WallOfFame, pv.textFace, 0)
}

func (pv *PlayView) ProcessEvents() {
	//Eventos
	switch pv.gameState {
	case play:
		pv.processEventsPlay()
	case gameover:
		pv.processEventsGameOver()
	case wingame:
		pv.processEventsWinGame()
	}
}

func (pv *PlayView) Update() {
	pv.soundCtrl.PlayBackgroundMusic()

	//Actualización.
	switch pv.gameState {
	case play:
		pv.updatePlay()
	case gameover:
		pv.updateGameOver()
	case wingame:
		pv.updateWinGame()
	case transition:
		pv.updateTransition()
	}
}

func (pv *PlayView) Draw(screen *ebiten.Image) {
	switch pv.gameState {
	case play:
		pv.drawPlay(screen)
	case gameover:
		pv.drawGameOver(screen)
	case wingame:
		pv.drawWinGame(screen)
	case transition:
		pv.drawTransition(screen)
	}
}

func (pv *PlayView) NextView() views.ViewId {
	return pv.nextViewId
}

// ----------------------------------------------------------------------------
// Implements Notifier Interface
// ----------------------------------------------------------------------------

func (pv *PlayView) OnCardFlipped(card common.ICard) {
	switch pv.playInnerState {
	case selecting:
		pv.cardsFlipped = append(pv.cardsFlipped, card)
		if len(pv.cardsFlipped) == pv.levels[pv.currentLevel].Matchs() { //2
			pv.playInnerState = waiting
		}
	case unflipping:
		isFlipping := false
		for _, cardFlipped := range pv.cardsFlipped {
			isFlipping = isFlipping || cardFlipped.IsFlipping()
		}
		if !isFlipping {
			pv.emptyCardsFlipped()
			pv.playInnerState = selecting
		}
	}
}

func (pv *PlayView) OnUIOverlayFinished() {
	switch pv.playInnerState {
	case targetPresentation:
		pv.playInnerState = selecting
	case feedbackFinishedLevel:
		pv.lastGameState = play
		pv.playInnerState = newgame
	}
}

func (pv *PlayView) OnUIOverlayScoreBonus() (uint, uint) {
	pv.score += config.ScoreBonusPoints
	if uint(pv.timeLeft)-1 > 0 {
		pv.timeLeft -= 1
	} else {
		pv.timeLeft = 0
		pv.emptyBoard()
	}
	pv.soundCtrl.PlayFx(sounds.BonusPoints)
	return pv.score, uint(pv.timeLeft)
}

//-----------------------------------------------------------------------------
// Procesa eventos segun estado
//-----------------------------------------------------------------------------

func (pv *PlayView) processEventsPlay() {
	if pv.playInnerState == selecting {
		if platform.IsPressEventJustRelease() {
			x, y := platform.PressPosition()
			x, y = utils.GetPositionInGameCoords(x, y, pv.context)

			//Transformamos las coordenadas físicas a logicas.
			xx := (x - int(config.OffsetX)) / 101
			yy := (y - int(config.OffsetY)) / 101

			//Transformamos xx e yy en posición del slice.
			pos := yy*4 + xx

			card := pv.cards[pos]
			if !card.IsFaceUpCard() && collider.CheckPointInsideRect(float64(x), float64(y), card) {
				if !card.IsFlipping() {
					card.DoFlip()
					pv.soundCtrl.PlayFx(sounds.FlipCard)
				}
			}
		}
	}
}

func (pv *PlayView) processEventsGameOver() {
	pv.playBtn.OnClick(func() {
		pv.emptyBoard()
		pv.startInternalTransition(pv.drawGameOver, pv.drawPlay)
		pv.soundCtrl.PlayFx(sounds.ClickButton)
	})

	pv.backBtn.OnClick(func() {
		pv.nextViewId = views.Menu
		pv.soundCtrl.PlayFx(sounds.ClickButton)
		pv.soundCtrl.StopBackgroundMusic()
	})
}

func (pv *PlayView) processEventsWinGame() {
	pv.backBtn.OnClick(func() {
		pv.nextViewId = views.Menu
		pv.soundCtrl.PlayFx(sounds.ClickButton)
		pv.soundCtrl.StopBackgroundMusic()
	})
}

// -----------------------------------------------------------------------------
// Update estado interno según estado
// -----------------------------------------------------------------------------

func (pv *PlayView) updatePlay() {
	switch pv.playInnerState {
	case waiting:
		pv.waitingTime += config.Dt
		if pv.waitingTime > waitingDelay {
			pv.waitingTime = 0
			pv.playInnerState = checking
		}
	case checking:
		lastId := ""
		idCardsMatchCount := 1
		for _, cardFlipped := range pv.cardsFlipped {
			if lastId == "" {
				lastId = cardFlipped.Id()
			} else {
				if lastId == cardFlipped.Id() {
					idCardsMatchCount += 1
				} else {
					break
				}
			}
		}

		if idCardsMatchCount == pv.levels[pv.currentLevel].Matchs() {
			pv.emptyCardsFlipped()
			pv.faceUpCardCount++
			pv.score += config.CardPoints

			pv.soundCtrl.PlayFx(sounds.TargetFound)
			pv.playUI.AddFeedbackNewPoints(config.CardPoints)

			if pv.levels[pv.currentLevel].Target() == pv.faceUpCardCount {
				pv.playInnerState = feedbackFinishedLevel
				pv.playUI.AddOverlayFinishLevel(pv.timeLeft, pv.score)
				pv.soundCtrl.PlayFx(sounds.CompletedLevel)
			} else {
				pv.playInnerState = selecting
			}
		} else {
			pv.playInnerState = unflipping
			for _, cardFlipped := range pv.cardsFlipped {
				cardFlipped.DoFlip()
				pv.soundCtrl.PlayFx(sounds.FlipCard)
			}
		}
	case newgame:
		if pv.lastGameState == gameover {
			pv.currentLevel = 0
			pv.score = 0
		} else {
			if pv.currentLevel+1 < uint(len(pv.levels)) {
				pv.currentLevel++
			} else {
				pv.startInternalTransition(pv.drawPlay, pv.drawWinGame)
				return
			}
		}
		pv.emptyCardsFlipped()
		pv.cards = createBoard(pv, pv.levels[pv.currentLevel])
		pv.faceUpCardCount = 0
		pv.waitingTime = 0
		pv.timeLeft = pv.levels[pv.currentLevel].Time()

		pv.playUI.AddOverlayPresentation(pv.currentLevel+1, pv.levels[pv.currentLevel].Target(), matchsDescription[pv.levels[pv.currentLevel].Matchs()])

		pv.playInnerState = targetPresentation
	}

	for _, card := range pv.cardsFlipped {
		card.Update()
	}

	if pv.playInnerState != targetPresentation && pv.playInnerState != feedbackFinishedLevel && pv.playInnerState != newgame {
		pv.timeLeft -= config.Dt
		if pv.timeLeft <= 0 {
			pv.timeLeft = 0
			pv.startInternalTransition(pv.drawPlay, pv.drawGameOver)
			pv.soundCtrl.PlayFx(sounds.FailedLevel)
		}
	}

	pv.playUI.UpdateOverlays()
}

func (pv *PlayView) updateGameOver() {
	pv.playBtn.Update()
	pv.backBtn.Update()
}

func (pv *PlayView) updateWinGame() {
	pv.wingameMessagePos -= 1
	if pv.wingameMessagePos+int(pv.wingameMessageWidth) < 0 {
		pv.wingameMessagePos = config.WindowWidth + 5
	}
	pv.backBtn.Update()
}

func (pv *PlayView) updateTransition() {
	if pv.transitionTime == transitionDelay {
		pv.transitionTime = 0
		pv.fromImg = nil
		pv.toImg = nil
		switch pv.lastGameState {
		case play:
			if pv.playInnerState == newgame {
				pv.gameState = wingame
			} else {
				pv.gameState = gameover
			}
		case gameover:
			pv.gameState = play
			pv.playInnerState = newgame
		}
	} else {
		pv.transitionTime += config.Dt
		if pv.transitionTime > transitionDelay {
			pv.transitionTime = transitionDelay
		}
	}
}

//-----------------------------------------------------------------------------
// Dibuja en pantalla según estado
//-----------------------------------------------------------------------------

func (pv *PlayView) drawPlay(screen *ebiten.Image) {
	pv.playUI.DrawHeader(screen, pv.timeLeft, pv.score)

	//Board
	for pos := range maxBoardCards {
		y := pos / 4
		x := pos % 4
		vector.StrokeRect(screen, float32(x*101)+config.OffsetX, float32(y*101)+config.OffsetY, cardWidth, cardHeight, 1, color.NRGBA{0xff, 0x00, 0x00, 0xff}, true)
	}

	if pv.playInnerState != newgame {
		for _, card := range pv.cards {
			card.Draw(screen)
		}
	}

	//UI extra
	pv.playUI.DrawOverlays(screen)
}

func (pv *PlayView) drawGameOver(screen *ebiten.Image) {
	uiGameOverText := language.Value.GameOver
	widthText, _ := text.Measure(uiGameOverText, pv.textFace, 0)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(config.WindowWidth)/2-widthText/2, float64(config.WindowHeight)/2-100)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, uiGameOverText, pv.textFace, op)

	uiScoreText := fmt.Sprintf("%s: %06d", language.Value.YourScore, pv.score)
	widthText, _ = text.Measure(uiScoreText, pv.textFace, 0)
	op = &text.DrawOptions{}
	op.GeoM.Translate(float64(config.WindowWidth)/2-widthText/2, float64(config.WindowHeight)/2+10)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, uiScoreText, pv.textFace, op)

	pv.playBtn.Draw(screen)
	pv.backBtn.Draw(screen)
}

func (pv *PlayView) drawWinGame(screen *ebiten.Image) {
	uiWinGameText := language.Value.YouWin
	widthText, _ := text.Measure(uiWinGameText, pv.textFace, 0)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(config.WindowWidth)/2-widthText/2, float64(config.WindowHeight)/2-100)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, uiWinGameText, pv.textFace, op)

	uiMessageText := language.Value.WallOfFame
	op = &text.DrawOptions{}
	op.GeoM.Translate(float64(pv.wingameMessagePos), float64(config.WindowHeight)/2-50)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, uiMessageText, pv.textFace, op)

	uiScoreText := fmt.Sprintf("%s: %06d", language.Value.YourScore, pv.score)
	widthText, _ = text.Measure(uiScoreText, pv.textFace, 0)
	op = &text.DrawOptions{}
	op.GeoM.Translate(float64(config.WindowWidth)/2-widthText/2, float64(config.WindowHeight)/2+10)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, uiScoreText, pv.textFace, op)

	pv.backBtn.Draw(screen)
}

func (pv *PlayView) drawTransition(screen *ebiten.Image) {
	alpha := utils.EaseInOutCubic(float64(pv.transitionTime / transitionDelay))

	// From state
	if pv.fromImg == nil {
		pv.fromImg = ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	}
	pv.fromImg.Clear()
	pv.transitionFromDraw(pv.fromImg)
	opFrom := &ebiten.DrawImageOptions{}
	opFrom.ColorScale.ScaleAlpha(float32(1 - alpha))
	screen.DrawImage(pv.fromImg, opFrom)

	// To state
	if pv.toImg == nil {
		pv.toImg = ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	}
	pv.toImg.Clear()
	pv.transitionToDraw(pv.toImg)
	opTo := &ebiten.DrawImageOptions{}
	opTo.ColorScale.ScaleAlpha(float32(alpha))
	screen.DrawImage(pv.toImg, opTo)
}

func (pv *PlayView) emptyCardsFlipped() {
	pv.cardsFlipped = slices.DeleteFunc(pv.cardsFlipped, func(card common.ICard) bool {
		return true
	})
}

func (pv *PlayView) emptyBoard() {
	pv.cards = slices.DeleteFunc(pv.cards, func(card common.ICard) bool {
		return true
	})
}

func (pv *PlayView) startInternalTransition(fromDraw, toDraw func(*ebiten.Image)) {
	pv.lastGameState = pv.gameState
	pv.gameState = transition
	pv.transitionFromDraw = fromDraw
	pv.transitionToDraw = toDraw
}

// Funciones helper

func createBoard(notifier common.Notifier, level *Level) []common.ICard {
	emojis := loadRandonEmojis(level.EmojisCount(), level.Matchs())

	cards := []common.ICard{}
	for pos := range maxBoardCards {
		y := pos / 4
		x := pos % 4
		emojiName := emojis[pos]
		card := board.NewCard(float64(x)*101+float64(config.OffsetX)+3, float64(y)*101+float64(config.OffsetY)+3, emojiName, notifier)
		cards = append(cards, card)
	}
	return cards
}

func loadRandonEmojis(emojisCount int, times int) []string {
	rnd := rand.New(rand.NewPCG(utils.RandomSeed()))

	var emojis []string

	//obtenemos los nombres de los ficheros.
	filenames := images.GetEmojisNamesFromFS()

	for range emojisCount {
		pos := rnd.IntN(len(filenames))
		emojiFileName := filenames[pos]
		filenames = slices.DeleteFunc(filenames, func(str string) bool {
			return str == emojiFileName
		})

		for range times {
			emojis = append(emojis, emojiFileName)
		}
	}

	if len(emojis) < maxBoardCards {
		for range maxBoardCards - len(emojis) {
			emojis = append(emojis, "empty")
		}
	}

	for range 100 {
		rnd.Shuffle(len(emojis), func(i, j int) {
			emojis[i], emojis[j] = emojis[j], emojis[i]
		})
	}
	return emojis
}

func loadLevels() []*Level {
	levels := []*Level{}

	//TODO: añadir premio si se consigue antes de un tiempo establecido.
	//Dos fichas emparejadas - 14 emojis.
	levels = append(levels, NewLevel(2, 45, 14, 2))
	levels = append(levels, NewLevel(3, 60, 14, 2))
	levels = append(levels, NewLevel(4, 60, 14, 2))
	levels = append(levels, NewLevel(5, 60, 14, 2))
	levels = append(levels, NewLevel(6, 70, 14, 2))
	levels = append(levels, NewLevel(7, 75, 14, 2))
	levels = append(levels, NewLevel(8, 75, 14, 2))
	levels = append(levels, NewLevel(9, 75, 14, 2))
	levels = append(levels, NewLevel(10, 80, 14, 2))
	levels = append(levels, NewLevel(11, 80, 14, 2))
	levels = append(levels, NewLevel(12, 85, 14, 2))
	levels = append(levels, NewLevel(13, 85, 14, 2))
	levels = append(levels, NewLevel(14, 90, 14, 2))

	//Tres fichas emparejadas - 9 emojis.
	levels = append(levels, NewLevel(2, 100, 9, 3))
	levels = append(levels, NewLevel(3, 115, 9, 3))
	levels = append(levels, NewLevel(4, 130, 9, 3))
	levels = append(levels, NewLevel(5, 145, 9, 3))
	levels = append(levels, NewLevel(6, 160, 9, 3))
	levels = append(levels, NewLevel(7, 175, 9, 3))
	levels = append(levels, NewLevel(8, 190, 9, 3))
	levels = append(levels, NewLevel(9, 205, 9, 3))

	//Cuatro fichas emparejadas - 7 emojis.
	levels = append(levels, NewLevel(2, 215, 7, 4))
	levels = append(levels, NewLevel(3, 230, 7, 4))
	levels = append(levels, NewLevel(4, 245, 7, 4))
	levels = append(levels, NewLevel(5, 260, 7, 4))
	levels = append(levels, NewLevel(6, 275, 7, 4))
	levels = append(levels, NewLevel(7, 290, 7, 4))

	//Cinco fichas emparejadas - 5 emojis.
	levels = append(levels, NewLevel(2, 300, 5, 5))
	levels = append(levels, NewLevel(3, 315, 5, 5))
	levels = append(levels, NewLevel(4, 330, 5, 5))
	levels = append(levels, NewLevel(5, 345, 5, 5))

	//Seis fichas emparejadas - 4 emojis.
	levels = append(levels, NewLevel(2, 355, 4, 6))
	levels = append(levels, NewLevel(3, 370, 4, 6))
	levels = append(levels, NewLevel(4, 395, 4, 6))

	//Siete fichas emparejadas - 4 emojis.
	levels = append(levels, NewLevel(2, 410, 4, 7))
	levels = append(levels, NewLevel(3, 425, 4, 7))
	levels = append(levels, NewLevel(4, 440, 4, 7))

	//Ocho fichas emparejadas - 3 emojis.
	levels = append(levels, NewLevel(2, 450, 3, 8))
	levels = append(levels, NewLevel(3, 465, 3, 8))

	//Nueve fichas emparejadas - 3 emojis.
	levels = append(levels, NewLevel(2, 475, 3, 9))
	levels = append(levels, NewLevel(3, 490, 3, 9))

	matchsDescription = make(map[int]string)
	matchsDescription[2] = language.Value.MatchType2
	matchsDescription[3] = language.Value.MatchType3
	matchsDescription[4] = language.Value.MatchType4
	matchsDescription[5] = language.Value.MatchType5
	matchsDescription[6] = language.Value.MatchType6
	matchsDescription[7] = language.Value.MatchType7
	matchsDescription[8] = language.Value.MatchType8
	matchsDescription[9] = language.Value.MatchType9

	return levels
}

const cardWidth float32 = 100
const cardHeight float32 = 100
const waitingDelay float32 = 0.75
const transitionDelay float32 = 0.65
const maxBoardCards int = 28

var matchsDescription map[int]string
