package game

import (
	"math"
	"runtime"
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
				// skip empty pixels
				if b != 0 {
					buf.Set(xp+x, t.pos.Y+y, t.colour.Brightness(b))
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

type Cube3D struct {
	colour               util.Pixel
	pos                  util.Vector2
	size, angle1, angle2 int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func drawLine(buf *util.ScreenBuffer, p1, p2 util.Vector2, colour util.Pixel) {
	dx := abs(p2.X - p1.X)
	dy := -abs(p2.Y - p1.Y)
	sx := -1
	if p1.X < p2.X {
		sx = 1
	}
	sy := -1
	if p1.Y < p2.Y {
		sy = 1
	}
	e := dx + dy

	for {
		buf.Set(p1.X, p1.Y, colour)
		if p1.X == p2.X && p1.Y == p2.Y {
			break
		}
		e2 := 2 * e
		if e2 >= dy {
			e += dy
			p1.X += sx
		}
		if e2 <= dx {
			e += dx
			p1.Y += sy
		}
	}
}

func (c *Cube3D) drawTo(buf *util.ScreenBuffer) {
	// rotating cube in integer 3D
	angle1 := float64(c.angle1) / 10
	angle2 := float64(c.angle2) / 10
	sin1, cos1 := math.Sin(angle1), math.Cos(angle1)
	sin2, cos2 := math.Sin(angle2), math.Cos(angle2)

	// 8 corners of the cube
	points := [8]util.Vector2{}
	for i := range points {
		x := float64((i>>0)&1*2-1) * float64(c.size)
		y := float64((i>>1)&1*2-1) * float64(c.size)
		z := float64((i>>2)&1*2-1) * float64(c.size)

		// rotate around Y axis
		xz := x*cos1 - z*sin1
		z = x*sin1 + z*cos1
		x = xz

		// rotate around X axis
		yz := y*cos2 - z*sin2
		z = y*sin2 + z*cos2
		y = yz

		// project 3D to 2D (simple orthographic projection)
		f := 20.0 / (z + 40) // perspective factor
		sx := int(x*f) + c.pos.X
		sy := int(y*f) + c.pos.Y

		points[i] = util.Vector2{X: sx, Y: sy}
	}

	// draw edges
	edges := [][2]int{
		{0, 1},
		{1, 3},
		{3, 2},
		{2, 0},
		{4, 5},
		{5, 7},
		{7, 6},
		{6, 4},
		{0, 4},
		{1, 5},
		{2, 6},
		{3, 7},
	}

	for _, e := range edges {
		drawLine(buf, points[e[0]], points[e[1]], c.colour)
	}
}

// intro to show, in the event that something else is loading or to show information
func Splash(en util.Engine) {
	freq := en.CPUFrequency()

	engineText := &Text{FontDex, "Grips Engine", util.Vector2{X: 25, Y: 5}, util.Red} // I'm calling it this because it's an anagram of Sprig
	buildText := &Text{FontUnifont, strconv.Itoa(runtime.NumCPU()) + " cores", util.Vector2{X: 5, Y: 35}, util.Blue}
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
	println("Splash complete")
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
	var leftPressed, rightPressed bool
	for i, b := range btns {
		if b.Pressed() {
			if i < int(util.ButtonsCount/2) {
				leftPressed = true
			} else {
				rightPressed = true
			}
			Texts[i].colour = util.Green
		} else {
			Texts[i].colour = util.Red
		}
	}

	if leftPressed {
		en.SetLeft(0xffff) // max brightness
	} else {
		en.SetLeft(0xfff) // off
	}

	if rightPressed {
		en.SetRight(0xffff) // max brightness
	} else {
		en.SetRight(0xfff) // off
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

	cube := &Cube3D{
		colour: util.Yellow,
		pos:    util.Vector2{X: util.Width / 2, Y: util.Height / 2},
		size:   18,
		angle1: s.f,
		angle2: s.f / 2,
	}

	ui := []UIElement{grid, x, y}
	for _, t := range Texts {
		ui = append(ui, t)
	}
	ui = append(ui, sine, cube)

	for _, e := range ui {
		e.drawTo(en.ScreenBuffer())
	}

	en.Render()
	s.f++
}
