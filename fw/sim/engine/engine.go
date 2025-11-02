package engine

import (
	"image"
	"image/color"
	"math"
	"time"

	"fw/util"

	"golang.org/x/exp/shiny/screen"
)

const (
	Scale  = 4
	width  = util.Width * Scale * 1.5
	height = util.Height * Scale
)

func bufferToImg(buf *util.ScreenBuffer, scale int, backlight bool) []uint8 {
	final := make([]uint8, util.Width*scale*util.Height*scale*4)
	for y := range util.Height {
		for x := range util.Width {
			r, g, b := buf[y][x].RGB()
			if !backlight {
				r >>= 3
				g >>= 3
				b >>= 3
			}

			for sy := range scale {
				for sx := range scale {
					i := ((y*scale+sy)*util.Width*scale + (x*scale + sx)) * 4
					final[i+0] = r
					final[i+1] = g
					final[i+2] = b
				}
			}
		}
	}
	return final
}

// huh, I don't think I've ever written a var with both a type and an initialiser in Go before
var bg color.Color = color.Gray{0x10}

const (
	buttonSize = 12 * Scale
	ledSize    = 4 * Scale
	pad        = 2 * Scale
)

// sprig has 8kro, my keyboard has 6kro, it's joever
// WASD IJKL
var buttonPoss = [util.ButtonsCount]util.Vector2{
	{X: buttonSize + pad, Y: buttonSize},
	{X: pad, Y: buttonSize * 2},
	{X: buttonSize + pad, Y: buttonSize * 3},
	{X: buttonSize*2 + pad, Y: buttonSize * 2},

	{X: width - buttonSize*2 - pad, Y: buttonSize},
	{X: width - buttonSize*3 - pad, Y: buttonSize * 2},
	{X: width - buttonSize*2 - pad, Y: buttonSize * 3},
	{X: width - buttonSize - pad, Y: buttonSize * 2},
}

var (
	ledLeftPos  = util.Vector2{X: pad, Y: pad}
	ledRightPos = util.Vector2{X: width - ledSize*1.5 - pad, Y: pad}
)

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
		c = color.Gray{Y: 0x80}
	}
	b.window.Fill(image.Rect(b.pos.X, b.pos.Y, b.pos.X+buttonSize, b.pos.Y+buttonSize), c, screen.Src)
}

type Buttons [util.ButtonsCount]*ButtonImpl

type Engine struct {
	backlight bool
	screen    screen.Screen
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
		backlight:   true,
		screen:      s,
		Window:      w,
		ButtonImpls: btnimpls,

		windowBuffer: wbuf,
		buffer:       &util.ScreenBuffer{},
		buttons:      btns,
	}
}

func (e *Engine) Backlight(on bool) {
	e.backlight = on
}

func (e *Engine) Buttons() util.Buttons {
	return e.buttons
}

func (e *Engine) CPUFrequency() uint32 {
	return 1_000_000_000 // bs it
}

func (e *Engine) Render() {
	img := e.windowBuffer.RGBA()
	copy(img.Pix, bufferToImg(e.buffer, Scale, e.backlight))

	e.Window.Upload(image.Point{util.Width * Scale * 0.25, 0}, e.windowBuffer, img.Bounds())
	clear(e.buffer[:])
}

// relatively alright log scale, as 0xffff brightness is brighter but not 16x as bright as 0xfff
func Log(x uint16) uint8 {
	return uint8((math.Log(float64(x)/0xffff)/math.Log(100) + 1) * 0xff)
}

func (e *Engine) SetLeft(on uint16) {
	c := bg
	if on != 0 {
		v := Log(on)
		c = color.RGBA{v, v, v, 0xff} // white
	}

	e.Window.Fill(image.Rect(ledLeftPos.X, ledLeftPos.Y, ledLeftPos.X+ledSize*1.5, ledLeftPos.Y+ledSize), c, screen.Src)
}

func (e *Engine) SetRight(on uint16) {
	c := bg
	if on != 0 {
		v := Log(on)
		c = color.RGBA{0, v, v, 0xff} // cyan
	}

	e.Window.Fill(image.Rect(ledRightPos.X, ledRightPos.Y, ledRightPos.X+ledSize*1.5, ledRightPos.Y+ledSize), c, screen.Src)
}

func (e *Engine) ScreenBuffer() *util.ScreenBuffer {
	return e.buffer
}
