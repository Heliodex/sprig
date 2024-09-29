package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gotracker/playback/format/it"
	"github.com/gotracker/playback/format/it/layout"
	"github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/song"
)

var instruments = []rune{'~', '^', '-', '/'}

const bpr = 4

func main() {
	data, err := it.IT.Load("./track.it", []feature.Feature{})
	if err != nil {
		panic(err)
	}

	bpm := float64(data.GetInitialBPM())
	ms := int((60/bpm) * 1000 / bpr)

	patterns := []index.Pattern{0}

	var toPrint strings.Builder

	for _, p := range patterns {
		pattern, err := data.GetPattern(p)
		if err != nil {
			panic(err)
		}

		rows := pattern.NumRows()
		for i := range index.Row(rows) {
			toPrint.WriteString(fmt.Sprintf("\n%d", ms))
			firstNote := true

			pattern.GetRow(i).(layout.Row[period.Linear]).ForEach(func(c index.Channel, d song.ChannelData[volume.Volume]) (bool, error) {
				n := d.GetNote()
				ns := n.String()
				if ns == "C-0" {
					return c < 4, nil
				}

				note := strings.ReplaceAll(ns, "-", "")
				instrument := instruments[d.GetInstrument() - 1]

				if firstNote {
					toPrint.WriteString(": ")
				} else {
					toPrint.WriteString(" + ")
				}
				firstNote = false

				toPrint.WriteString(fmt.Sprintf("%s%c%d", note, instrument, ms))

				return c < 4, nil // 5 channelz
			})

			toPrint.WriteString(",")
		}
	}

	// write to file
	f, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString(toPrint.String())
	if err != nil {
		panic(err)
	}
}
