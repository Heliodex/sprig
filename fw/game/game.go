package game

import (
	"runtime"
	"strconv"
	"time"

	"fw/util"
)

type UIElement interface {
	drawTo(*util.ScreenBuffer)
}

// intro to show, in the event that something else is loading or to show information
func Splash(en util.Engine) {
	freq := en.CPUFrequency()

	engineText := &Text{FontDex, "Grips Engine", util.NewV2(25, 5), util.Red} // I'm calling it this because it's an anagram of Sprig
	buildText := &Text{FontUnifont, strconv.Itoa(runtime.NumCPU()) + " cores", util.NewV2(5, 35), util.Blue}
	freqText := &Text{FontUnifont, strconv.Itoa(int(freq/1_000_000)) + "MHz", util.NewV2(5, 50), util.Green}

	ui := []UIElement{engineText, buildText, freqText}
	if freq >= 270_000_000 {
		ui = append(ui, &Text{FontUnifont, "OVERCLOCKED!", util.NewV2(5, 65), util.Red})
	} else {
		ui = append(ui, &Text{FontUnifont, "(could be better)", util.NewV2(5, 65), util.White})
	}

	for _, e := range ui {
		e.drawTo(en.ScreenBuffer())
	}
	en.Render()

	time.Sleep(time.Second)
	println("Splash complete")
}

type Line struct {
	start, end util.Vector2
	colour     util.Pixel
}

var (
	lines []*Line
	f     int
	pos   util.Vector2
)

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

	x := &Text{FontUnifont, strconv.Itoa(pos.X), util.NewV2(70, 2), util.Green}
	y := &Text{FontUnifont, strconv.Itoa(pos.Y), util.NewV2(70, 22), util.Green}

	Texts := [util.ButtonsCount]*Text{
		{FontDex, "W", util.NewV2(2+20-2, 2), util.Red},
		{FontDex, "A", util.NewV2(2, 26), util.Red},
		{FontDex, "S", util.NewV2(2+20, 50), util.Red},
		{FontDex, "D", util.NewV2(2+40, 26), util.Red},
		{FontDex, "I", util.NewV2(109+20+1, 2), util.Red},
		{FontDex, "J", util.NewV2(109, 26), util.Red},
		{FontDex, "K", util.NewV2(109+20, 50), util.Red},
		{FontDex, "L", util.NewV2(109+40, 26), util.Red},
	}

	// read button states
	// var leftPressed, rightPressed bool
	for i, b := range btns {
		if b.Pressed() {
			// if i < int(util.ButtonsCount/2) {
			// 	leftPressed = true
			// } else {
			// 	rightPressed = true
			// }
			Texts[i].colour = util.Green
		} else {
			Texts[i].colour = util.Red
		}
	}

	// if leftPressed {
	// 	en.SetLeft(0xffff) // max brightness
	// } else {
	// 	en.SetLeft(0xfff) // off
	// }

	// if rightPressed {
	// 	en.SetRight(0xffff) // max brightness
	// } else {
	// 	en.SetRight(0xfff) // off
	// }

	// lines = append(lines, &Line{
	// 	util.NewV2(rand.Intn(util.Width), rand.Intn(util.Height)),
	// 	util.NewV2(rand.Intn(util.Width), rand.Intn(util.Height)),
	// 	util.MakePixel(uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256))),
	// })
	// println(len(lines))

	grid := &Grid{
		colour1:  util.Grey1,
		colour2:  util.Grey2,
		pos:      util.NewV2(0, 70),
		offset:   pos,
		size:     util.NewV2(util.Width, util.Height-70),
		cellSize: util.NewV2(10, 10),
	}

	// sine wave
	sine := &SineWave{
		colour:     util.Cyan,
		pos:        util.NewV2(0, util.Height*0.75),
		amplitude:  8,
		wavelength: 40,
		phase:      f,
	}

	cube := &Cube3D{
		colour: util.Yellow,
		pos:    util.NewV2(util.Width/2, util.Height/2),
		size:   18,
		angle1: f,
		angle2: f / 2,
	}

	ui := []UIElement{grid, x, y}
	for _, t := range Texts {
		ui = append(ui, t)
	}
	ui = append(ui, sine, cube)

	for _, e := range ui {
		e.drawTo(en.ScreenBuffer())
	}

	for i := 0; i < 55; i += 5 {
		drawLine(en.ScreenBuffer(), util.NewV2(0, i), util.NewV2(util.Width, util.Height), util.Red)
		drawLine(en.ScreenBuffer(), util.NewV2(i, 0), util.NewV2(util.Width/2, util.Height), util.Yellow)
		drawLine(en.ScreenBuffer(), util.NewV2(i, util.Height), util.NewV2(util.Width, i), util.Green)
		drawLine(en.ScreenBuffer(), util.NewV2(i, util.Height/2), util.NewV2(util.Width, i), util.Blue)
	}
	// for _, line := range lines {
	// 	drawLineAntialiased(en.ScreenBuffer(), line.start, line.end, line.colour)
	// }

	en.Render()
	f++
}
