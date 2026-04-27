package game

import (
	"embed"
	"encoding/hex"
	"strings"
)

// we're doing ascii for now for shits n giggles

//go:embed dex
var dirDex embed.FS

//go:embed unifont
var dirUnifont embed.FS

var brightnessMap = map[byte]uint8{
	' ': 0,
	'-': 64, // We make em slightly brighter than originally
	'=': 128,
	'#': 192,
	'@': 255,
}

type Char struct {
	width   uint8
	content [][]uint8
}

type Font struct {
	name          string
	dir           embed.FS
	crush, height uint8
	charset       [256]Char
}

var FontDex = &Font{
	name: "dex",
	dir:  dirDex,
}

var FontUnifont = &Font{
	name:  "unifont",
	dir:   dirUnifont,
	crush: 8,
}

func loadChar(font *Font, i uint8, dir embed.FS) {
	hexbyte := hex.EncodeToString([]byte{i})
	data, err := dir.ReadFile(font.name + "/" + hexbyte)
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	lines = lines[:len(lines)-1] // remove last empty line

	char := &font.charset[i]
	char.width = uint8(len(lines[0]))
	font.height = uint8(len(lines))
	char.content = make([][]uint8, font.height)

	for l, line := range lines {
		ll := len(line)
		lx := make([]uint8, ll)

		for x, c := range []byte(line) {
			if b, ok := brightnessMap[c]; ok { // valid brightness character
				lx[x] = b
			}
		}

		char.content[l] = lx
	}
}

func init() {
	// load all characters
	for i := range uint8(255) {
		loadChar(FontDex, i, dirDex)
		loadChar(FontUnifont, i, dirUnifont)
	}
}

func textToChars(font *Font, text string) (chars []Char) {
	bt := []byte(text)

	chars = make([]Char, len(bt))
	for i, b := range bt {
		chars[i] = font.charset[b]
	}

	return
}

// func main() {
// 	text := "Hello, world!"
// 	characters := textToChars("unifont", text)

// 	for _, char := range characters {
// 		for _, row := range char.content {
// 			for _, b := range row {
// 				fmt.Printf("%c", byteMap[b])
// 			}
// 			fmt.Println()
// 		}
// 	}
// }
