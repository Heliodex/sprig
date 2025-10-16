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

var sb = &util.ScreenBuffer{}

type Engine struct {
	buffer            *util.ScreenBuffer
	buttons           util.Buttons
	display           *Display
}

var (
	ledLeft  = machine.PWM6 // GP28, channel A
	ledRight = machine.PWM2 // GP4, channel A
)

func New() *Engine {
	display := NewDisplay()

	ledLeft.Configure(machine.PWMConfig{})
	ledRight.Configure(machine.PWMConfig{})

	machine.GP28.Configure(machine.PinConfig{Mode: machine.PinPWM})
	machine.GP4.Configure(machine.PinConfig{Mode: machine.PinPWM})

	realButtons := []Button{
		{machine.GP5},
		{machine.GP6},
		{machine.GP7},
		{machine.GP8},
		{machine.GP12},
		{machine.GP13},
		{machine.GP14},
		{machine.GP15},
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
		buffer:   sb,
		buttons:  buttons,
		display:  display,
	}
}

func (e *Engine) Backlight(on bool) {
	e.display.d.Backlight(on)
}

func (e *Engine) Buttons() util.Buttons {
	return e.buttons
}

func (e *Engine) CPUFrequency() uint32 {
	return machine.CPUFrequency()
}

func (e *Engine) Render() {
	e.display.Render(e.buffer)
	clear(e.buffer[:])
}

func (e *Engine) SetLeft(value uint16) {
	ledLeft.Set(0, uint32(value)) // this is turned to uint16 after calling though...
}

func (e *Engine) SetRight(value uint16) {
	ledRight.Set(0, uint32(value))
}

func (e *Engine) ScreenBuffer() *util.ScreenBuffer {
	return e.buffer
}
