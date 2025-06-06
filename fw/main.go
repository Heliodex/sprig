package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/st7735"
)

const (
	rx  = machine.GP16
	sck = machine.GP18
	tx  = machine.GP19

	cs  = machine.GP20
	dc  = machine.GP22
	rst = machine.GP26
)

func main() {
	// display things
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 36_000_000, // 36 I think is the design limit
		SDI:       rx,
		SCK:       sck,
		SDO:       tx,
	})

	d := st7735.New(machine.SPI0, rst, dc, cs, machine.GP17)
	d.Configure(st7735.Config{})

	// init left led (GP28, PWM6 channel A)
	ledLeft := machine.PWM6
	ledLeft.Configure(machine.PWMConfig{})
	lch, err := ledLeft.Channel(machine.GP28)
	if err != nil {
		panic(err) // where do the errors even go?
	}

	// init right led (GP4, PWM2 channel A)
	ledRight := machine.PWM2
	ledRight.Configure(machine.PWMConfig{})
	rch, err := ledRight.Channel(machine.GP4)
	if err != nil {
		panic(err)
	}

	for range 3 {
		d.FillScreen(color.RGBA{0xff, 0x00, 0x00, 0xff})
		ledLeft.Set(lch, 65535/32)
		ledRight.Set(rch, 0)
		time.Sleep(200 * time.Millisecond)

		d.FillScreen(color.RGBA{0x00, 0x00, 0x00, 0xff})
		ledLeft.Set(lch, 0)
		ledRight.Set(rch, 65535/32)
		time.Sleep(200 * time.Millisecond)
	}

	ledLeft.Set(lch, 0)
	ledRight.Set(rch, 0)
	d.FillScreen(color.RGBA{0x00, 0x00, 0x00, 0xff})
}
