package engine

import (
	"image"

	"fw/util"

	"golang.org/x/exp/shiny/screen"
)

const Scale = 4

func bufferToImg(buf *util.ScreenBuffer, scale int) []uint8 {
	final := make([]uint8, util.Width*scale*util.Height*scale*4)
	for y := range util.Height {
		for x := range util.Width {
			c := buf[y][x].RGBA()

			for sy := range scale {
				for sx := range scale {
					i := ((y*scale+sy)*(util.Width*scale) + (x*scale + sx)) * 4
					final[i+0] = c.R
					final[i+1] = c.G
					final[i+2] = c.B
					final[i+3] = c.A
				}
			}
		}
	}
	return final
}

type Engine struct {
	screen screen.Screen
	screen.Window

	buffer  *util.ScreenBuffer
	buttons util.Buttons
}

func New(s screen.Screen) *Engine {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Width:  util.Width * Scale,
		Height: util.Height * Scale,
	})
	if err != nil {
		panic(err)
	}

	return &Engine{
		screen: s,
		Window: w,
		buffer: &util.ScreenBuffer{},
	}
}

func (e *Engine) Buttons() util.Buttons {
	return e.buttons
}

func (e *Engine) CPUFrequency() uint32 {
	return 1_000_000_000 // bs it
}

func (e *Engine) Render() {
	buf, err := e.screen.NewBuffer(image.Point{util.Width * Scale, util.Height * Scale})
	if err != nil {
		panic(err)
	}

	img := buf.RGBA()
	copy(img.Pix, bufferToImg(e.buffer, Scale))

	e.Window.Upload(image.Point{0, 0}, buf, img.Bounds())
}

func (e *Engine) ScreenBuffer() *util.ScreenBuffer {
	return e.buffer
}
