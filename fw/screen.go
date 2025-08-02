package main

import (
	"image/color"
	"machine"

	"fw/st7735"

	"tinygo.org/x/drivers/pixel"
)

const (
	width  = st7735.Height // yes, really
	height = st7735.Width
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
	black = pixel.NewRGB565BE(0x00, 0x00, 0x00)
	white = pixel.NewRGB565BE(0xff, 0xff, 0xff)
	red   = pixel.NewRGB565BE(0xff, 0x00, 0x00)
	green = pixel.NewRGB565BE(0x00, 0xff, 0x00)
	blue  = pixel.NewRGB565BE(0x00, 0x00, 0xff)
)

const bpp = 2

type Display struct {
	d st7735.Device
}

func (d *Display) SetPixel(x, y int16, c RGB) {
	if x < 0 || x >= int16(width) || y < 0 || y >= int16(height) {
		return // out of bounds
	}
	d.d.SetPixel(x, y, c.RGBA())
}

// idky
const interlaced = false

func (d *Display) Render(buf *st7735.ScreenBuffer) {
	if interlaced {
		for i := int16(0); i < height; i+=2 {
			d.d.SetPixelLel(buf, i)
		}
		for i := int16(1); i < height; i+=2 {
			d.d.SetPixelLel(buf, i)
		}
	} else {
		for i := range int16(height) {
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

	return &Display{
		d: d,
	}
}
