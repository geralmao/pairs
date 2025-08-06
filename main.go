//go:build !android

package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/programatta/pairs/internal"
	"github.com/programatta/pairs/internal/config"
)

func main() {
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)
	ebiten.SetWindowTitle("Match Emojis")

	game := internal.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
