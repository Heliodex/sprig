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
		dt, scale, G, floor, angle, speed float64
		origin                     util.Vector2[int]
		masses                     []*Mass
	}
	State []util.Vector2[float64]
)

func (sim *Simulation) drawTo(sb *util.ScreenBuffer) {
	for _, m := range sim.masses {
		radius := int(math.Sqrt(m.mass/math.Pi) * 5)
		circle := &Circle{sim.origin.Add(m.position.Mul(sim.scale).Int()), radius, m.colour}
		circle.drawTo(sb)
	}

	origin := util.V2(sim.origin.X-int(sim.scale), sim.origin.Y)
	marker := &Crosshair{origin, util.Grey3, 3}
	marker.drawTo(sb)

	// draw 10px line in angle direction
	angleRad := sim.angle * (math.Pi / 180)
	lineEnd := util.V2(
		int(10*math.Cos(angleRad)),
		int(10*math.Sin(angleRad)),
	).Add(origin)
	drawLine(sb, origin, lineEnd, util.Yellow)

	// draw floor
	y := int(float64(sim.origin.Y) + sim.floor*sim.scale)
	// drawHLine(sb, 50, 0, 100, util.Red)
	drawHLine(sb, y, 0, util.Width, util.Grey3)
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
		floor:  1,
		angle:  -45,
		speed:  1,
		origin: util.V2(80, 64),
	}
	trace     []*State
	toggleAdd = &Toggle{
		onPress: func() {
			colours := []util.Pixel{util.Red, util.Green, util.Blue, util.White, util.Yellow, util.Cyan, util.Magenta, util.Grey2, util.Grey3}
			c := colours[len(sim.masses)%len(colours)]

			// base velocity on angle
			angleRad := sim.angle * (math.Pi / 180)
			vel := util.V2(math.Cos(angleRad), math.Sin(angleRad)).Mul(sim.speed)

			newmass := &Mass{
				mass:         1.5,
				initPosition: util.V2[float64](-1, 0),
				initVelocity: vel,
				colour:       c,
			}

			// if len(sim.masses) < 10 {
			sim.masses = append(sim.masses, newmass)
			// } else {
			// 	sim.masses = sim.masses[1:]
			// 	sim.masses = append(sim.masses, newmass)
			// }
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

	if btns[4].Pressed() {
		sim.angle -= 2
	}
	if btns[6].Pressed() {
		sim.angle += 2
	}

	toggleAdd.Update(btns[7].Pressed())

	// update projectile simulation
	for _, m := range sim.masses {
		m.time += sim.dt

		m.position = util.V2(
			m.initPosition.X+m.initVelocity.X*m.time,
			m.initPosition.Y+m.initVelocity.Y*m.time+(sim.G*m.time*m.time)/2,
		)

		m.velocity = util.V2(
			m.initVelocity.X,
			m.initVelocity.Y+sim.G*m.time,
		)
	}

	state := make(State, len(sim.masses))
	for i, m := range sim.masses {
		state[i] = m.position

		rf := sim.floor - 0.07
		if m.position.Y < rf {
			continue
		}

		// reset position and time
		m.initPosition = util.V2(
			m.position.X,
			rf,
		)
		m.time = 0
		// dampen velocity
		m.initVelocity = util.V2(
			m.velocity.X*0.7,
			-m.velocity.Y*0.7,
		)

		if math.Abs(m.initVelocity.Y) < 0.1 {
			m.initVelocity.Y = 0
		}
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

	en.Render()
	f++
}
