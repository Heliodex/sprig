package main

import (
	"fw/engine"
	"fw/game"
)

func main() {
	en := engine.New()
	game.Splash(en)

	for state := (&game.State{}); ; {
		state.Update(en)
	}
}
