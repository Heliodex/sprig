package game

import (
	"math"
	"strconv"
	"time"

	"fw/util"
)

type UIElement interface {
	drawTo(*util.ScreenBuffer)
}

type Text struct {
	font *Font
	text string
	// xPos, yPos int
	pos    util.Vector2
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
	pos, offset, size, cellSize util.Vector2
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
	pos                          util.Vector2
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

	engineText := &Text{FontDex, "Grips Engine", util.Vector2{X: 25, Y: 5}, util.Red} // I'm calling it this because it's an anagram of Sprig
	buildText := &Text{FontUnifont, "build 21", util.Vector2{X: 5, Y: 35}, util.Blue}
	freqText := &Text{FontUnifont, strconv.Itoa(int(freq/1_000_000)) + "MHz", util.Vector2{X: 5, Y: 50}, util.Green}

	ui := []UIElement{engineText, buildText, freqText}
	if freq >= 270_000_000 {
		ui = append(ui, &Text{FontUnifont, "OVERCLOCKED!", util.Vector2{X: 5, Y: 65}, util.Red})
	} else {
		ui = append(ui, &Text{FontUnifont, "(could be better)", util.Vector2{X: 5, Y: 65}, util.White})
	}

	for _, e := range ui {
		e.drawTo(en.ScreenBuffer())
	}
	en.Render()

	time.Sleep(time.Second)
}

type State struct {
	f   int
	pos util.Vector2
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

	x := &Text{FontUnifont, strconv.Itoa(s.pos.X), util.Vector2{X: 70, Y: 2}, util.Green}
	y := &Text{FontUnifont, strconv.Itoa(s.pos.Y), util.Vector2{X: 70, Y: 22}, util.Green}

	Texts := [util.ButtonsCount]*Text{
		{FontDex, "W", util.Vector2{X: 2 + 20 - 2, Y: 2}, util.Red},
		{FontDex, "A", util.Vector2{X: 2, Y: 26}, util.Red},
		{FontDex, "S", util.Vector2{X: 2 + 20, Y: 50}, util.Red},
		{FontDex, "D", util.Vector2{X: 2 + 40, Y: 26}, util.Red},
		{FontDex, "I", util.Vector2{X: 109 + 20 + 1, Y: 2}, util.Red},
		{FontDex, "J", util.Vector2{X: 109, Y: 26}, util.Red},
		{FontDex, "K", util.Vector2{X: 109 + 20, Y: 50}, util.Red},
		{FontDex, "L", util.Vector2{X: 109 + 40, Y: 26}, util.Red},
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
		pos:      util.Vector2{X: 0, Y: 70},
		offset:   s.pos,
		size:     util.Vector2{X: util.Width, Y: util.Height - 70},
		cellSize: util.Vector2{X: 10, Y: 10},
	}

	// sine wave
	sine := &SineWave{
		colour:     util.Cyan,
		pos:        util.Vector2{X: 0, Y: util.Height * 0.75},
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
