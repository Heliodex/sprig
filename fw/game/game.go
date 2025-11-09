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

type Simulation struct {
	g,
	len1, len2,
	angle1, angle2,
	mass1, mass2,
	aVel1, aVel2,
	dt float64
	// aAcc1, aAcc2 float32
}

type State [2]util.Vector2

const maxTraceLen = 30

var (
	f   int
	sim = &Simulation{
		g:     2,
		len1:  40,
		len2:  40,
		mass1: 4,
		mass2: 4,
		dt:    0.4,
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

	if btns[util.W].Pressed() {
		sim.g = 0
	} else {
		sim.g = 2
	}

	if btns[util.S].Pressed() {
		// reset
		sim = &Simulation{
			g:     2,
			len1:  40,
			len2:  40,
			mass1: 4,
			mass2: 4,
			dt:    0.4,
		}
		trace = nil
	}

	if btns[util.A].Pressed() {
		sim.aVel1 -= 0.01
	}

	if btns[util.D].Pressed() {
		sim.aVel1 += 0.01
	}

	if btns[util.I].Pressed() {
		sim.len1 = min(sim.len1+1, 100)
	}

	if btns[util.K].Pressed() {
		sim.len1 = max(sim.len1-1, 10)
	}

	if btns[util.L].Pressed() {
		sim.mass2 = min(sim.mass2+0.5, 20)
	}

	if btns[util.J].Pressed() {
		sim.mass2 = max(sim.mass2-0.5, 1)
	}

	// update simulation
	// ek := 0.5*sim.mass1*math.Square(sim.len1)*math.Square(sim.aVel1) +
	// 	0.5*sim.mass2*(math.Square(sim.len1)*math.Square(sim.aVel1)+
	// 		math.Square(sim.len2)*math.Square(sim.aVel2)+
	// 		2*sim.len1*sim.len2*sim.aVel1*sim.aVel2*math.Cos(sim.angle1-sim.angle2))

	// ep := (sim.mass1+sim.mass2)*sim.gravity*sim.len1*(1-math.Cos(sim.angle1)) +
	// 	sim.mass2*sim.gravity*sim.len2*(1-math.Cos(sim.angle2))

	// L := ek - ep

	diff := sim.angle2 - sim.angle1

	d1 := (sim.mass1+sim.mass2)*sim.len1 -
		sim.mass2*sim.len1*square(math.Cos(diff))
	d2 := (sim.len2 / sim.len1) * d1

	aAccel1 := (sim.mass2*sim.len1*square(sim.aVel1)*math.Sin(diff)*math.Cos(diff) +
		sim.mass2*sim.g*math.Sin(sim.angle2)*math.Cos(diff) +
		sim.mass2*sim.len2*square(sim.aVel2)*math.Sin(diff) -
		(sim.mass1+sim.mass2)*sim.g*math.Sin(sim.angle1)) / d1

	aAccel2 := (-sim.mass2*sim.len2*square(sim.aVel2)*math.Sin(diff)*math.Cos(diff) +
		(sim.mass1+sim.mass2)*sim.g*math.Sin(sim.angle1)*math.Cos(diff) -
		(sim.mass1+sim.mass2)*sim.len1*square(sim.aVel1)*math.Sin(diff) -
		(sim.mass1+sim.mass2)*sim.g*math.Sin(sim.angle2)) / d2

	// speed clamping, so it can't helicopter itself to NaNistan
	sim.aVel1 = min(sim.aVel1+aAccel1*sim.dt, 1)
	sim.aVel2 = min(sim.aVel2+aAccel2*sim.dt, 1)
	sim.angle1 += sim.aVel1 * sim.dt
	sim.angle2 += sim.aVel2 * sim.dt

	pendulum := &DoublePendulum{
		origin: util.V2(80, 24),
		p1: Pendulum{
			length: sim.len1,
			angle:  sim.angle1,
			mass:   sim.mass1,
		},
		p2: Pendulum{
			length: sim.len2,
			angle:  sim.angle2,
			mass:   sim.mass2,
		},
	}

	p1p := pendulum.p1.position(pendulum.origin)
	p2p := pendulum.p2.position(p1p)

	trace = append(trace, State{p1p, p2p})
	if len(trace) > maxTraceLen {
		trace = trace[1:]
	}

	for i := 1; i < len(trace); i++ {
		drawLine(en.ScreenBuffer(), trace[i-1][0], trace[i][0], util.Red)
		drawLine(en.ScreenBuffer(), trace[i-1][1], trace[i][1], util.Blue)
	}

	ui := []UIElement{pendulum}
	for _, t := range Texts {
		ui = append(ui, t)
	}

	for _, e := range ui {
		e.drawTo(en.ScreenBuffer())
	}

	en.Render()
	f++
}
