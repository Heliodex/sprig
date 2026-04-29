package game

import (
	"fw/util"
	"math"
	"runtime"
	"strconv"
	"time"
)

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
	f         int
	prevFrame time.Time
	pos       util.Vector2[int] = util.V2(0, 0)
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
var terrain = &Terrain{}

// ran every frame (or, more like this is what makes the frames)
func Update(en util.Engine) {
	btns := en.Buttons()

	if btns[util.W].Pressed() {
		pos.Y -= 4
	}

	if btns[util.S].Pressed() {
		pos.Y += 4
	}

	if btns[util.A].Pressed() {
		pos.X -= 4
	}

	if btns[util.D].Pressed() {
		pos.X += 4
	}

	// read button states
	for i, b := range btns {
		if b.Pressed() {
			Texts[i].colour = util.Green
		} else {
			Texts[i].colour = util.Grey4
		}
	}

	const mul = 2
	const div = 4
	for x := range terrainX {
		for y := range terrainY {
			h := math.Sin(float64(x)/div)*mul + math.Cos(float64(y)/div)*mul + mul*2 + 1
			for z := range int(h) {
				terrain[x][y].Set(z, true)
				// terrain[x][y].Set(z, (x+y+z)%2 == 0)
			}
		}
	}
	// terrain[0][0].Set(6, true)
	// terrain[1][1].Set(7, true)

	rectsize := util.V2(20, 40)
	rect := &Rect{
		pos:    util.V2((util.Width-rectsize.X)/2, (util.Height-rectsize.Y)/2),
		size:   rectsize,
		colour: util.Magenta,
	}

	size := 10
	// szsize := terrainDiagonal - 1
	// sz := abs(szsize - f%(szsize*2))

	grid := &IsometricProjection{
		terrainHeight: 8,
		terrain:       terrain,

		topColour:   util.Green,
		baseColour1: util.Brown,
		baseColour2: util.Brown2,
		minFactor:   0x20,
		maxFactor:   0xff,
		offset:      pos,
		pos:         util.V2(util.Width/2-size, util.Height/2-size),
		size:        util.V2(util.Width, util.Height),
		cellSize:    size,
		cellHeight:  20,

		sprite:  rect,
		spriteD: 25,
	}

	// fps counter
	now := time.Now()
	elapsed := now.Sub(prevFrame)
	prevFrame = now
	fps := int(1 / elapsed.Seconds())

	coords := &Text{FontUnifont, strconv.Itoa(pos.X) + " " + strconv.Itoa(pos.Y), util.V2(55, 2), util.Red}
	fpst := &Text{FontUnifont, strconv.Itoa(fps) + " FPS", util.V2(55, 2+15), util.Red}

	ui := []UIElement{grid, coords, fpst}
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
