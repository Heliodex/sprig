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

	buildText := &Text{fontUnifont, "build 16", 2, 2, blue, true}
	successText := &Text{fontUnifont, "success", 2, 46, green, false}

	btW := &Text{fontDex, "W", 2 + 20 - 1, 24, red, true}
	btA := &Text{fontDex, "A", 2, 48, red, true}
	btS := &Text{fontDex, "S", 2 + 20, 72, red, true}
	btD := &Text{fontDex, "D", 2 + 40, 48, red, true}
	btI := &Text{fontDex, "I", width/2 + 2 + 20 + 1, 24, red, true}
	btJ := &Text{fontDex, "J", width/2 + 2, 48, red, true}
	btK := &Text{fontDex, "K", width/2 + 2 + 20, 72, red, true}
	btL := &Text{fontDex, "L", width/2 + 2 + 40, 48, red, true}

	ui := []UIElement{
		buildText, successText,
		btW, btA, btS, btD,
		btI, btJ, btK, btL,
	}

	// event loop i guess
	for {
		// read button states
		if engine.W.Pressed() {
			btW.colour = white
		} else {
			btW.colour = red
		}
		if engine.A.Pressed() {
			btA.colour = white
		} else {
			btA.colour = red
		}
		if engine.S.Pressed() {
			btS.colour = white
		} else {
			btS.colour = red
		}
		if engine.D.Pressed() {
			btD.colour = white
		} else {
			btD.colour = red
		}
		if engine.I.Pressed() {
			btI.colour = white
		} else {
			btI.colour = red
		}
		if engine.J.Pressed() {
			btJ.colour = white
		} else {
			btJ.colour = red
		}
		if engine.K.Pressed() {
			btK.colour = white
		} else {
			btK.colour = red
		}
		if engine.L.Pressed() {
			btL.colour = white
		} else {
			btL.colour = red
		}

		for _, e := range ui {
			e.drawTo(engine.screenmem)
		}

		// d.Render(engine.screenmem)
		engine.display.Render(engine.screenmem)
		clear(engine.screenmem[:])
	}
}
