package main

import (
	"fw/engine"
	"fw/game"
)

func main() {
	en := engine.New()
	game.Splash(en)

	for {
		game.Update(en)
	}
}
