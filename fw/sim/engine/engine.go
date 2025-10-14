package engine

import (
	"image"
	"image/color"
	"time"

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

type ButtonImpl struct {
	State bool
}

func (b *ButtonImpl) Pressed() bool {
	return b.State
}

func (b *ButtonImpl) Set(pressed bool) {
	b.State = pressed
}

type Buttons [util.ButtonsCount]ButtonImpl

type Engine struct {
	screen screen.Screen
	screen.Window
	ButtonImpls Buttons

	windowBuffer screen.Buffer
	buffer       *util.ScreenBuffer
	buttons      util.Buttons
}

func New(s screen.Screen) *Engine {
	width, height := int(util.Width*Scale*1.5), int(util.Height*Scale)
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Width:  width,
		Height: height,
	})
	if err != nil {
		panic(err)
	}

	wbuf, err := s.NewBuffer(image.Point{util.Width * Scale, util.Height * Scale})
	if err != nil {
		panic(err)
	}

	var buttonimpls Buttons
	var btns util.Buttons
	for i := range btns {
		btns[i] = &buttonimpls[i]
	}

	time.Sleep(20 * time.Millisecond) // we need at least 4ms to initialise afaict

	w.Fill(image.Rect(0, 0, width, height), color.Gray{0x10}, screen.Src)

	return &Engine{
		screen:      s,
		Window:      w,
		ButtonImpls: buttonimpls,

		windowBuffer: wbuf,
		buffer:       &util.ScreenBuffer{},
		buttons:      btns,
	}
}

func (e *Engine) Buttons() util.Buttons {
	return e.buttons
}

func (e *Engine) CPUFrequency() uint32 {
	return 1_000_000_000 // bs it
}

func (e *Engine) Render() {
	img := e.windowBuffer.RGBA()
	copy(img.Pix, bufferToImg(e.buffer, Scale))

	e.Window.Upload(image.Point{util.Width * Scale * 0.25, 0}, e.windowBuffer, img.Bounds())
}

func (e *Engine) ScreenBuffer() *util.ScreenBuffer {
	return e.buffer
}
