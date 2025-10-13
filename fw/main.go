package main

import (
	"machine"
	"math"
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
	pos    Vector2
	colour engine.Pixel
}

func (t *Text) drawTo(buf *engine.ScreenBuffer) {
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
}

func (g *Grid) drawTo(buf *engine.ScreenBuffer) {
	// no "funny business" around the x and y axes
	for y := range g.size.Y {
		for x := range g.size.X {
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

type SineWave struct {
	colour                       engine.Pixel
	pos                          Vector2
	amplitude, wavelength, phase int
}

func (s *SineWave) drawTo(buf *engine.ScreenBuffer) {
	for x := s.pos.X; x < engine.Width; x++ {
		y := s.pos.Y + int(float64(s.amplitude)*math.Sin(float64(x)/float64(s.wavelength)*2*math.Pi+float64(s.phase)/10))
		buf.Set(x, y, s.colour)
	}
}

// intro to show, in the event that something else is loading or to show information
func splash(en *engine.Engine) {
	freq := machine.CPUFrequency()

	engineText := &Text{FontDex, "Grips Engine", Vector2{25, 5}, red} // I'm calling it this because it's an anagram of Sprig
	buildText := &Text{FontUnifont, "build 21", Vector2{5, 35}, blue}
	freqText := &Text{FontUnifont, strconv.Itoa(int(freq/1_000_000)) + "MHz", Vector2{5, 50}, green}

	ui := []UIElement{engineText, buildText, freqText}
	if freq >= 270_000_000 {
		ui = append(ui, &Text{FontUnifont, "OVERCLOCKED!", Vector2{5, 65}, red})
	} else {
		ui = append(ui, &Text{FontUnifont, "(could be better)", Vector2{5, 65}, white})
	}

	for _, e := range ui {
		e.drawTo(en.ScreenBuffer)
	}
	en.Render()

	time.Sleep(time.Second)
}

type State struct {
	f   int
	pos Vector2
}

// ran every frame (or, more like this is what makes the frames)
func (s *State) Update(en *engine.Engine) {
	if en.Buttons[engine.W].Pressed() {
		s.pos.Y--
	}

	if en.Buttons[engine.S].Pressed() {
		s.pos.Y++
	}

	if en.Buttons[engine.A].Pressed() {
		s.pos.X--
	}

	if en.Buttons[engine.D].Pressed() {
		s.pos.X++
	}

	x := &Text{FontUnifont, strconv.Itoa(s.pos.X), Vector2{2, 2}, green}
	y := &Text{FontUnifont, strconv.Itoa(s.pos.Y), Vector2{2, 22}, green}

	Texts := [engine.ButtonsCount]*Text{
		{FontDex, "W", Vector2{2 + 20 - 2, 2}, red},
		{FontDex, "A", Vector2{2, 26}, red},
		{FontDex, "S", Vector2{2 + 20, 50}, red},
		{FontDex, "D", Vector2{2 + 40, 26}, red},
		{FontDex, "I", Vector2{109 + 20 + 1, 2}, red},
		{FontDex, "J", Vector2{109, 26}, red},
		{FontDex, "K", Vector2{109 + 20, 50}, red},
		{FontDex, "L", Vector2{109 + 40, 26}, red},
	}

	// read button states
	for i, b := range en.Buttons {
		if b.Pressed() {
			Texts[i].colour = green
		} else {
			Texts[i].colour = red
		}
	}

	grid := &Grid{
		colour1:  grey1,
		colour2:  grey2,
		pos:      Vector2{0, 70},
		offset:   s.pos,
		size:     Vector2{engine.Width, engine.Height - 70},
		cellSize: Vector2{10, 10},
	}

	// sine wave
	sine := &SineWave{
		colour:     cyan,
		pos:        Vector2{0, engine.Height * 0.75},
		amplitude:  8,
		wavelength: 40,
		phase:      s.f,
	}

	ui := []UIElement{grid, x, y}
	for _, t := range Texts {
		ui = append(ui, t)
	}
	ui = append(ui, sine)

	for _, e := range ui {
		e.drawTo(en.ScreenBuffer)
	}

	en.Render()
	s.f++
}

func main() {
	en := engine.New()
	splash(en)

	for state := (&State{}); ; {
		state.Update(en)
	}
}
