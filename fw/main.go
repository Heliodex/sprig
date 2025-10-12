package main

import (
	"strconv"

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

	var xPos, yPos int

	buildText := &Text{FontUnifont, "build 19", 2, 2, blue, true}

	x := &Text{FontUnifont, "0", 2, 22, green, true}
	y := &Text{FontUnifont, "0", 2, 42, green, true}

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

	ui := []UIElement{buildText, x, y}
	for _, t := range Texts {
		ui = append(ui, t)
	}

	// event loop i guess
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
			yPos--
		}

		if en.Buttons[engine.S].Pressed() {
			yPos++
		}

		if en.Buttons[engine.A].Pressed() {
			xPos--
		}

		if en.Buttons[engine.D].Pressed() {
			xPos++
		}

		x.text = strconv.Itoa(xPos)
		y.text = strconv.Itoa(yPos)

		for _, e := range ui {
			e.drawTo(en.ScreenBuffer)
		}

		en.Render()
	}
}
