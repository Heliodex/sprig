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
	setLeft, setRight func(uint16)
}

func New() *Engine {
	display := NewDisplay()

	// init left led (GP28, PWM6 channel A)
	ledLeft := machine.PWM6
	ledLeft.Configure(machine.PWMConfig{})
	lch, _ := ledLeft.Channel(machine.GP28)
	println("Left channel", lch)
	setLeft := func(value uint16) {
		ledLeft.Set(lch, uint32(value))
	}

	// init right led (GP4, PWM2 channel A)
	ledRight := machine.PWM2
	ledRight.Configure(machine.PWMConfig{})
	rch, _ := ledRight.Channel(machine.GP4)
	println("Right channel", rch)
	setRight := func(value uint16) {
		ledRight.Set(rch, uint32(value))
	}

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
		setLeft:  setLeft,
		setRight: setRight,
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

func (e *Engine) SetLeft(on uint16) {
	e.setLeft(on)
}

func (e *Engine) SetRight(on uint16) {
	e.setRight(on)
}

func (e *Engine) ScreenBuffer() *util.ScreenBuffer {
	return e.buffer
}
