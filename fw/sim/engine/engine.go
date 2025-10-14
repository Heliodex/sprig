package engine

import (
	"image"
	"image/color"
	"time"

	"fw/util"

	"golang.org/x/exp/shiny/screen"
)

const (
	Scale  = 4
	width  = util.Width * Scale * 1.5
	height = util.Height * Scale
)

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

// huh, I don't think I've ever written a var with both a type and an initialiser in Go before
var bg color.Color = color.Gray{0x10}

const buttonSize = 50

// sprig has 8kro, my keyboard has 6kro, it's joever
// WASD IJKL
var buttonPoss = [util.ButtonsCount]util.Vector2{
	{X: buttonSize + 5, Y: 0},
	{X: 5, Y: buttonSize},
	{X: buttonSize + 5, Y: buttonSize * 2},
	{X: buttonSize*2 + 5, Y: buttonSize},

	{X: width - buttonSize*2 - 5, Y: 0},
	{X: width - buttonSize*3 - 5, Y: buttonSize},
	{X: width - buttonSize*2 - 5, Y: buttonSize * 2},
	{X: width - buttonSize - 5, Y: buttonSize},
}

type ButtonImpl struct {
	pos    util.Vector2
	window screen.Window
	state  bool
}

func (b *ButtonImpl) Pressed() bool {
	return b.state
}

func (b *ButtonImpl) Set(pressed bool) {
	b.state = pressed

	c := bg
	if pressed {
		c = color.White
	}
	b.window.Fill(image.Rect(b.pos.X, b.pos.Y, b.pos.X+50, b.pos.Y+50), c, screen.Src)
}

type Buttons [util.ButtonsCount]*ButtonImpl

type Engine struct {
	screen screen.Screen
	screen.Window
	ButtonImpls Buttons

	windowBuffer screen.Buffer
	buffer       *util.ScreenBuffer
	buttons      util.Buttons
}

func New(s screen.Screen) *Engine {
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

	var btnimpls Buttons
	var btns util.Buttons
	for i := range btns {
		bi := &ButtonImpl{
			pos:    buttonPoss[i],
			window: w,
		}
		btnimpls[i] = bi
		btns[i] = bi
	}

	time.Sleep(20 * time.Millisecond) // we need at least 4ms to initialise afaict

	w.Fill(image.Rect(0, 0, width, height), bg, screen.Src)

	return &Engine{
		screen:      s,
		Window:      w,
		ButtonImpls: btnimpls,

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
