package main

import (
	"machine"

	"fw/st7735"
	// "tinygo.org/x/drivers/pixel"
)

// 270MHz

const (
	rx  = machine.GP16
	sck = machine.GP18
	tx  = machine.GP19

	cs  = machine.GP20
	dc  = machine.GP22
	rst = machine.GP26
)

// MUST be declared as top-level
var screenmem = &st7735.ScreenBuffer{}

type Button struct {
	pin   machine.Pin
	prev  bool
	event func(bool)
}

func (b *Button) Pressed() bool {
	return !b.pin.Get()
}

type Engine struct {
	display           *Display
	screenmem         *st7735.ScreenBuffer
	setLeft, setRight func(uint32)
	W, A, S, D, I, J, K, L Button
}

func NewEngine() *Engine {
	display := NewDisplay()

	// init left led (GP28, PWM6 channel A)
	ledLeft := machine.PWM6
	ledLeft.Configure(machine.PWMConfig{})
	lch, _ := ledLeft.Channel(machine.GP28)
	setLeft := func(value uint32) {
		ledLeft.Set(lch, value)
	}

	// init right led (GP4, PWM2 channel A)
	ledRight := machine.PWM2
	ledRight.Configure(machine.PWMConfig{})
	rch, _ := ledRight.Channel(machine.GP4)
	setRight := func(value uint32) {
		ledRight.Set(rch, value)
	}

	e := &Engine{
		display:   display,
		screenmem: screenmem,
		setLeft:   setLeft,
		setRight:  setRight,
		W:         Button{pin: machine.GPIO5},
		A:         Button{pin: machine.GPIO6},
		S:         Button{pin: machine.GPIO7},
		D:         Button{pin: machine.GPIO8},
		I:         Button{pin: machine.GPIO12},
		J:         Button{pin: machine.GPIO13},
		K:         Button{pin: machine.GPIO14},
		L:         Button{pin: machine.GPIO15},
	}

	buttons := []Button{e.W, e.A, e.S, e.D, e.I, e.J, e.K, e.L}
	for _, b := range buttons {
		b.pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	return e
}

func (e *Engine) Render() {
	e.display.Render(e.screenmem)
}
