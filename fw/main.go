package main

import (
	"machine"

	"fw/st7735"

	"tinygo.org/x/drivers/pixel"
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

var screenmem = &st7735.ScreenBuffer{}

type Button struct {
	pin   machine.Pin
	prev  bool
	event func(bool)
}

type UIElement interface {
	drawTo(buf *st7735.ScreenBuffer)
}

type Text struct {
	font       *Font
	text       string
	xPos, yPos int
	colour     pixel.RGB565BE
	visible    bool
}

func (t *Text) drawTo(buf *st7735.ScreenBuffer) {
	if !t.visible {
		return
	}

	xp := t.xPos

	chars := textToChars(t.font, t.text)
	for _, char := range chars {
		for y, row := range char.content {
			for x, b := range row {
				if b == 0 {
					continue // skip empty pixels
				}

				yl := t.yPos + y
				// r := buf[yl]

				xl := xp + x

				buf[yl][xl] = t.colour
			}
		}

		xp += int(char.width) - int(t.font.crush)
	}
}

var buttons = map[byte]*Button{
	'W': {pin: machine.GPIO5},
	'A': {pin: machine.GPIO6},
	'S': {pin: machine.GPIO7},
	'D': {pin: machine.GPIO8},

	'I': {pin: machine.GPIO12},
	'J': {pin: machine.GPIO13},
	'K': {pin: machine.GPIO14},
	'L': {pin: machine.GPIO15},
}

func main() {
	// display things
	d := NewDisplay()
	d.d.ClearScreen()

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

	for _, button := range buttons {
		button.pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	buildText := &Text{fontUnifont, "build 9", 2, 2, blue, true}
	// successText := &Text{fontUnifont, "success", 2, 46, green, false}

	buttonTextsL := map[byte]*Text{
		'W': {fontDex, "W", 2 + 20 - 1, 24, red, true},
		'A': {fontDex, "A", 2, 48, red, true},
		'S': {fontDex, "S", 2 + 20, 72, red, true},
		'D': {fontDex, "D", 2 + 40, 48, red, true},
	}
	buttonTextsR := map[byte]*Text{
		'I': {fontDex, "I", width/2 + 2 + 20 + 1, 24, red, true},
		'J': {fontDex, "J", width/2 + 2, 48, red, true},
		'K': {fontDex, "K", width/2 + 2 + 20, 72, red, true},
		'L': {fontDex, "L", width/2 + 2 + 40, 48, red, true},
	}

	ui := []UIElement{buildText}
	for k, text := range buttonTextsL {
		buttons[k].event = func(state bool) {
			if state {
				text.colour = green
				setLeft(0xffffffff)
			} else {
				text.colour = red
				setLeft(0)
			}
		}
		ui = append(ui, text)
	}
	for k, text := range buttonTextsR {
		buttons[k].event = func(state bool) {
			if state {
				text.colour = green
				setRight(0xffffffff)
			} else {
				text.colour = red
				setRight(0)
			}
		}
		ui = append(ui, text)
	}

	// event loop i guess
	for {
		// read button states
		for _, button := range buttons {
			if button.event == nil {
				continue
			}
			state := !button.pin.Get() // inverted, so false = not pressed, true = pressed
			if state == button.prev {
				continue
			}
			button.prev = state
			button.event(state)
		}

		for _, e := range ui {
			e.drawTo(screenmem)
		}

		d.Render(screenmem)
		clear(screenmem[:])
	}
}
