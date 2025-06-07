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

// nothing cares about the opacity anyway
type RGB struct {
	R, G, B uint8
}

func (c RGB) RGBA() color.RGBA {
	return color.RGBA{c.R, c.G, c.B, 0xff}
}

func Greyscale(v uint8) RGB {
	return RGB{v, v, v}
}

func FromRGBA(c color.RGBA) RGB {
	return RGB{
		uint8(uint16(c.R) * uint16(c.A) / 0xff),
		uint8(uint16(c.G) * uint16(c.A) / 0xff),
		uint8(uint16(c.B) * uint16(c.A) / 0xff),
	}
}

var (
	black = RGB{0x00, 0x00, 0x00}
	white = RGB{0xff, 0xff, 0xff}
	red   = RGB{0xff, 0x00, 0x00}
	green = RGB{0x00, 0xff, 0x00}
	blue  = RGB{0x00, 0x00, 0xff}
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

					draw(xPos+uint8(x), yPos+uint8(y), colour)
				}
			}

			xPos += char.width
		}
	}

	fillScreen := func(c RGB) {
		d.FillScreen(c.RGBA())
	}

	const txt = "Sprig text drawing"

	for {
		const n = 60

		for i := range n {
			drawText(txt, 0, uint8(i), Greyscale(uint8(i*0xff/n)))
		}
		for i := range n {
			drawText(txt, 0, uint8(i), Greyscale(0xff-uint8(i*0xff/n)))
		}
	}

	setLeft(0)
	setRight(0)
	fillScreen(black)
	// d.EnableBacklight(false)
}
