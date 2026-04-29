package game

import (
	"fmt"
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

func init() {
	const mul = 2
	const div = 4
	for x := range terrainX {
		for y := range terrainY {
			h := math.Sin(float64(x)/div)*mul + math.Cos(float64(y)/div)*mul + mul*2 + 1
			for z := range int(h) {
				terrain[x][y].Set(z, true)
				// terrain[x][y].Set(z, (x+y+z)%2 == 0)
			}
			// for z := range 7 {
			// 	terrain[x][y].Set(z, true)
			// }
		}
	}
}

type Movement struct {
	gridPosScaled util.Vector3[int]
	scale         int
	jumping       uint8
}

// func (m *Movement) Unscaled() util.Vector3[int] {
// 	fmt.Println(m.gridPosScaled)
// 	return m.gridPosScaled.FloorDiv(m.scale)
// }

func unscaled(v util.Vector3[int], scale int) util.Vector3[int] {
	return v.FloorDiv(scale)
}

func (m *Movement) Update() {
	floorSquare := unscaled(m.gridPosScaled.Add(util.V3(0, 0, -1)), m.scale)
	currentFloorTile := terrain.GetV3(floorSquare)

	if m.jumping > 0 {
		println("jumping", m.jumping)
		m.gridPosScaled.Z += 2
		m.jumping--
	}

	if !currentFloorTile {
		fmt.Println("no tile at", floorSquare)

		if m.jumping == 0 {
			m.gridPosScaled.Z -= 2 // experience gravity
		}
		// } else {
		// 	m.gridPosScaled.Z += 2 // experience normal force
	}
	// gridSquare := m.gridPosScaled
}

func (m *Movement) ValidX(x int) bool {
	gridSquare := unscaled(m.gridPosScaled, m.scale)
	gridSquareNew := unscaled(m.gridPosScaled.Add(util.V3(x, 0, 0)), m.scale)
	if gridSquareNew.X < 0 || gridSquareNew.X >= terrainX ||
		gridSquareNew.Y < 0 || gridSquareNew.Y >= terrainY {
		return false
	}

	currentTile := terrain.GetV3(gridSquare)
	if currentTile {
		m.gridPosScaled.Z += 2 // step up
		return true            // whatever bruh
	}

	newTile := terrain.GetV3(gridSquareNew)
	return !newTile
}

func (m *Movement) ValidY(y int) bool {
	gridSquare := unscaled(m.gridPosScaled, m.scale)
	gridSquareNew := unscaled(m.gridPosScaled.Add(util.V3(0, y, 0)), m.scale)
	if gridSquareNew.X < 0 || gridSquareNew.X >= terrainX ||
		gridSquareNew.Y < 0 || gridSquareNew.Y >= terrainY {
		return false
	}

	currentTile := terrain.GetV3(gridSquare)
	if currentTile {
		m.gridPosScaled.Z += 2 // step up
		return true            // whatever bruh
	}

	newTile := terrain.GetV3(gridSquareNew)
	return !newTile
}

func (m *Movement) Apply(x, y int) {
	if m.ValidX(x) {
		m.gridPosScaled.X += x
	}

	if m.ValidY(y) {
		m.gridPosScaled.Y += y
	}
}

func (m *Movement) ApplyZ(z int) {
	m.gridPosScaled.Z += z
}

func (m *Movement) ToScreen() util.Vector2[int] {
	sx := (m.gridPosScaled.X - m.gridPosScaled.Y - 68) * m.scale
	sy := (m.gridPosScaled.X + m.gridPosScaled.Y - 90 - m.gridPosScaled.Z*2) * m.scale / 2
	return util.V2(sx, sy).FloorDiv(m.scale)
}

// func x2screen(x int) int {}
// func y2screen(x int) int {}

const size = 10

var mvmt = &Movement{
	gridPosScaled: util.V3(0, 0, 8*size),
	scale:         size,
}

// ran every frame (or, more like this is what makes the frames)
func Update(en util.Engine) {
	mvmt.Update()

	btns := en.Buttons()

	if btns[util.W].Pressed() {
		mvmt.Apply(-2, -2)
	}

	if btns[util.S].Pressed() {
		mvmt.Apply(2, 2)
	}

	if btns[util.A].Pressed() {
		mvmt.Apply(-1, 1)
	}

	if btns[util.D].Pressed() {
		mvmt.Apply(1, -1)
	}

	if btns[util.I].Pressed() {
		mvmt.jumping = 5
	}

	// read button states
	for i, b := range btns {
		if b.Pressed() {
			Texts[i].colour = util.Green
		} else {
			Texts[i].colour = util.Grey4
		}
	}

	// terrain[0][0].Set(6, true)
	// terrain[1][1].Set(7, true)

	rectsize := util.V2(10, 20)
	rect := &Rect{
		pos:    util.V2((util.Width-rectsize.X)/2, (util.Height/2 - rectsize.Y)),
		size:   rectsize,
		colour: util.Magenta,
	}

	// szsize := terrainDiagonal - 1
	// sz := abs(szsize - f%(szsize*2))

	// rp := pos.Div(10)

	grid := &DimetricProjection{
		terrainHeight: 8,
		terrain:       terrain,

		topColour:   util.Green,
		baseColour1: util.Brown,
		baseColour2: util.Brown2,
		minFactor:   0x20,
		maxFactor:   0xff,
		offset:      mvmt.ToScreen(),
		cellSize:    size,
		cellHeight:  20,

		sprite:  rect,
		spriteZ: unscaled(mvmt.gridPosScaled, mvmt.scale).Z,
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
