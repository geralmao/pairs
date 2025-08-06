package common

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/programatta/pairs/internal/collider"
)

type Notifier interface {
	OnCardFlipped(card ICard)
	OnUIOverlayFinished()
	OnUIOverlayScoreBonus() (uint, uint)
}

type Flipper interface {
	IsFaceUpCard() bool
	IsFlipping() bool
	DoFlip()
}

type ICard interface {
	collider.Collider
	Flipper
	Id() string
	Draw(screen *ebiten.Image)
	Update()
}
