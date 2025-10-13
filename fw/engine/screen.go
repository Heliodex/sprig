package engine

import (
	"image/color"
	"machine"

	"fw/st7735"

	"tinygo.org/x/drivers/pixel"
)

const (
	Width  = st7735.Height // yes, really
	Height = st7735.Width
)

type Pixel = pixel.RGB565BE

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

const bpp = 2

type Display struct {
	d st7735.Device
}

// idky
const interlaced = false

func (d *Display) Render(buf *st7735.ScreenBuffer) {
	if interlaced {
		for i := int16(0); i < Height; i += 2 {
			d.d.SetPixelLel(buf, i)
		}
		for i := int16(1); i < Height; i += 2 {
			d.d.SetPixelLel(buf, i)
		}
	} else {
		for i := range int16(Height) {
			d.d.SetPixelLel(buf, i)
		}
	}
}

func NewDisplay() *Display {
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 36_000_000, // 36 I think is the design limit
		SDI:       rx,
		SCK:       sck,
		SDO:       tx,
	})

	d := st7735.New(machine.SPI0, rst, dc, cs, machine.GP17)
	// d.ClearScreen()

	return &Display{
		d: d,
	}
}
