package util

import (
	"math"

	"tinygo.org/x/drivers/pixel"
)

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

func (sb *ScreenBuffer) ForceSetAlpha(x, y int, px Pixel, alpha uint8) {
	r1, g1, b1 := (*sb)[y][x].RGB()
	r2, g2, b2 := px.RGB()

	ua := uint16(alpha)
	nua := 0xff - ua
	r := uint8((uint16(r1)*nua + uint16(r2)*ua) / 0xff)
	g := uint8((uint16(g1)*nua + uint16(g2)*ua) / 0xff)
	b := uint8((uint16(b1)*nua + uint16(b2)*ua) / 0xff)

	(*sb)[y][x] = MakePixel(r, g, b)
}

func (sb *ScreenBuffer) SetAlpha(x, y int, px Pixel, alpha uint8) {
	if x < 0 || x >= Width || y < 0 || y >= Height {
		return
	}
	sb.ForceSetAlpha(x, y, px, alpha)
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

	Grey2 = MakePixel(0x20, 0x20, 0x20)
	Grey3 = MakePixel(0x30, 0x30, 0x30)
	Grey4 = MakePixel(0x40, 0x40, 0x40)
	Grey8 = MakePixel(0x80, 0x80, 0x80)
	GreyB = MakePixel(0xb0, 0xb0, 0xb0)

	Brown  = MakePixel(0x60, 0x30, 0x00)
	Brown2 = MakePixel(0xc0, 0x60, 0x00)
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
	SetLeft(on uint16)
	SetRight(on uint16)
	ScreenBuffer() *ScreenBuffer
}

type Num interface {
	int | float64
}

func sqrt[T Num](x T) T {
	return T(math.Sqrt(float64(x)))
}

type Vector2[T Num] struct {
	X, Y T
}

func V2[T Num](x, y T) Vector2[T] {
	return Vector2[T]{X: x, Y: y}
}

func (v Vector2[T]) Swap() Vector2[T] {
	return Vector2[T]{X: v.Y, Y: v.X}
}

func (v Vector2[T]) Add(o Vector2[T]) Vector2[T] {
	return Vector2[T]{X: v.X + o.X, Y: v.Y + o.Y}
}

func (v Vector2[T]) Sub(o Vector2[T]) Vector2[T] {
	return Vector2[T]{X: v.X - o.X, Y: v.Y - o.Y}
}

func (v Vector2[T]) Mul(scalar float64) Vector2[T] {
	return Vector2[T]{X: T(float64(v.X) * scalar), Y: T(float64(v.Y) * scalar)}
}

func (v Vector2[T]) Div(scalar float64) Vector2[T] {
	return Vector2[T]{X: T(float64(v.X) / scalar), Y: T(float64(v.Y) / scalar)}
}

func (v Vector2[T]) Neg() Vector2[T] {
	return Vector2[T]{X: -v.X, Y: -v.Y}
}

func (v Vector2[T]) Mag() float64 {
	return math.Sqrt(math.Abs(float64(v.X + v.Y)))
}

func (v Vector2[T]) Norm() Vector2[T] {
	mag := v.Mag()
	return v.Div(mag)
}

func (v Vector2[T]) Len() T {
	x2 := v.X * v.X
	y2 := v.Y * v.Y
	return sqrt(x2 + y2)
}

func (v Vector2[float64]) Int() Vector2[int] {
	return V2(int(v.X), int(v.Y))
}

func (v Vector2[int]) Float64() Vector2[float64] {
	return V2(float64(v.X), float64(v.Y))
}
