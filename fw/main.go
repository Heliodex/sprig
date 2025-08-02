package main

import (
	"machine"
	"time"

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

type UIElement struct {
	xPos, yPos uint8
	visible    bool
	drawTo     func(buf *st7735.ScreenBuffer)
}

// func (e UIElement) drawTo(buf *st7735.ScreenBuffer) {
// 	// buf[10000+int(e.xPos)+int(e.yPos)] = pixel.NewRGB565BE(0xff, 0x00, 0x00) // just a test

// 	// for y := range e.height {
// 	for y := range uint8(len(e.content)) {
// 		// for x := range e.width {
// 		for x := range uint8(len(e.content[y])) {
// 			pp := int(e.yPos+y)*st7735.Width + int(e.xPos+x)
// 			if pp < len(buf) {
// 				buf[pp] = e.content[y][x]
// 			}
// 		}
// 	}
// }

var buttons = map[byte]*Button{
	'W': {pin: machine.GPIO5},
	'A': {pin: machine.GPIO6},
	'S': {pin: machine.GPIO7},
	'D': {pin: machine.GPIO8},

	'I': {pin: machine.GPIO12},
	'J': {pin: machine.GPIO13},
	'K': {pin: machine.GPIO14},
	'L': {pin: machine.GPIO14},
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

	Text := func(font *Font, text string, xPos, yPos uint8, colour pixel.RGB565BE) *UIElement {
		drawTo := func(buf *st7735.ScreenBuffer) {
			xp := int(xPos)

			chars := textToChars(font, text)
			for _, char := range chars {
				for y, row := range char.content {
					for x, b := range row {
						if b == 0 {
							continue // skip empty pixels
						}

						buf[(int(yPos)+y)*width+xp+x] = colour
					}
				}

				xp += int(char.width) - 8 // unicrushed
			}
		}

		return &UIElement{
			xPos:   xPos,
			yPos:   yPos,
			drawTo: drawTo,
		}
	}

	buildText := Text(fontUnifont, "build 7", 2, 24, green)
	buildText.visible = true

	successText := Text(fontUnifont, "success", 2, 46, green)

	ui := []*UIElement{buildText, successText}

	buttons['W'].event = func(state bool) {
		if state {
			setLeft(0)
			setRight(0xffffffff)
		} else {
			setLeft(0xffffffff)
			setRight(0)
		}

		successText.visible = state
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
			if e.visible {
				e.drawTo(screenmem)
			}
		}

		d.Render(screenmem)
		clear(screenmem[:])
		time.Sleep(time.Millisecond * 100)
	}
}
