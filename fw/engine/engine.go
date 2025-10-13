package engine

import (
	"machine"

	"fw/st7735"
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

type Button struct {
	pin   machine.Pin
	prev  bool
	event func(bool)
}

func (b *Button) Pressed() bool {
	return !b.pin.Get()
}

type buttonId uint8

const (
	W buttonId = iota
	A
	S
	D
	I
	J
	K
	L
	ButtonsCount
)

// WASD IJKL
type Buttons [ButtonsCount]Button

type ScreenBuffer = st7735.ScreenBuffer

// MUST be declared as top-level (why? nobody knows)
var sb = &ScreenBuffer{}

type Engine struct {
	display *Display
	*ScreenBuffer
	SetLeft, SetRight func(uint32)
	Buttons
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

	buttons := Buttons{{pin: machine.GP5}, {pin: machine.GP6}, {pin: machine.GP7}, {pin: machine.GP8}, {pin: machine.GP12}, {pin: machine.GP13}, {pin: machine.GP14}, {pin: machine.GP15}}

	for _, b := range buttons {
		b.pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	// and fire up the engine
	return &Engine{
		display:      display,
		ScreenBuffer: sb,
		SetLeft:      setLeft,
		SetRight:     setRight,
		Buttons:      buttons,
	}
}

func (e *Engine) Render() {
	e.display.Render(e.ScreenBuffer)
	clear(e.ScreenBuffer[:])
}

// func (e *Engine) Backlight(on bool) {
// 	e.display.d.Backlight(on)
// }
