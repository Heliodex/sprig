package util

import "tinygo.org/x/drivers/pixel"

const (
	Width  = 160
	Height = 128
)

type Pixel pixel.RGB565BE

func MakePixel(r, g, b uint8) Pixel {
	return Pixel(pixel.NewRGB565BE(r, g, b))
}

func (p Pixel) RGB() (r, g, b uint8) {
	rgb := pixel.RGB565BE(p).RGBA()
	return rgb.R, rgb.G, rgb.B
}

func (p Pixel) Brightness(factor uint8) Pixel {
	r, g, b := p.RGB()
	r = uint8(uint16(r) * uint16(factor) / 0xff)
	g = uint8(uint16(g) * uint16(factor) / 0xff)
	b = uint8(uint16(b) * uint16(factor) / 0xff)
	return MakePixel(r, g, b)
}

type ScreenBuffer [Height][Width]Pixel

func (sb *ScreenBuffer) Set(x, y int, px Pixel) {
	if x < 0 || x >= Width || y < 0 || y >= Height {
		return
	}
	(*sb)[y][x] = px
}

type DisplayDevice interface {
	Backlight(on bool)
	SetScreen(buf *ScreenBuffer, i int16)
}

var (
	Black   = MakePixel(0x00, 0x00, 0x00)
	White   = MakePixel(0xff, 0xff, 0xff)
	Red     = MakePixel(0xff, 0x00, 0x00)
	Green   = MakePixel(0x00, 0xff, 0x00)
	Blue    = MakePixel(0x00, 0x00, 0xff)
	Cyan    = MakePixel(0x00, 0xff, 0xff)
	Yellow  = MakePixel(0xff, 0xff, 0x00)
	Magenta = MakePixel(0xff, 0x00, 0xff)

	Grey1 = MakePixel(0x20, 0x20, 0x20)
	Grey2 = MakePixel(0x40, 0x40, 0x40)
)

type buttonId uint8

const (
	W buttonId = iota
	A
	S
	D
	I
	J
	K
	L
	ButtonsCount
)

type Button interface {
	Pressed() bool
}

// WASD IJKL
type Buttons [ButtonsCount]Button

type Engine interface {
	Backlight(on bool)
	Buttons() Buttons
	CPUFrequency() uint32
	Render()
	SetLeft(on uint32)
	SetRight(on uint32)
	ScreenBuffer() *ScreenBuffer
}

type Vector2 struct {
	X, Y int
}
