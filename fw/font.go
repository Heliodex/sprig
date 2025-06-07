package main

import (
	"embed"
	"encoding/hex"
	"strings"
)

// we're doing ascii for now for shits n giggles

// Dex Display time!!!!!!
//
//go:embed font/*
var fontdir embed.FS

var brightnessMap = map[byte]uint8{
	' ': 0,
	'-': 64, // We make em slightly brighter than originally
	'=': 128,
	'#': 192,
	'@': 255,
}

const textHeight = 20

type Char struct {
	width   uint8
	content [20][]uint8
}

var chars [256]Char

func init() {
	// load all characters
	for i := range uint8(255) {
		hexbyte := hex.EncodeToString([]byte{i})
		data, err := fontdir.ReadFile("font/" + hexbyte)
		if err != nil {
			continue
		}

		lines := strings.Split(string(data), "\n")
		lines = lines[:textHeight]

		if len(lines) > textHeight {
			continue // invalid character data
		}

		chars[i].width = uint8(len(lines[0]))

		for y, line := range lines {
			if y >= textHeight {
				continue // too tall
			}

			for _, char := range []byte(line) {
				b, ok := brightnessMap[char]
				if !ok {
					continue // invalid brightness character
				}

				chars[i].content[y] = append(chars[i].content[y], b)
			}
		}
	}
}

func textToChars(text string) []Char {
	bt := []byte(text)

	result := make([]Char, 0, len(bt))
	for _, b := range bt {
		result = append(result, chars[b])
	}

	return result
}

// func main() {
// 	text := "Hello, world!"
// 	characters := textToChars(text)

// 	for _, char := range characters {
// 		for _, row := range char.content {
// 			for _, b := range row {
// 				fmt.Printf("%c", byteMap[b])
// 			}
// 			fmt.Println()
// 		}
// 	}
// }
