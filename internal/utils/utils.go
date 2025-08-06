package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/programatta/pairs/internal/assets/fonts"
	"github.com/programatta/pairs/internal/config"
)

// GenerateImage crea una imagen con el emoji cargado y centrado en la imagen devuelta.
func GenerateImage(width, heigh int, emojiBytes []byte) *ebiten.Image {
	imgTmp, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(emojiBytes))
	if err != nil {
		panic(err)
	}

	opTmp := &ebiten.DrawImageOptions{}
	opTmp.GeoM.Translate(float64(width)/2-float64(imgTmp.Bounds().Dx())/2, float64(heigh)/2-float64(imgTmp.Bounds().Dy())/2)

	img := ebiten.NewImage(width, heigh)
	img.Fill(color.White)
	img.DrawImage(imgTmp, opTmp)

	return img
}

func LoadEmbeddedFont(size float64) *text.GoTextFace {
	faceSource, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.FontFiles))
	if err != nil {
		log.Fatal(err)
	}

	return &text.GoTextFace{
		Source: faceSource,
		Size:   size,
	}
}

func RandomSeed() (uint64, uint64) {
	var b [16]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic("no se pudo generar semilla aleatoria")
	}
	seed1 := binary.LittleEndian.Uint64(b[:8])
	seed2 := binary.LittleEndian.Uint64(b[8:])
	return seed1, seed2
}

func GetPositionInGameCoords(x, y int, gameContext *config.GameContext) (int, int) {
	mx, my := float64(x), float64(y)

	gameX := mx
	gameY := my
	if gameContext != nil {
		// Restar el offset y dividir entre escala para obtener coord. internas
		gameX = (mx - gameContext.OffsetX) / gameContext.Scale
		gameY = (my - gameContext.OffsetY) / gameContext.Scale
	}

	return int(gameX), int(gameY)
}

//-----------------------------------------------------------------------------
// Funciones easing
//-----------------------------------------------------------------------------

func EaseOutQuad(t float64) float64 {
	if t > 1 {
		t = 1
	}
	return 1 - (1-t)*(1-t)
}

func EaseInQuint(t float64) float64 {
	if t < 0 {
		return 0
	}
	if t > 1 {
		return 1
	}
	return t * t * t * t * t
}

func EaseInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	f := -2*t + 2
	return 1 - 0.5*f*f*f
}
