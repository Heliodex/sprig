// Package st7735 implements a driver for the ST7735 TFT displays, it comes in various screen sizes.
//
// Datasheet: https://www.crystalfontz.com/controllers/Sitronix/ST7735R/319/
package st7735 // import "tinygo.org/x/drivers/st7735"

import (
	"machine"
	"time"
	"unsafe"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/pixel"
)

const (
	// remember these fuckers are reversed from what you think they are
	Width       = 128
	Height      = 160
	BatchLength = Height
)

// Device wraps an SPI connection.
type Device struct {
	bus                           drivers.SPI
	dcPin, resetPin, csPin, blPin machine.Pin
	rotation                      drivers.Rotation
	batchData                     pixel.Image[pixel.RGB565BE] // "image" with width, height of (batchLength, 1)
}

// Configure initializes the display with default configuration
func (d *Device) Configure() {
	d.batchData = pixel.NewImage[pixel.RGB565BE](BatchLength, 1)

	// reset the device
	d.resetPin.High()
	time.Sleep(5 * time.Millisecond)
	d.resetPin.Low()
	time.Sleep(10 * time.Millisecond)
	d.resetPin.High()
	time.Sleep(10 * time.Millisecond)

	// Common initialization
	d.Command(SWRESET)
	time.Sleep(150 * time.Millisecond)
	d.Command(SLPOUT)
	time.Sleep(500 * time.Millisecond)
	d.Command(FRMCTR1)
	d.Data(0x01)
	d.Data(0x2C)
	d.Data(0x2D)
	d.Command(FRMCTR2)
	d.Data(0x01)
	d.Data(0x2C)
	d.Data(0x2D)
	d.Command(FRMCTR3)
	for range 2 {
		d.Data(0x01)
		d.Data(0x2C)
		d.Data(0x2D)
	}
	d.Command(INVCTR)
	d.Data(0x07)
	d.Command(PWCTR1)
	d.Data(0xA2)
	d.Data(0x02)
	d.Data(0x84)
	d.Command(PWCTR2)
	d.Data(0xC5)
	d.Command(PWCTR3)
	d.Data(0x0A)
	d.Data(0x00)
	d.Command(PWCTR4)
	d.Data(0x8A)
	d.Data(0x2A)
	d.Command(PWCTR5)
	d.Data(0x8A)
	d.Data(0xEE)
	d.Command(VMCTR1)
	d.Data(0x0E)

	d.Invert(false)
	d.SetRotation()

	// Set the color format depending on the generic type.
	d.Command(COLMOD)
	d.Data(0x05) // 16 bits per pixel

	// common color adjustment
	d.Command(GMCTRP1)
	d.Data(0x02)
	d.Data(0x1C)
	d.Data(0x07)
	d.Data(0x12)
	d.Data(0x37)
	d.Data(0x32)
	d.Data(0x29)
	d.Data(0x2D)
	d.Data(0x29)
	d.Data(0x25)
	d.Data(0x2B)
	d.Data(0x39)
	d.Data(0x00)
	d.Data(0x01)
	d.Data(0x03)
	d.Data(0x10)
	d.Command(GMCTRN1)
	d.Data(0x03)
	d.Data(0x1D)
	d.Data(0x07)
	d.Data(0x06)
	d.Data(0x2E)
	d.Data(0x2C)
	d.Data(0x29)
	d.Data(0x2D)
	d.Data(0x2E)
	d.Data(0x2E)
	d.Data(0x37)
	d.Data(0x3F)
	d.Data(0x00)
	d.Data(0x00)
	d.Data(0x02)
	d.Data(0x10)

	d.Command(NORON)
	time.Sleep(10 * time.Millisecond)
	d.Command(DISPON)
	time.Sleep(100 * time.Millisecond)

	d.blPin.High()
}

// New creates a new ST7735 connection. The SPI wire must already be configured.
func New(bus drivers.SPI, resetPin, dcPin, csPin, blPin machine.Pin) (d Device) {
	dcPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	resetPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	csPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	blPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	d = Device{
		bus:      bus,
		dcPin:    dcPin,
		resetPin: resetPin,
		csPin:    csPin,
		blPin:    blPin,
	}
	d.Configure()
	return
}

// setWindow prepares the screen to be modified at a given rectangle
func (d *Device) setWindow(x, y, w, h int16) {
	d.Command(CASET)
	d.Tx([]byte{byte(x >> 8), byte(x), byte((x + h - 1) >> 8), byte(x + h - 1)}, false)
	d.Command(RASET)
	d.Tx([]byte{byte(y >> 8), byte(y), byte((y + w - 1) >> 8), byte(y + w - 1)}, false)
	d.Command(RAMWR)
}

type ScreenBuffer [Width][Height]pixel.RGB565BE

const bpp = 2

func (d *Device) SetPixelLel(buf *ScreenBuffer, i int16) {
	d.setWindow(0, i, 1, Height)

	bs := unsafe.Slice((*byte)(unsafe.Pointer(&buf[i])), Height*2)

	d.Tx(bs, false) // little endian
}

// FillScreen fills the screen with a given color
func (d *Device) ClearScreen() {
	d.setWindow(0, 0, Width, Height)

	for range Width {
		d.Tx(make([]byte, Height*bpp), false) // fill with 0s
	}
}

// SetRotation changes the rotation of the device (clock-wise)
func (d *Device) SetRotation() {
	d.Command(MADCTL)
	d.Data(MADCTL_MX | MADCTL_MV) // we like it this way
}

// Command sends a command to the display
func (d *Device) Command(command uint8) {
	d.Tx([]byte{command}, true)
}

// Command sends a data to the display
func (d *Device) Data(data uint8) {
	d.Tx([]byte{data}, false)
}

// Tx sends data to the display
func (d *Device) Tx(data []byte, isCommand bool) {
	d.dcPin.Set(!isCommand)
	d.bus.Tx(data, nil)
}

// Backlight enables or disables the backlight
func (d *Device) Backlight(on bool) {
	d.blPin.Set(on)
}

// Set the sleep mode for this LCD panel. When sleeping, the panel uses a lot
// less power. The LCD won't display an image anymore, but the memory contents
// will be kept.
func (d *Device) Sleep(sleepEnabled bool) {
	if sleepEnabled {
		// Shut down LCD panel.
		d.Command(SLPIN)
		time.Sleep(5 * time.Millisecond) // 5ms required by the datasheet
	} else {
		// Turn the LCD panel back on.
		d.Command(SLPOUT)
		// The st7735 datasheet says it is necessary to wait 120ms before
		// sending another command.
		time.Sleep(120 * time.Millisecond)
	}
}

// Invert inverts the colors of the screen (pretty instant!)
func (d *Device) Invert(invert bool) {
	if invert {
		d.Command(INVON)
	} else {
		d.Command(INVOFF)
	}
}
