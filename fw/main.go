package main

import (
	"strconv"
	"time"

	"fw/engine"
)

type UIElement interface {
	drawTo(*engine.ScreenBuffer)
}

type Text struct {
	font       *Font
	text       string
	xPos, yPos int
	colour     engine.Pixel
	visible    bool
}

func (t *Text) drawTo(buf *engine.ScreenBuffer) {
	if !t.visible {
		return
	}

	xp := t.xPos

	for _, char := range textToChars(t.font, t.text) {
		for y, row := range char.content {
			for x, b := range row {
				if b != 0 {
					// skip empty pixels
					buf.Set(xp+x, t.yPos+y, t.colour)
				}
			}
		}

		xp += int(char.width) - int(t.font.crush)
	}
}

func main() {
	// display things
	en := engine.New()

	buildText := &Text{FontUnifont, "build 19", 2, 2, blue, true}
	fps := &Text{FontUnifont, "0", 2, 22, green, true}

	Texts := [engine.ButtonsCount]*Text{
		{FontDex, "W", 2 + 20 - 1, 24, red, true},
		{FontDex, "A", 2, 48, red, true},
		{FontDex, "S", 2 + 20, 72, red, true},
		{FontDex, "D", 2 + 40, 48, red, true},
		{FontDex, "I", engine.Width/2 + 2 + 20 + 1, 24, red, true},
		{FontDex, "J", engine.Width/2 + 2, 48, red, true},
		{FontDex, "K", engine.Width/2 + 2 + 20, 72, red, true},
		{FontDex, "L", engine.Width/2 + 2 + 40, 48, red, true},
	}

	ui := []UIElement{buildText}
	for _, t := range Texts {
		ui = append(ui, t)
	}
	ui = append(ui, fps)

	// event loop i guess
	lastDifferences := []int{}
	lastTime := time.Now()
	for {
		// read button states
		for i, b := range en.Buttons {
			if b.Pressed() {
				Texts[i].colour = green
			} else {
				Texts[i].colour = red
			}
		}

		if en.Buttons[engine.W].Pressed() {
			buildText.yPos--
		}

		if en.Buttons[engine.S].Pressed() {
			buildText.yPos++
		}

		if en.Buttons[engine.A].Pressed() {
			buildText.xPos--
		}

		if en.Buttons[engine.D].Pressed() {
			buildText.xPos++
		}

		for _, e := range ui {
			e.drawTo(en.ScreenBuffer)
		}

		lastDifferences = append(lastDifferences, int(time.Since(lastTime).Milliseconds()))
		if len(lastDifferences) > 10 {
			lastDifferences = lastDifferences[1:]
		}
		
		total := 0
		for _, v := range lastDifferences {
			total += v
		}
		avg := float64(total) / float64(len(lastDifferences))
		if avg == 0 {
			avg = 1
		}
		fps.text = strconv.Itoa(int(1000.0 / avg))
		
		lastTime = time.Now()

		en.Render()
	}
}
