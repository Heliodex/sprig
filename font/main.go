package main

import (
	"encoding/hex"
	"fmt"
	"image/png"
	"os"
	"strings"
)

const (
	charsDir = "./chars"
	height   = 24
	fifth    = 0xffff / 5
)

func main() {
	dir, err := os.ReadDir(charsDir)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll("out", 0o755)
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

		err = os.WriteFile("out/"+nhex, []byte(p.String()), 0o644)
		if err != nil {
			panic(err)
		}
	}
}
