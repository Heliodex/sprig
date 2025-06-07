package main

import (
	"machine"
)

const (
	rx  = machine.GP16
	sck = machine.GP18
	tx  = machine.GP19

	cs  = machine.GP20
	dc  = machine.GP22
	rst = machine.GP26
)

func main() {
	// display things
	d := NewDisplay()

	d.FillScreen(black)

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

	draw := func(x, y uint8, c RGB) {
		if x >= width || y >= height {
			return // out of bounds
		}
		d.SetPixel(int16(x), int16(y), c)
	}

	drawText := func(text string, xPos, yPos uint8, colour RGB) {
		chars := textToChars(text)

		for _, char := range chars {
			if char.width == 0 {
				xPos += 4
				continue // skip empty characters
			}

			for y, row := range char.content {
				for x, b := range row {
					if b == 0 {
						continue // skip empty pixels
					}

					colour := colour.RGBA()
					colour.A = b

					draw(xPos+uint8(x), yPos+uint8(y), FromRGBA(colour))
				}
			}

			xPos += char.width
		}
	}

	const txt = "Sprig text drawing"

	for {
		const n = 60

		for i := range n {
			drawText(txt, 0, uint8(i), RGB{uint8(i*0xff/n), 0x00, 0x00})
		}
		for i := range n {
			drawText(txt, 0, uint8(i), RGB{0x00, 0xff - uint8(i*0xff/n), 0x00})
		}
	}

	setLeft(0)
	setRight(0)
	d.FillScreen(black)
	// d.EnableBacklight(false)
}
