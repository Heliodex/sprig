package game

import (
	"fw/util"
	"math"
	"math/rand/v2"
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
		dt, scale, G float64
		origin       util.Vector2[int]
		masses       []*Mass
	}
	State []util.Vector2[float64]
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
	f   int
	sim = &Simulation{
		dt:     0.06,
		scale:  50,
		G:      1,
		origin: util.V2(64, 64),
		masses: []*Mass{
			{
				mass:     2,
				position: util.V2(-0.372008640907423, 0),
				velocity: util.V2(0, 1.21800411067968),
				colour:   util.Blue,
			},
			{
				mass:     2,
				position: util.V2[float64](1, 0),
				velocity: util.V2(0, 0.4531080538336022),
				colour:   util.Red,
			},
			{
				mass:     2,
				position: util.V2[float64](0, 0),
				velocity: util.V2(0, -(1.21800411067968 + 0.4531080538336022)),
				colour:   util.Green,
			},
		},
	}
	trace     []*State
	toggleAdd = &Toggle{
		onPress: func() {
			colours := []util.Pixel{util.Red, util.Green, util.Blue, util.White, util.Yellow, util.Cyan, util.Magenta, util.Grey2, util.Grey3}
			c := colours[len(sim.masses)%len(colours)]

			// centre on window crosshairs
			centre := util.V2[float64](util.Width/2, util.Height/2).Sub(sim.origin.Float64()).Div(sim.scale)

			newmass := &Mass{
				mass:     rand.Float64()*1.5 + 0.5,
				position: centre,
				velocity: util.V2(rand.Float64()*2-1, rand.Float64()*2-1).Mul(0.3),
				colour:   c,
			}
			sim.masses = append(sim.masses, newmass)
		},
	}
	toggleRemove = &Toggle{
		onPress: func() {
			if len(sim.masses) == 0 {
				return
			}
			sim.masses = sim.masses[:len(sim.masses)-1]
		},
	}
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

var Crosshairs = &Crosshair{util.V2(util.Width/2, util.Height/2), util.Grey3, 3}

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

	if btns[4].Pressed() {
		sim.dt += 0.001
	}

	if btns[6].Pressed() {
		sim.dt -= 0.001
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
	toggleRemove.Update(btns[5].Pressed())

	// update simulation
	forces := make([]util.Vector2[float64], len(sim.masses))

	for i, mi := range sim.masses {
		for j := i + 1; j < len(sim.masses); j++ {
			mj := sim.masses[j]

			d := mj.position.Sub(mi.position)

			force := sim.G * mi.mass * mj.mass
			distance := max(d.Len(), 1)

			f := d.Mul(force).Div(distance * distance * distance)

			forces[i] = forces[i].Add(f)
			forces[j] = forces[j].Sub(f)
		}
	}

	for i, m := range sim.masses {
		a := forces[i].Div(m.mass)

		m.velocity = m.velocity.Add(a.Mul(sim.dt))
		m.position = m.position.Add(m.velocity.Mul(sim.dt))

		sim.masses[i] = m // todo ptr or smth
	}

	state := make(State, len(sim.masses))
	for i, m := range sim.masses {
		state[i] = m.position
	}

	maxTraceLen := max(30-len(sim.masses), 0)

	trace = append(trace, &state)
	if diff := len(trace) - maxTraceLen; diff > 0 {
		trace = trace[diff:]
	}

	for i := 1; i < len(trace); i++ {
		for j := range sim.masses {
			a, b := *(trace[i-1]), *(trace[i])

			if j >= len(a) || j >= len(b) {
				continue
			}

			start := a[j].Mul(sim.scale).Add(util.V2(float64(sim.origin.X), float64(sim.origin.Y))).Int()
			end := b[j].Mul(sim.scale).Add(util.V2(float64(sim.origin.X), float64(sim.origin.Y))).Int()
			drawLine(en.ScreenBuffer(), start, end, sim.masses[j].colour)
		}
	}

	sim.drawTo(en.ScreenBuffer())

	for _, e := range Texts {
		e.drawTo(en.ScreenBuffer())
	}

	Crosshairs.drawTo(en.ScreenBuffer())

	en.Render()
	f++
}
