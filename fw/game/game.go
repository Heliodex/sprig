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

type (
	Simulation struct {
		dt, scale float64
		origin    util.Vector2[int]
		masses    []*Mass
	}
	Trace []util.Vector2[float64]
)

func (sim *Simulation) drawTo(sb *util.ScreenBuffer) {
	for _, m := range sim.masses {
		radius := int(math.Sqrt(m.mass/math.Pi) * 5)
		circle := &Circle{sim.origin.Add(m.position.Mul(sim.scale).Int()), radius, m.colour}
		circle.drawTo(sb)
	}
}

type Toggle struct {
	last, state bool
	onPress     func()
}

func (t *Toggle) Update(current bool) {
	// if the value has changed since last frame
	if current && current != t.last {
		t.state = !t.state
		// and it's now pressed
		if t.onPress != nil {
			t.onPress()
		}
	}
	t.last = current
}

var (
	// f   int
	defaultOrigin = util.V2(80, 64)
	sim           = &Simulation{
		dt:     0.06,
		scale:  50,
		origin: defaultOrigin,
	}
	traces    []*Trace
	toggleAdd = &Toggle{
		onPress: func() {
			// base colour on charge
			// charge := rand.Float64()*2 - 1 // -1 to 1
			charge := float64(1 - (len(sim.masses)%2)*2)

			h := uint8((charge)*0x7f + 0x80)
			c := util.MakePixel(h, 0, -h)

			// centre on window crosshairs
			centre := util.V2[float64](util.Width/2, util.Height/2).Sub(sim.origin.Float64()).Div(sim.scale)

			newmass := &Mass{
				mass:     1.5,
				charge:   charge,
				position: centre,
				colour:   c,
			}

			sim.masses = append(sim.masses, newmass)
		},
	}
)

const (
	hOffset = 15
	vOffset = 18
)

var (
	Texts = [util.ButtonsCount]*Text{
		{FontDex, "W", util.V2(2+hOffset-2, 2), util.Red},
		{FontDex, "A", util.V2(2, 2+vOffset), util.Red},
		{FontDex, "S", util.V2(2+hOffset, 2+vOffset*2), util.Red},
		{FontDex, "D", util.V2(2+hOffset*2, 2+vOffset), util.Red},
		{FontDex, "I", util.V2(2+util.Width-hOffset*2, 2), util.Red},
		{FontDex, "J", util.V2(2+util.Width-hOffset*3, 2+vOffset), util.Red},
		{FontDex, "K", util.V2(2+util.Width-hOffset*2, 2+vOffset*2), util.Red},
		{FontDex, "L", util.V2(2+util.Width-hOffset, 2+vOffset), util.Red},
	}
	Crosshairs = &Crosshair{util.V2(util.Width/2, util.Height/2), util.Grey3, 3}
)

const K = 9_000_000_000

// ran every frame (or, more like this is what makes the frames)
func Update(en util.Engine) {
	btns := en.Buttons()

	// read button states
	for i, b := range btns {
		if b.Pressed() {
			Texts[i].colour = util.Green
		} else {
			Texts[i].colour = util.Grey2
		}
	}

	if btns[2].Pressed() {
		sim.origin = sim.origin.Add(util.V2(0, -2))
	}
	if btns[0].Pressed() {
		sim.origin = sim.origin.Add(util.V2(0, 2))
	}
	if btns[1].Pressed() {
		sim.origin = sim.origin.Add(util.V2(2, 0))
	}
	if btns[3].Pressed() {
		sim.origin = sim.origin.Add(util.V2(-2, 0))
	}

	toggleAdd.Update(btns[7].Pressed())

	//

	maxTraceLen := max(30-len(sim.masses), 0)

	for i, m := range sim.masses {
		// calculate trace
		if len(traces) <= i {
			traces = append(traces, &Trace{})
		}
		t := traces[i]

		*t = append(*t, m.position)
		if len(*t) > maxTraceLen {
			*t = (*t)[1:]
		}

		// update mass position

		// for each pair of charges
		var tf util.Vector2[float64]

		for j, m2 := range sim.masses {
			if i == j {
				continue // same mass
			}

			// apply force to the first mass
			diff := m.position.Sub(m2.position)

			dir := diff.Norm()
			println(dir.Mag())

			mag := max(diff.Mag(), 1)
			F := m.charge * m2.charge / square(mag)

			// fmt.Println(diff.Mag(), mag)
			tf = tf.Add(dir.Mul(F))
		}

		m.position = m.position.Add(tf.Mul(0.01))
	}

	for i, trace := range traces {
		for j := 1; j < len(*trace); j++ {
			a, b := (*trace)[j-1], (*trace)[j]

			start := a.Mul(sim.scale).Add(util.V2(float64(sim.origin.X), float64(sim.origin.Y))).Int()
			end := b.Mul(sim.scale).Add(util.V2(float64(sim.origin.X), float64(sim.origin.Y))).Int()
			drawLine(en.ScreenBuffer(), start, end, sim.masses[i].colour)
		}
	}

	sim.drawTo(en.ScreenBuffer())

	for _, e := range Texts {
		e.drawTo(en.ScreenBuffer())
	}

	Crosshairs.drawTo(en.ScreenBuffer())

	en.Render()
	// f++
}
