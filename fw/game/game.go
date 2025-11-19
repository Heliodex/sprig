package game

import (
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
	State []util.Vector2
)

const maxTraceLen = 30

var (
	f   int
	sim = &Simulation{
		masses: []Mass{
			{position: util.V2(60, 24), mass: 5, colour: util.Red},
			{position: util.V2(100, 24), mass: 5, colour: util.Blue},
			{position: util.V2(80, 48), mass: 5, colour: util.Green},
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

	m1p := sim.masses[0].position
	m2p := sim.masses[1].position
	m3p := sim.masses[2].position

	trace = append(trace, State{m1p, m2p, m3p})
	if len(trace) > maxTraceLen {
		trace = trace[1:]
	}

	for i := 1; i < len(trace); i++ {
		drawLine(en.ScreenBuffer(), trace[i-1][0], trace[i][0], util.Red)
		drawLine(en.ScreenBuffer(), trace[i-1][1], trace[i][1], util.Blue)
		drawLine(en.ScreenBuffer(), trace[i-1][2], trace[i][2], util.Green)
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
