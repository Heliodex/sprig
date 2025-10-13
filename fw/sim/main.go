package main

import (
	"image"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"

	"fw/game"
	"fw/sim/engine"
	"fw/util"
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

func startUI(s screen.Screen) {
	en := engine.New(s)
	game.Splash(en)

	buf, err := s.NewBuffer(image.Point{util.Width * engine.Scale, util.Height * engine.Scale})
	if err != nil {
		panic(err)
	}

	img := buf.RGBA()
	copy(img.Pix, bufferToImg(en.ScreenBuffer(), engine.Scale))

	for {

		en.Window.Upload(image.Point{0, 0}, buf, img.Bounds())

		switch e := en.Window.NextEvent().(type) {
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
