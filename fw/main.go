package main

import (
	"machine"
	"strconv"

	"tinygo.org/x/drivers/pixel"
)

// 270MHz

const (
	rx  = machine.GP16
	sck = machine.GP18
	tx  = machine.GP19

	cs  = machine.GP20
	dc  = machine.GP22
	rst = machine.GP26
)

// type MegaBuffer struct {
// 	bs [][]byte
// }

// func NewMegaBuffer() MegaBuffer {
// 	return MegaBuffer{
// 		bs: [][]byte{{}},
// 	}
// }

// func (mb *MegaBuffer) append(b byte) {
// 	lb := len(mb.bs)-1
// 	if len(mb.bs[lb]) < 512 {
// 		mb.bs[lb] = append(mb.bs[lb], b)
// 	} else {
// 		mb.bs = append(mb.bs, []byte{b})
// 	}
// }

// func (mb *MegaBuffer) len() (l int) {
// 	for _, b := range mb.bs {
// 		l += len(b)
// 	}
// 	return
// }

var mem =  &ScreenBuffer2{}

func main() {
	// display things
	d := NewDisplay()
	d.FillScreen(black)

	// init left led (GP28, PWM6 channel A)
	ledLeft := machine.PWM6
	ledLeft.Configure(machine.PWMConfig{})
	lch, _ := ledLeft.Channel(machine.GP28)
	setLeft := func(value uint32) {
		ledLeft.Set(lch, value)
	}

	// init right led (GP4, PWM2 channel A)
	ledRight := machine.PWM2
	ledRight.Configure(machine.PWMConfig{})
	rch, _ := ledRight.Channel(machine.GP4)
	setRight := func(value uint32) {
		ledRight.Set(rch, value)
	}

	drawText := func(font *Font, text string, xPos, yPos uint8, colour RGB) {
		chars := textToChars(font, text)

		for _, char := range chars {
			image := pixel.NewImage[pixel.RGB565BE](int(char.width), int(font.height))
			for y, row := range char.content {
				for x, b := range row {
					if b == 0 {
						continue // skip empty pixels
					}

					image.Set(x, y, pixel.NewRGB565BE(colour.R, colour.G, colour.B))
				}
			}

			d.d.DrawBitmap(int16(xPos), int16(yPos), image)

			xPos += char.width - 8 // unicrushed
		}
	}

	// const txt = "Hello, worl!"

	for {
		drawText(fontUnifont, strconv.Itoa(len(mem)), 2, 2, red)
		d.d.FillBuffermap(mem)
	}

	// drawText(fontUnifont, "build 6", 2, 24, green)

	// size := unsafe.Sizeof(buf)
	// drawText(fontUnifont, strconv.Itoa(int(size)), 2, 2, green)

	// drawText(fontUnifont, "success", 2, 46, green)

	// for {
	// 	const n = 60

	// 	for i := range n {
	// 		drawText(txt, 0, uint8(i), RGB{uint8(i * 0xff / n), 0x00, 0x00})
	// 	}
	// 	d.FillScreen(white)
	// 	for i := range n {
	// 		drawText(txt, 0, uint8(i), RGB{0x00, 0xff - uint8(i*0xff/n), 0x00})
	// 	}
	// 	d.FillScreen(black)
	// }

	setLeft(0)
	setRight(0)
	// d.FillScreen(black)
	// d.EnableBacklight(false)
}
