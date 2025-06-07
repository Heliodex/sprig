package main

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/st7735"
)

const (
	rx  = machine.GP16
	sck = machine.GP18
	tx  = machine.GP19

	cs  = machine.GP20
	dc  = machine.GP22
	rst = machine.GP26

	width  = 160
	height = 128
)

func main() {
	// display things
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 36_000_000, // 36 I think is the design limit
		SDI:       rx,
		SCK:       sck,
		SDO:       tx,
	})

	d := st7735.New(machine.SPI0, rst, dc, cs, machine.GP17)
	d.Configure(st7735.Config{
		Rotation: drivers.Rotation270, // better coordinates
	})
	d.FillScreen(color.RGBA{0x00, 0x00, 0x00, 0xff})

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

	draw := func(x, y uint8, c color.RGBA) {
		if x >= width || y >= height {
			return // out of bounds
		}
		d.SetPixel(int16(x), int16(y), c)
	}

	drawText := func(text string, xPos, yPos uint8) {
		chars := textToChars(text)

		for _, char := range chars {
			if char.width == 0 {
				xPos += 4
				continue // skip empty characters
			}

			for y, row := range char.content {
				for x, b := range row {
					draw(xPos+uint8(x), yPos+uint8(y), color.RGBA{
						R: b,
						G: b,
						B: b,
						A: 0xff,
					})
				}
			}

			xPos += char.width
		}
	}

	drawText("Hello, Sprig!", 2, 2)

	setLeft(0)
	setRight(0)
	// d.FillScreen(color.RGBA{0x00, 0x00, 0x00, 0xff})
	// d.EnableBacklight(false)
}
