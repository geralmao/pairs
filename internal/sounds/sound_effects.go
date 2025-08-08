package sounds

import (
	"bytes"
	"fmt"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/programatta/pairs/internal/assets/sounds"
)

type soundType int

const (
	t_wav soundType = iota
	t_ogg
)

type SoundEffects struct {
	clickButton     *audio.Player
	flipCard        *audio.Player
	targetFound     *audio.Player
	completedLevel  *audio.Player
	failedLevel     *audio.Player
	bonusPoint      *audio.Player
	backgroundMusic *audio.Player
}

func NewSoundEffects() *SoundEffects {
	soundEffects := &SoundEffects{}

	const sampleRate = 44100
	audioContext := audio.NewContext(sampleRate)

	soundEffects.clickButton = loadSound(t_ogg, audioContext, loadFx("clickButton.ogg"))
	soundEffects.flipCard = loadSound(t_ogg, audioContext, loadFx("flipCard.ogg"))
	soundEffects.targetFound = loadSound(t_ogg, audioContext, loadFx("targetFound.ogg"))
	soundEffects.completedLevel = loadSound(t_ogg, audioContext, loadFx("completedLevel.ogg"))
	soundEffects.failedLevel = loadSound(t_ogg, audioContext, loadFx("failedLevel.ogg"))
	soundEffects.bonusPoint = loadSound(t_ogg, audioContext, loadFx("bonusPoint.ogg"))
	soundEffects.backgroundMusic = loadSound(t_ogg, audioContext, loadMusic("bgFlightHome.ogg"))

	return soundEffects
}

func (se SoundEffects) PlayFlipCard() {
	se.resetPlayer(se.flipCard)
}

func (se SoundEffects) PlayTargetFound() {
	se.resetPlayer(se.targetFound)
}

func (se SoundEffects) PlayCompletedLevel() {
	se.resetPlayer(se.completedLevel)
}

func (se SoundEffects) PlayFailedLevel() {
	se.resetPlayer(se.failedLevel)
}

func (se SoundEffects) PlayClickButton() {
	se.resetPlayer(se.clickButton)
}

func (se SoundEffects) PlayBonusPoint() {
	se.resetPlayer(se.bonusPoint)
}

func (se SoundEffects) PlayBackgroundMusic() {
	if !se.backgroundMusic.IsPlaying() {
		se.resetPlayer(se.backgroundMusic)
		se.backgroundMusic.SetVolume(0.3)
	}
}

func (se SoundEffects) StopBackgroundMusic() {
	if se.backgroundMusic.IsPlaying() {
		se.backgroundMusic.Pause()
	}
}

func (se SoundEffects) resetPlayer(player *audio.Player) {
	player.Rewind()
	player.Play()
}

// loadFx loads sound fx resource from embed file system and returns a byte array.
func loadFx(soundFxName string) []byte {
	sourceSound, sourceErr := sounds.SoundFxFS.ReadFile(fmt.Sprintf("fx/%s", soundFxName))
	if sourceErr != nil {
		panic(sourceErr)
	}
	return sourceSound
}

// loadMusic loads music resource from embed file system and returns a byte array.
func loadMusic(soundMusicName string) []byte {
	sourceSound, sourceErr := sounds.BackgroundMusicFS.ReadFile(fmt.Sprintf("music/%s", soundMusicName))
	if sourceErr != nil {
		panic(sourceErr)
	}
	return sourceSound
}

// loadSound decodes a source sound with sample rate context and return a player.
func loadSound(soundType soundType, audioContext *audio.Context, sourceSound []byte) *audio.Player {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}
	var soundStream audioStream = nil
	var decodeErr error

	switch soundType {
	case t_wav:
		soundStream, decodeErr = wav.DecodeWithSampleRate(audioContext.SampleRate(), bytes.NewReader(sourceSound))
		if decodeErr != nil {
			panic(decodeErr)
		}
	case t_ogg:
		soundStream, decodeErr = vorbis.DecodeWithSampleRate(audioContext.SampleRate(), bytes.NewReader(sourceSound))
		if decodeErr != nil {
			panic(decodeErr)
		}
	}

	player, playerErr := audioContext.NewPlayer(soundStream)
	if playerErr != nil {
		panic(playerErr)
	}
	return player
}
