package main

import (
	"machine"
	"strconv"
	"time"

	"fw/engine"
)

type Vector2 struct {
	X, Y int
}

type UIElement interface {
	drawTo(*engine.ScreenBuffer)
}

type Text struct {
	font *Font
	text string
	// xPos, yPos int
	pos     Vector2
	colour  engine.Pixel
	visible bool
}

func (t *Text) drawTo(buf *engine.ScreenBuffer) {
	if !t.visible {
		return
	}

	xp := t.pos.X

	for _, char := range textToChars(t.font, t.text) {
		for y, row := range char.content {
			for x, b := range row {
				if b != 0 {
					// skip empty pixels
					buf.Set(xp+x, t.pos.Y+y, t.colour)
				}
			}
		}

		xp += int(char.width) - int(t.font.crush)
	}
}

type Grid struct {
	colour1, colour2            engine.Pixel
	pos, offset, size, cellSize Vector2
	visible                     bool
}

func (g *Grid) drawTo(buf *engine.ScreenBuffer) {
	if !g.visible {
		return
	}

	// no "funny business" around the x and y axes
	for y := 0; y < g.size.Y; y++ {
		for x := 0; x < g.size.X; x++ {
			ox, oy := x+g.offset.X, y+g.offset.Y
			cc := ox/g.cellSize.X + oy/g.cellSize.Y
			if ox < 0 {
				cc--
			}
			if oy < 0 {
				cc--
			}


			if cc%2 == 0 {
				buf.Set(g.pos.X+x, g.pos.Y+y, g.colour1)
			} else {
				buf.Set(g.pos.X+x, g.pos.Y+y, g.colour2)
			}
		}
	}
}

func splash(en *engine.Engine) {
	freq := machine.CPUFrequency()

	engineText := &Text{FontDex, "Grips Engine", Vector2{25, 5}, red, true} // I'm calling it this because it's an anagram of Sprig
	buildText := &Text{FontUnifont, "build 20", Vector2{5, 35}, blue, true}
	freqText := &Text{FontUnifont, strconv.Itoa(int(freq/1_000_000)) + "MHz", Vector2{5, 50}, green, true}

	ui := []UIElement{engineText, buildText, freqText}
	if freq >= 270_000_000 {
		ui = append(ui, &Text{FontUnifont, "OVERCLOCKED!", Vector2{5, 65}, red, true})
	} else {
		ui = append(ui, &Text{FontUnifont, "(could be better)", Vector2{5, 65}, white, true})
	}

	for _, e := range ui {
		e.drawTo(en.ScreenBuffer)
	}
	en.Render()

	time.Sleep(time.Second)
}

func main() {
	// display things
	en := engine.New()

	// splash(en)

	var xPos, yPos int

	// event loop i guess
	for {
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

		x := &Text{FontUnifont, strconv.Itoa(xPos), Vector2{2, 2}, green, true}
		y := &Text{FontUnifont, strconv.Itoa(yPos), Vector2{2, 22}, green, true}

		Texts := [engine.ButtonsCount]*Text{
			{FontDex, "W", Vector2{2 + 20 - 2, 2}, red, true},
			{FontDex, "A", Vector2{2, 26}, red, true},
			{FontDex, "S", Vector2{2 + 20, 50}, red, true},
			{FontDex, "D", Vector2{2 + 40, 26}, red, true},
			{FontDex, "I", Vector2{109 + 20 + 1, 2}, red, true},
			{FontDex, "J", Vector2{109, 26}, red, true},
			{FontDex, "K", Vector2{109 + 20, 50}, red, true},
			{FontDex, "L", Vector2{109 + 40, 26}, red, true},
		}

		// read button states
		for i, b := range en.Buttons {
			if b.Pressed() {
				Texts[i].colour = green
			} else {
				Texts[i].colour = red
			}
		}

		// draw a grid
		grid := &Grid{
			colour1:  grey1,
			colour2:  grey2,
			pos:      Vector2{0, 70},
			offset:   Vector2{xPos, yPos},
			size:     Vector2{engine.Width, engine.Height - 70},
			cellSize: Vector2{10, 10},
			visible:  true,
		}
		grid.drawTo(en.ScreenBuffer)

		ui := []UIElement{x, y}
		for _, t := range Texts {
			ui = append(ui, t)
		}

		for _, e := range ui {
			e.drawTo(en.ScreenBuffer)
		}

		en.Render()
	}
}
