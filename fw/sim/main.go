package main

import (
	"os"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"

	"fw/game"
	"fw/sim/engine"
	// "fw/util"
)

// these are Dvorak because I'm a nerd
var codeMap = map[key.Code]uint8{
	key.CodeComma: 0,
	key.CodeA:     1,
	key.CodeO:     2,
	key.CodeE:     3,
	key.CodeC:     4,
	key.CodeH:     5,
	key.CodeT:     6,
	key.CodeN:     7,
}

func startUI(s screen.Screen) {
	en := engine.New(s)

	go func() {
		// keypress events fire on both press and release
		// all we need to do is check whether a key is pressed at a given time
		// if a key is held, multiple press events will be fired, however these are easy to tell apart because normal fast keypresses have a few milliseconds between them, whereas held keys only have microseconds or nanoseconds between them (I guess we can realistically hope that no frame occurs in this time)
		for {
			switch e := en.Window.NextEvent().(type) {
			case key.Event:
				btn, ok := codeMap[e.Code]
				if !ok {
					continue
				}

				en.ButtonImpls[btn].Set(e.Direction == key.DirPress)
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					os.Exit(0)
				}
			}
		}
	}()

	game.Splash(en)

	for state := (&game.State{}); ; {
		state.Update(en)
		// wait 1 frame
		time.Sleep(40 * time.Millisecond) // ~25fps, which is how fast my console actually runs, ymmv
	}
}

func main() {
	driver.Main(startUI)
}
