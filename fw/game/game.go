package game

import (
	"fw/util"
	"math"
	"runtime"
	"strconv"
	"time"
)

type UIElement interface {
	drawTo(*util.ScreenBuffer)
}

type Cube3D struct {
	colour               util.Pixel
	pos                  util.Vector2[int]
	size, angle1, angle2 int
}

func (c *Cube3D) drawTo(buf *util.ScreenBuffer) {
	// rotating cube in integer 3D
	angle1 := float64(c.angle1) / 10
	angle2 := float64(c.angle2) / 10
	sin1, cos1 := math.Sin(angle1), math.Cos(angle1)
	sin2, cos2 := math.Sin(angle2), math.Cos(angle2)

	// 8 corners of the cube
	points := [8]util.Vector2[int]{}
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

		points[i] = util.V2(sx, sy)
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

	engineText := &Text{FontDex, "Grips Engine", util.V2(25, 5), util.Red} // I'm calling it this because it's an anagram of Sprig
	buildText := &Text{FontUnifont, strconv.Itoa(runtime.NumCPU()) + " cores", util.V2(5, 35), util.Blue}
	freqText := &Text{FontUnifont, strconv.Itoa(int(freq/1_000_000)) + "MHz", util.V2(5, 50), util.Green}

	ui := []UIElement{engineText, buildText, freqText}
	if freq >= 270_000_000 {
		ui = append(ui, &Text{FontUnifont, "OVERCLOCKED!", util.V2(5, 65), util.Red})
	} else {
		ui = append(ui, &Text{FontUnifont, "(could be better)", util.V2(5, 65), util.White})
	}

	for _, e := range ui {
		e.drawTo(en.ScreenBuffer())
	}
	en.Render()

	time.Sleep(time.Second)
	println("Splash complete")
}

var (
	f   int
	pos util.Vector2[int]
)

const (
	hOffset = 15
	vOffset = 18
)

var Texts = [util.ButtonsCount]*Text{
	{FontDex, "W", util.V2(2+hOffset-2, 2), util.Red},
	{FontDex, "A", util.V2(2, 2+vOffset), util.Red},
	{FontDex, "S", util.V2(2+hOffset, 2+vOffset*2), util.Red},
	{FontDex, "D", util.V2(2+hOffset*2, 2+vOffset), util.Red},
	{FontDex, "I", util.V2(2+util.Width-hOffset*2, 2), util.Red},
	{FontDex, "J", util.V2(2+util.Width-hOffset*3, 2+vOffset), util.Red},
	{FontDex, "K", util.V2(2+util.Width-hOffset*2, 2+vOffset*2), util.Red},
	{FontDex, "L", util.V2(2+util.Width-hOffset, 2+vOffset), util.Red},
}

// ran every frame (or, more like this is what makes the frames)
func Update(en util.Engine) {
	btns := en.Buttons()

	if btns[util.W].Pressed() {
		pos.Y--
	}

	if btns[util.S].Pressed() {
		pos.Y++
	}

	if btns[util.A].Pressed() {
		pos.X--
	}

	if btns[util.D].Pressed() {
		pos.X++
	}

	x := &Text{FontUnifont, strconv.Itoa(pos.X), util.V2(70, 2), util.Green}
	y := &Text{FontUnifont, strconv.Itoa(pos.Y), util.V2(70, 22), util.Green}

	// read button states
	for i, b := range btns {
		if b.Pressed() {
			Texts[i].colour = util.Green
		} else {
			Texts[i].colour = util.Grey4
		}
	}

	grid := &DiamondGrid{
		colour1:  util.Grey2,
		colour2:  util.Grey3,
		pos:      util.V2(0, 0),
		offset:   pos,
		size:     util.V2(util.Width, util.Height),
		cellSize: 10,
		heightMulti: 2,
	}


	ui := []UIElement{grid, x, y}
	for _, t := range Texts {
		ui = append(ui, t)
	}

	for _, e := range ui {
		e.drawTo(en.ScreenBuffer())
	}

	for _, e := range Texts {
		e.drawTo(en.ScreenBuffer())
	}

	en.Render()
	f++
}
