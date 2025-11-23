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
		dt, scale, G float64
		origin       util.Vector2[int]
		masses       []Mass
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

const maxTraceLen = 30

var (
	f   int
	sim = &Simulation{
		dt:     0.06,
		scale:  75,
		G:      1,
		origin: util.V2(64, 64),
		masses: []Mass{
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
	trace []State
)

// ran every frame (or, more like this is what makes the frames)
func Update(en util.Engine) {
	btns := en.Buttons()

	Texts := [util.ButtonsCount]*Text{
		{FontDex, "W", util.V2(2+20-2, 2), util.Red},
		{FontDex, "A", util.V2(2, 26), util.Red},
		{FontDex, "S", util.V2(2+20, 50), util.Red},
		{FontDex, "D", util.V2(2+40, 26), util.Red},
		{FontDex, "I", util.V2(109+20+1, 2), util.Red},
		{FontDex, "J", util.V2(109, 26), util.Red},
		{FontDex, "K", util.V2(109+20, 50), util.Red},
		{FontDex, "L", util.V2(109+40, 26), util.Red},
	}

	// read button states
	for i, b := range btns {
		if b.Pressed() {
			Texts[i].colour = util.Green
		} else {
			Texts[i].colour = util.Grey2
		}
	}

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

	trace = append(trace, state)
	if len(trace) > maxTraceLen {
		trace = trace[1:]
	}

	for i := 1; i < len(trace); i++ {
		for j := range sim.masses {
			drawLine(en.ScreenBuffer(), trace[i-1][j].Mul(sim.scale).Add(util.V2(float64(sim.origin.X), float64(sim.origin.Y))).Int(),
				trace[i][j].Mul(sim.scale).Add(util.V2(float64(sim.origin.X), float64(sim.origin.Y))).Int(),
				sim.masses[j].colour)
		}
	}

	sim.drawTo(en.ScreenBuffer())

	// for _, e := range Texts {
	// 	e.drawTo(en.ScreenBuffer())
	// }

	en.Render()
	f++
}
