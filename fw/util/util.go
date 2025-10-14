package util

import "tinygo.org/x/drivers/pixel"

const (
	Width  = 160
	Height = 128
)

type Pixel = pixel.RGB565BE

type ScreenBuffer [Height][Width]Pixel

func (sb *ScreenBuffer) Set(x, y int, px pixel.RGB565BE) {
	if x < 0 || x >= Width || y < 0 || y >= Height {
		return
	}
	(*sb)[y][x] = px
}

type DisplayDevice interface {
	SetScreen(buf *ScreenBuffer, i int16)
}

var (
	Black   = pixel.NewRGB565BE(0x00, 0x00, 0x00)
	White   = pixel.NewRGB565BE(0xff, 0xff, 0xff)
	Red     = pixel.NewRGB565BE(0xff, 0x00, 0x00)
	Green   = pixel.NewRGB565BE(0x00, 0xff, 0x00)
	Blue    = pixel.NewRGB565BE(0x00, 0x00, 0xff)
	Cyan    = pixel.NewRGB565BE(0x00, 0xff, 0xff)
	Yellow  = pixel.NewRGB565BE(0xff, 0xff, 0x00)
	Magenta = pixel.NewRGB565BE(0xff, 0x00, 0xff)

	Grey1 = pixel.NewRGB565BE(0x20, 0x20, 0x20)
	Grey2 = pixel.NewRGB565BE(0x40, 0x40, 0x40)
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
	Buttons() Buttons
	CPUFrequency() uint32
	Render()
	ScreenBuffer() *ScreenBuffer
}

type Vector2 struct {
	X, Y int
}
