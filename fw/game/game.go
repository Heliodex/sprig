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

type Line struct {
	start, end util.Vector2
	colour     util.Pixel
}

var f int

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
			Texts[i].colour = util.Red
		}
	}

	// circle := &Circle{util.V2(50, 50), 10, util.Blue}

	pendulum := &DoublePendulum{
		origin: util.V2(80, 0),
		p1: Pendulum{
			length:          50,
			angle: float32(f) * 0.01,
		},
		p2: Pendulum{
			length:          50,
			angle: float32(f) * 0.1,
		},
	}

	ui := []UIElement{}
	for _, t := range Texts {
		ui = append(ui, t)
	}
	ui = append(ui, pendulum)

	for _, e := range ui {
		e.drawTo(en.ScreenBuffer())
	}

	en.Render()
	f++
}
