package main

import (
	"fmt"

	"github.com/gotracker/playback/format/it"
	"github.com/gotracker/playback/format/it/layout"
	"github.com/gotracker/playback/format/it/volume"
	"github.com/gotracker/playback/index"
	"github.com/gotracker/playback/period"
	"github.com/gotracker/playback/player/feature"
	"github.com/gotracker/playback/song"
)

func main() {
	data, err := it.IT.Load("./track.it", []feature.Feature{})
	if err != nil {
		fmt.Println(err)
		return
	}

	pattern, err := data.GetPattern(0)
	if err != nil {
		fmt.Println(err)
		return
	}

	row := pattern.GetRow(0).(layout.Row[period.Linear])

	row.ForEach(func(c index.Channel, d song.ChannelData[volume.Volume]) (bool, error) {
		note := d.GetNote()

		fmt.Printf("Channel %d: Note %s\n", c, note)

		return c < 4, nil // 5 channel
	})
}
