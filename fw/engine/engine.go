package engine

import (
	"machine"

	"fw/util"
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

type Button struct{ machine.Pin }

func (b *Button) Pressed() bool {
	return !b.Get()
}

// MUST be declared as top-level (why? nobody knows)
var sb = &util.ScreenBuffer{}

type Engine struct {
	display           *Display
	buffer            *util.ScreenBuffer
	SetLeft, SetRight func(uint32)
	buttons           util.Buttons
}

func New() *Engine {
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

	realButtons := []Button{
		Button{machine.GP5},
		Button{machine.GP6},
		Button{machine.GP7},
		Button{machine.GP8},
		Button{machine.GP12},
		Button{machine.GP13},
		Button{machine.GP14},
		Button{machine.GP15},
	}

	for _, b := range realButtons {
		b.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	var buttons util.Buttons
	for i := range buttons {
		buttons[i] = &realButtons[i]
	}

	// and fire up the engine
	return &Engine{
		display:  display,
		buffer:   sb,
		SetLeft:  setLeft,
		SetRight: setRight,
		buttons:  buttons,
	}
}

func (e *Engine) Buttons() util.Buttons {
	return e.buttons
}

func (e *Engine) Render() {
	e.display.Render(e.buffer)
	clear(e.buffer[:])
}

func (e *Engine) ScreenBuffer() *util.ScreenBuffer {
	return e.buffer
}

// func (e *Engine) Backlight(on bool) {
// 	e.display.d.Backlight(on)
// }
