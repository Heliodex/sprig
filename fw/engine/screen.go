package engine

import (
	"image/color"
	"machine"

	"fw/engine/st7735"
	"fw/util"

	"tinygo.org/x/drivers/pixel"
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
	d util.DisplayDevice
}

// idky
const interlaced = false

func (d *Display) Render(buf *util.ScreenBuffer) {
	if interlaced {
		for i := int16(0); i < util.Height; i += 2 {
			d.d.SetScreen(buf, i)
		}
		for i := int16(1); i < util.Height; i += 2 {
			d.d.SetScreen(buf, i)
		}
	} else {
		for i := range int16(util.Height) {
			d.d.SetScreen(buf, i)
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

	return &Display{
		d: &d,
	}
}
