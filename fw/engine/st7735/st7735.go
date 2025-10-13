// Package st7735 implements a driver for the ST7735 TFT displays, it comes in various screen sizes.
//
// Datasheet: https://www.crystalfontz.com/controllers/Sitronix/ST7735R/319/
package st7735 // import "tinygo.org/x/drivers/st7735"

import (
	"machine"
	"time"
	"unsafe"

	"fw/util"

	"tinygo.org/x/drivers"
)

const (
	// yes, rly
	Width       = util.Height
	Height      = util.Width
	BatchLength = Height
)

// Device wraps an SPI connection.
type Device struct {
	bus                                drivers.SPI
	dcPin, resetPin /*csPin, */, blPin machine.Pin
	rotation                           drivers.Rotation
}

// Tx sends data to the display
func (d *Device) Tx(data []byte, isCommand bool) {
	d.dcPin.Set(!isCommand)
	d.bus.Tx(data, nil)
}

// Command sends a command to the display
func (d *Device) Command(command uint8) {
	d.Tx([]byte{command}, true)
}

// Command sends a data to the display
func (d *Device) Data(data ...byte) {
	d.Tx(data, false)
}

// SetRotation changes the rotation of the device (clock-wise)
func (d *Device) SetRotation() {
	d.Command(MADCTL)
	d.Data(MADCTL_MX | MADCTL_MV) // we like it this way
}

// Configure initializes the display with default configuration
func (d *Device) Configure() {
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
	d.Data(0x01, 0x2C, 0x2D)
	d.Command(FRMCTR2)
	d.Data(0x01, 0x2C, 0x2D)
	d.Command(FRMCTR3)
	for range 2 {
		d.Data(0x01, 0x2C, 0x2D)
	}
	d.Command(INVCTR)
	d.Data(0x07)
	d.Command(PWCTR1)
	d.Data(0xA2, 0x02, 0x84)
	d.Command(PWCTR2)
	d.Data(0xC5)
	d.Command(PWCTR3)
	d.Data(0x0A, 0x00)
	d.Command(PWCTR4)
	d.Data(0x8A, 0x2A)
	d.Command(PWCTR5)
	d.Data(0x8A, 0xEE)
	d.Command(VMCTR1)
	d.Data(0x0E)

	d.Invert(false)
	d.SetRotation()

	// Set the color format depending on the generic type.
	d.Command(COLMOD)
	d.Data(0x05) // 16 bits per pixel

	// common color adjustment
	d.Command(GMCTRP1)
	d.Data(0x02, 0x1C, 0x07, 0x12, 0x37, 0x32, 0x29, 0x2D, 0x29, 0x25, 0x2B, 0x39, 0x00, 0x01, 0x03, 0x10)
	d.Command(GMCTRN1)
	d.Data(0x03, 0x1D, 0x07, 0x06, 0x2E, 0x2C, 0x29, 0x2D, 0x2E, 0x2E, 0x37, 0x3F, 0x00, 0x00, 0x02, 0x10)
	d.Command(NORON)
	time.Sleep(10 * time.Millisecond)
	d.Command(DISPON)
	time.Sleep(100 * time.Millisecond)

	d.blPin.High()
}

// New creates a new ST7735 connection. The SPI wire must already be configured.
func New(bus drivers.SPI, resetPin, dcPin, csPin, blPin machine.Pin) Device {
	pc := machine.PinConfig{Mode: machine.PinOutput}
	dcPin.Configure(pc)
	resetPin.Configure(pc)
	csPin.Configure(pc)
	blPin.Configure(pc)

	d := Device{
		bus:      bus,
		dcPin:    dcPin,
		resetPin: resetPin,
		// csPin:    csPin,
		blPin: blPin,
	}
	d.Configure()
	return d
}

// setWindow prepares the screen to be modified at a given rectangle
func (d *Device) setWindow(y, h int16) {
	d.Command(CASET)
	d.Data(0, 0, byte((h-1)>>8), byte(h-1))
	d.Command(RASET)
	d.Data(byte(y>>8), byte(y), byte((y)>>8), byte(y))
	d.Command(RAMWR)
}

const bpp = 2

func (d *Device) SetScreen(buf *util.ScreenBuffer, i int16) {
	d.setWindow(i, Height)

	bs := unsafe.Slice((*byte)(unsafe.Pointer(&buf[i])), Height*2)

	d.Data(bs...) // little endian
}

// Backlight enables or disables the backlight
func (d *Device) Backlight(on bool) {
	d.blPin.Set(on)
}

// Set the sleep mode for this LCD panel. When sleeping, the panel uses a lot less power. The LCD won't display an image anymore, but the memory contents will be kept.
func (d *Device) Sleep(enable bool) {
	if enable {
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
