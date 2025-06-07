package main

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/st7735"
)

const (
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

type ScreenBuffer [width][height]RGB

type Display struct {
	d                  st7735.Device
	// buffer, nextBuffer *ScreenBuffer
}

func (d *Display) SetPixel(x, y int16, c RGB) {
	if x < 0 || x >= int16(width) || y < 0 || y >= int16(height) {
		return // out of bounds
	}
	d.d.SetPixel(x, y, c.RGBA())
}

func (d *Display) FillScreen(c RGB) {
	d.d.FillScreen(c.RGBA())
}

func NewDisplay() *Display {
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

	return &Display{
		d: d,
	}
}
