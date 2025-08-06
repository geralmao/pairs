package sounds

import "github.com/programatta/pairs/internal/config"

type SoundController struct {
	context      *config.GameContext
	soundEffects *SoundEffects
}

type FXSounds int

const (
	BonusPoints FXSounds = iota
	FlipCard
	ClickButton
	TargetFound
	CompletedLevel
	FailedLevel
)

func NewSoundController(context *config.GameContext) *SoundController {
	soundController := &SoundController{
		context: context,
	}

	soundController.soundEffects = NewSoundEffects()

	return soundController
}

func (sc SoundController) PlayBackgroundMusic() {
	if sc.context.IsSoundActive {
		sc.soundEffects.PlayBackgroundMusic()
	}
}

func (sc SoundController) StopBackgroundMusic() {
	sc.soundEffects.StopBackgroundMusic()
}

func (sc SoundController) PlayFx(fxPlay FXSounds) {
	if sc.context.IsFxActive {
		switch fxPlay {
		case BonusPoints:
			sc.soundEffects.PlayBonusPoint()
		case FlipCard:
			sc.soundEffects.PlayFlipCard()
		case ClickButton:
			sc.soundEffects.PlayClickButton()
		case TargetFound:
			sc.soundEffects.PlayTargetFound()
		case CompletedLevel:
			sc.soundEffects.PlayCompletedLevel()
		case FailedLevel:
			sc.soundEffects.PlayFailedLevel()
		}
	}
}
