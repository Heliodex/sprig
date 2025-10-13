package util

import "tinygo.org/x/drivers/pixel"

const (
	Width  = 160
	Height = 128
)

type ScreenBuffer [Height][Width]pixel.RGB565BE

func (sb *ScreenBuffer) Set(x, y int, px pixel.RGB565BE) {
	if x < 0 || x >= Width || y < 0 || y >= Height {
		return
	}
	(*sb)[y][x] = px
}

type DisplayDevice interface {
	SetScreen(buf *ScreenBuffer, i int16)
}
