package game

import (
	"math"
	"strconv"
	"time"

	"fw/util"
)

type Vector2 struct {
	X, Y int
}

type UIElement interface {
	drawTo(*util.ScreenBuffer)
}

type Text struct {
	font *Font
	text string
	// xPos, yPos int
	pos    Vector2
	colour util.Pixel
}

func (t *Text) drawTo(buf *util.ScreenBuffer) {
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
	colour1, colour2            util.Pixel
	pos, offset, size, cellSize Vector2
}

func (g *Grid) drawTo(buf *util.ScreenBuffer) {
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
	colour                       util.Pixel
	pos                          Vector2
	amplitude, wavelength, phase int
}

func (s *SineWave) drawTo(buf *util.ScreenBuffer) {
	for x := s.pos.X; x < util.Width; x++ {
		y := s.pos.Y + int(float64(s.amplitude)*math.Sin(float64(x)/float64(s.wavelength)*2*math.Pi+float64(s.phase)/10))
		buf.Set(x, y, s.colour)
	}
}

// intro to show, in the event that something else is loading or to show information
func Splash(en util.Engine) {
	freq := en.CPUFrequency()

	engineText := &Text{FontDex, "Grips Engine", Vector2{25, 5}, util.Red} // I'm calling it this because it's an anagram of Sprig
	buildText := &Text{FontUnifont, "build 21", Vector2{5, 35}, util.Blue}
	freqText := &Text{FontUnifont, strconv.Itoa(int(freq/1_000_000)) + "MHz", Vector2{5, 50}, util.Green}

	ui := []UIElement{engineText, buildText, freqText}
	if freq >= 270_000_000 {
		ui = append(ui, &Text{FontUnifont, "OVERCLOCKED!", Vector2{5, 65}, util.Red})
	} else {
		ui = append(ui, &Text{FontUnifont, "(could be better)", Vector2{5, 65}, util.White})
	}

	for _, e := range ui {
		e.drawTo(en.ScreenBuffer())
	}
	en.Render()

	time.Sleep(time.Second)
}

type State struct {
	f   int
	pos Vector2
}

// ran every frame (or, more like this is what makes the frames)
func (s *State) Update(en util.Engine) {
	btns := en.Buttons()

	if btns[util.W].Pressed() {
		s.pos.Y--
	}

	if btns[util.S].Pressed() {
		s.pos.Y++
	}

	if btns[util.A].Pressed() {
		s.pos.X--
	}

	if btns[util.D].Pressed() {
		s.pos.X++
	}

	x := &Text{FontUnifont, strconv.Itoa(s.pos.X), Vector2{2, 2}, util.Green}
	y := &Text{FontUnifont, strconv.Itoa(s.pos.Y), Vector2{2, 22}, util.Green}

	Texts := [util.ButtonsCount]*Text{
		{FontDex, "W", Vector2{2 + 20 - 2, 2}, util.Red},
		{FontDex, "A", Vector2{2, 26}, util.Red},
		{FontDex, "S", Vector2{2 + 20, 50}, util.Red},
		{FontDex, "D", Vector2{2 + 40, 26}, util.Red},
		{FontDex, "I", Vector2{109 + 20 + 1, 2}, util.Red},
		{FontDex, "J", Vector2{109, 26}, util.Red},
		{FontDex, "K", Vector2{109 + 20, 50}, util.Red},
		{FontDex, "L", Vector2{109 + 40, 26}, util.Red},
	}

	// read button states
	for i, b := range btns {
		if b.Pressed() {
			Texts[i].colour = util.Green
		} else {
			Texts[i].colour = util.Red
		}
	}

	grid := &Grid{
		colour1:  util.Grey1,
		colour2:  util.Grey2,
		pos:      Vector2{0, 70},
		offset:   s.pos,
		size:     Vector2{util.Width, util.Height - 70},
		cellSize: Vector2{10, 10},
	}

	// sine wave
	sine := &SineWave{
		colour:     util.Cyan,
		pos:        Vector2{0, util.Height * 0.75},
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
		e.drawTo(en.ScreenBuffer())
	}

	en.Render()
	s.f++
}
