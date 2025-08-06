package board

import (
	"crypto/md5"
	"fmt"
	"image/color"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/programatta/pairs/internal/assets/images"
	"github.com/programatta/pairs/internal/config"
	"github.com/programatta/pairs/internal/utils"
	"github.com/programatta/pairs/internal/views/play/common"
)

type Card struct {
	id         string
	posX       float64
	posY       float64
	front      *ebiten.Image
	back       *ebiten.Image
	isBack     bool
	isFlipping bool
	flipTime   float32
	notifier   common.Notifier
}

func NewCard(posX, posY float64, emojiName string, notifier common.Notifier) *Card {
	card := &Card{posX: posX, posY: posY, isBack: true, notifier: notifier}

	hashMd5 := md5.New()
	io.WriteString(hashMd5, emojiName)
	card.id = fmt.Sprintf("%v", hashMd5.Sum(nil))
	card.back = ebiten.NewImage(94, 94)
	card.back.Fill(color.NRGBA{0x00, 0xff, 0x00, 0xff})
	if emojiName != "empty" {
		emojiBytes, err := images.EmojisDataFS.ReadFile(emojiName)
		if err != nil {
			panic(err)
		}
		card.front = utils.GenerateImage(94, 94, emojiBytes)
	} else {
		card.front = ebiten.NewImage(94, 94)
		card.front.Fill(color.NRGBA{0xcf, 0xba, 0xf0, 0xff})
	}

	return card
}

// ----------------------------------------------------------------------------
// Implementa ICard Interface
// ----------------------------------------------------------------------------

func (c *Card) Id() string {
	return c.id
}

func (c *Card) Update() {
	if c.isFlipping {
		c.flipTime += config.Dt
		if c.flipTime >= flipDelay {
			c.flipTime = 0
			c.isFlipping = false
			c.isBack = !c.isBack
			if c.isBack {
				c.notifier.OnCardFlipped(c)
			}
		}
	}
}

func (c *Card) Draw(screen *ebiten.Image) {
	var imageTmp *ebiten.Image = nil

	op := &ebiten.DrawImageOptions{}

	// cardWidth := c.front.Bounds().Dx()
	// cardHeight := c.front.Bounds().Dy()

	if c.isFlipping {
		process := float64(c.flipTime) / float64(flipDelay)
		scaleX := 1.0 - process*2
		// fmt.Printf("\nscaleX:%f - process:%f", scaleX, process)
		if process > 0.5 {
			scaleX = -scaleX
		}

		op.GeoM.Scale(float64(scaleX), 1.0)

		if process > 0.5 {
			if c.isBack {
				imageTmp = c.front
			} else {
				imageTmp = c.back
			}
		} else {
			if c.isBack {
				imageTmp = c.back
			} else {
				imageTmp = c.front
			}
		}
	} else {
		if c.isBack {
			imageTmp = c.back
		} else {
			imageTmp = c.front
		}
	}
	op.GeoM.Translate(c.posX, c.posY)
	screen.DrawImage(imageTmp, op)
}

// ----------------------------------------------------------------------------
// Implementa ICard(Flipper) Interface
// ----------------------------------------------------------------------------

func (c *Card) IsFlipping() bool {
	return c.isFlipping
}

func (c *Card) DoFlip() {
	c.isFlipping = true
	c.flipTime = 0
	if c.isBack {
		c.notifier.OnCardFlipped(c)
	}
}

func (c *Card) IsFaceUpCard() bool {
	return !c.isBack
}

// ----------------------------------------------------------------------------
// Implementa ICard(Collider) Interface
// ----------------------------------------------------------------------------

func (c *Card) Rect() (float64, float64, float64, float64) {
	return c.posX, c.posY, float64(c.front.Bounds().Dx()), float64(c.front.Bounds().Dy())
}

const flipDelay float32 = 0.35
