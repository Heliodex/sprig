package main

import (
	"fw/game"
	"fw/engine"
)

func main() {
	en := engine.New()
	game.Splash(en)

	for state := (&game.State{}); ; {
		state.Update(en)
	}
}
