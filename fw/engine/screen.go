package engine

import (
	"fw/engine/st7735"
	"fw/util"
	"image/color"
	"machine"
)

type (
	Display struct {
		d util.DisplayDevice
	}

	// nothing cares about the opacity anyway
	RGB struct {
		R, G, B uint8
	}
)

const (
	bpp = 2

	// idky
	interlaced = false
)

var spi = machine.SPI0

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

func FromRGBA(c color.RGBA) RGB {
	return RGB{
		uint8(uint16(c.R) * uint16(c.A) / 0xff),
		uint8(uint16(c.G) * uint16(c.A) / 0xff),
		uint8(uint16(c.B) * uint16(c.A) / 0xff),
	}
}

func Greyscale(v uint8) RGB {
	return RGB{v, v, v}
}

func NewDisplay() *Display {
	spi.Configure(machine.SPIConfig{
		Frequency: 36_000_000, // 36 I think is the design limit
		SDI: rx,
		SCK: sck,
		SDO: tx,
	})

	d := st7735.New(spi, rst, dc, cs, machine.GP17)

	return &Display{
		d: &d,
	}
}

func (c RGB) RGBA() color.RGBA {
	return color.RGBA{c.R, c.G, c.B, 0xff}
}
