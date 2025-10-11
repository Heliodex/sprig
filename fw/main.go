package main

import (
	"fw/st7735"

	"tinygo.org/x/drivers/pixel"
)

type UIElement interface {
	drawTo(buf *st7735.ScreenBuffer)
}

type Text struct {
	font       *Font
	text       string
	xPos, yPos int
	colour     pixel.RGB565BE
	visible    bool
}

func (t *Text) drawTo(buf *st7735.ScreenBuffer) {
	if !t.visible {
		return
	}

	xp := t.xPos

	chars := textToChars(t.font, t.text)
	for _, char := range chars {
		for y, row := range char.content {
			for x, b := range row {
				if b == 0 {
					continue // skip empty pixels
				}

				yl := t.yPos + y
				// r := buf[yl]

				xl := xp + x

				buf[yl][xl] = t.colour
			}
		}

		xp += int(char.width) - int(t.font.crush)
	}
}

func main() {
	// display things
	engine := NewEngine()

	buildText := &Text{fontUnifont, "build 17", 2, 2, blue, true}
	successText := &Text{fontUnifont, "success", 2, 46, green, false}

	Texts := [ButtonsCount]*Text{
		{fontDex, "W", 2 + 20 - 1, 24, red, true},
		{fontDex, "A", 2, 48, red, true},
		{fontDex, "S", 2 + 20, 72, red, true},
		{fontDex, "D", 2 + 40, 48, red, true},
		{fontDex, "I", width/2 + 2 + 20 + 1, 24, red, true},
		{fontDex, "J", width/2 + 2, 48, red, true},
		{fontDex, "K", width/2 + 2 + 20, 72, red, true},
		{fontDex, "L", width/2 + 2 + 40, 48, red, true},
	}

	ui := []UIElement{buildText, successText}
	for _, t := range Texts {
		ui = append(ui, t)
	}

	// event loop i guess
	for {
		// read button states
		for i, b := range engine.Buttons {
			if b.Pressed() {
				Texts[i].colour = green
			} else {
				Texts[i].colour = red
			}
		}

		for _, e := range ui {
			e.drawTo(engine.screenmem)
		}

		engine.display.Render(engine.screenmem)
	}
}
