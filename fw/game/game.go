package game

import (
	"fmt"
	"fw/util"
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
		masses []Mass
	}
	State []util.Vector2[float64]
)

const maxTraceLen = 30

var (
	f   int
	sim = &Simulation{
		masses: []Mass{
			{position: util.V2[float64](20, 24), mass: 5, colour: util.Red},
			{position: util.V2[float64](90, 34), mass: 5, colour: util.Blue},
			{position: util.V2[float64](75, 98), mass: 5, colour: util.Green},
		},
	}
	trace []State
)

func ForceBetween(m1, m2 *Mass) util.Vector2[float64] {
	dir := m2.position.Sub(m1.position)
	fmt.Println("ForceBetween:", dir)
	dist := dir.Len()
	if dist == 0 {
		return util.Vector2[float64]{}
	}

	xNormal := float64(dir.X) / float64(dist)
	yNormal := float64(dir.Y) / float64(dist)
	forceMagnitude := (m1.mass * m2.mass) / float64(dist*dist) * 10
	forceX := xNormal * forceMagnitude
	forceY := yNormal * forceMagnitude

	return util.Vector2[float64]{X: forceX, Y: forceY}
}

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
	for i := range sim.masses {
		f := util.Vector2[float64]{}
		for j := range sim.masses {
			if i != j {
				f = f.Add(ForceBetween(&sim.masses[i], &sim.masses[j]))
			}
		}
		sim.masses[i].force = f
		fmt.Println("Mass", i, "force:", sim.masses[i].force)
	}

	for i := range sim.masses {
		acc := sim.masses[i].force.Mul(1 / sim.masses[i].mass)
		sim.masses[i].position = sim.masses[i].position.Add(acc)
	}

	m1p := sim.masses[0].position
	m2p := sim.masses[1].position
	m3p := sim.masses[2].position

	trace = append(trace, State{m1p, m2p, m3p})
	if len(trace) > maxTraceLen {
		trace = trace[1:]
	}

	for i := 1; i < len(trace); i++ {
		drawLine(en.ScreenBuffer(), trace[i-1][0].Int(), trace[i][0].Int(), util.Red)
		drawLine(en.ScreenBuffer(), trace[i-1][1].Int(), trace[i][1].Int(), util.Blue)
		drawLine(en.ScreenBuffer(), trace[i-1][2].Int(), trace[i][2].Int(), util.Green)
	}

	for _, m := range sim.masses {
		m.drawTo(en.ScreenBuffer())
	}
	for _, e := range Texts {
		e.drawTo(en.ScreenBuffer())
	}

	en.Render()
	f++
}
