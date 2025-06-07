package main

import (
	"encoding/hex"
	"fmt"
	"image/png"
	"os"
	"strings"
)

func dexDisplay() {
	const (
		height = 24
		fifth  = 0xffff / 5
		charsDir = "./dex"
	)

	dir, err := os.ReadDir(charsDir)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll("out/dex", 0o755)
	if err != nil {
		panic(err)
	}

	for _, e := range dir {
		n := e.Name()
		n = n[:len(n)-4] // remove .png extension

		file, err := os.Open(charsDir + "/" + e.Name())
		if err != nil {
			panic(err)
		}
		defer file.Close()

		fmt.Println(n)
		img, err := png.Decode(file)
		if err != nil {
			panic(err)
		}

		bounds := img.Bounds()
		width := bounds.Dx()
		if bounds.Dy() != height {
			panic("bad height")
		}

		var p strings.Builder

		for y := range height {
			if y < 3 || y > height-2 {
				continue
			}
			for x := range width {
				b, _, _, _ := img.At(x, y).RGBA()

				if b < 1*fifth {
					p.WriteByte(' ')
				} else if b < 2*fifth {
					p.WriteByte('-')
				} else if b < 3*fifth {
					p.WriteByte('=')
				} else if b < 4*fifth {
					p.WriteByte('#')
				} else {
					p.WriteByte('@')
				}
			}
			p.WriteByte('\n')
		}

		nhex := hex.EncodeToString([]byte(n))

		err = os.WriteFile("out/dex/"+nhex, []byte(p.String()), 0o644)
		if err != nil {
			panic(err)
		}
	}
}

func unifont() {
	file, err := os.Open("./uni.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = os.MkdirAll("out/unifont", 0o755)
	if err != nil {
		panic(err)
	}

	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}

	const (
		startx, starty = 52, 39
		width, height  = 16, 16
	)
	bounds := img.Bounds()
	dx, dy := bounds.Dx(), bounds.Dy()

	var i uint8
	for x := startx; x < dx; x += width * 2 {
		for y := starty; y < dy; y += height * 2 {

			var p strings.Builder

			for py := range height {
				for px := range width {
					c := img.At(x+px, y+py)
					r, _, _, _ := c.RGBA()

					if r > 0x8000 {
						p.WriteByte(' ')
						continue
					}
					p.WriteByte('@')
				}
				p.WriteByte('\n')
			}

			nhex := hex.EncodeToString([]byte{i})
			err = os.WriteFile("out/unifont/"+nhex, []byte(p.String()), 0o644)
			if err != nil {
				panic(err)
			}

			i++
		}
	}
}

func main() {
	dexDisplay()
	unifont()
}
