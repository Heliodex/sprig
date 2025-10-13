package main

import (
	"image"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"

	"fw/util"
)

func bufferToImg(buf util.ScreenBuffer, scale int) []uint8 {
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

func startUI(s screen.Screen) {
	const scale = 4

	w, err := s.NewWindow(&screen.NewWindowOptions{
		Width:  util.Width * scale,
		Height: util.Height * scale,
	})
	if err != nil {
		panic(err)
	}
	defer w.Release()

	var screenBuffer util.ScreenBuffer

	screenBuffer.Set(5, 5, util.Red)

	for {
		buf, err := s.NewBuffer(image.Point{util.Width * scale, util.Height * scale})
		if err != nil {
			panic(err)
		}

		img := buf.RGBA()
		copy(img.Pix, bufferToImg(screenBuffer, scale))

		w.Upload(image.Point{0, 0}, buf, img.Bounds())

		switch e := w.NextEvent().(type) {
		case lifecycle.Event:
			if e.To == lifecycle.StageDead {
				return
			}
		}
	}
}

func main() {
	driver.Main(startUI)
}
