// Package st7735 implements a driver for the ST7735 TFT displays, it comes in various screen sizes.
//
// Datasheet: https://www.crystalfontz.com/controllers/Sitronix/ST7735R/319/
package st7735 // import "tinygo.org/x/drivers/st7735"

import (
	"errors"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/pixel"
)

// Pixel formats supported by the st7735 driver.
type Color interface {
	pixel.RGB444BE | pixel.RGB565BE

	pixel.BaseColor
}

const (
	Width       = 128
	Height      = 160
	BatchLength = max(Width, Height)
)

var errOutOfBounds = errors.New("rectangle coordinates outside display area")

// Device wraps an SPI connection.
type Device struct {
	bus          drivers.SPI
	dcPin        machine.Pin
	resetPin     machine.Pin
	csPin        machine.Pin
	blPin        machine.Pin
	columnOffset int16
	rowOffset    int16
	rotation     drivers.Rotation
	batchData    pixel.Image[pixel.RGB565BE] // "image" with width, height of (batchLength, 1)
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
	x += d.rowOffset
	y += d.columnOffset

	d.Command(CASET)
	d.Tx([]uint8{uint8(x >> 8), uint8(x), uint8((x + h - 1) >> 8), uint8(x + h - 1)}, false)
	d.Command(RASET)
	d.Tx([]uint8{uint8(y >> 8), uint8(y), uint8((y + w - 1) >> 8), uint8(y + w - 1)}, false)
	d.Command(RAMWR)
}

// SetScrollWindow sets an area to scroll with fixed top and bottom parts of the display
func (d *Device) SetScrollArea(topFixedArea, bottomFixedArea int16) {
	// TODO: this code is broken, see the st7789 and ili9341 implementations for
	// how to do this correctly.
	d.Command(VSCRDEF)
	d.Tx([]uint8{
		uint8(topFixedArea >> 8), uint8(topFixedArea),
		uint8(Height - topFixedArea - bottomFixedArea>>8), uint8(Height - topFixedArea - bottomFixedArea),
		uint8(bottomFixedArea >> 8), uint8(bottomFixedArea),
	},
		false)
}

// SetScroll sets the vertical scroll address of the display.
func (d *Device) SetScroll(line int16) {
	d.Command(VSCRSADD)
	d.Tx([]uint8{uint8(line >> 8), uint8(line)}, false)
}

// StopScroll returns the display to its normal state
func (d *Device) StopScroll() {
	d.Command(NORON)
}

// FillRectangle fills a rectangle at a given coordinates with a color
func (d *Device) FillRectangle(x, y, width, height int16, c color.RGBA) error {
	k, i := d.Size()
	if x < 0 || y < 0 || width <= 0 || height <= 0 ||
		x >= k || (x+height) > k || y >= i || (y+width) > i {
		return errors.New("rectangle coordinates outside display area")
	}
	d.setWindow(x, y, width, height)

	d.batchData.FillSolidColor(pixel.NewColor[pixel.RGB565BE](c.R, c.G, c.B))
	i = width * height
	for i > 0 {
		if i >= BatchLength {
			d.Tx(d.batchData.RawBuffer(), false)
		} else {
			d.Tx(d.batchData.Rescale(int(i), 1).RawBuffer(), false)
		}
		i -= BatchLength
	}
	return nil
}

// SetPixel sets a pixel in the screen
func (d *Device) SetPixel(x int16, y int16, c color.RGBA) {
	w, h := d.Size()
	if x < 0 || y < 0 || x >= w || y >= h {
		return
	}
	d.FillRectangle(x, y, 1, 1, c)
}

// DrawBitmap copies the bitmap to the internal buffer on the screen at the
// given coordinates. It returns once the image data has been sent completely.
func (d *Device) DrawBitmap(x, y int16, bitmap pixel.Image[pixel.RGB565BE]) error {
	width, height := bitmap.Size()
	h, w := int16(width), int16(height)

	k, i := d.Size()
	if x < 0 || y < 0 || h <= 0 || w <= 0 ||
		x >= k || (x+h) > k || y >= i || (y+w) > i {
		return errOutOfBounds
	}
	d.setWindow(x, y, w, h)
	d.Tx(bitmap.RawBuffer(), false)
	return nil
}

// FillRectangle fills a rectangle at a given coordinates with a buffer
func (d *Device) FillRectangleWithBuffer(x, y, width, height int16, buffer []color.RGBA) error {
	k, l := d.Size()
	if x < 0 || y < 0 || width <= 0 || height <= 0 ||
		x >= k || (x+height) > k || y >= l || (y+width) > l {
		return errors.New("rectangle coordinates outside display area")
	}
	k = width * height
	l = int16(len(buffer))
	if k != l {
		return errors.New("buffer length does not match with rectangle size")
	}

	d.setWindow(x, y, width, height)

	offset := int16(0)
	for k > 0 {
		for i := int16(0); i < BatchLength; i++ {
			if offset+i < l {
				c := buffer[offset+i]
				d.batchData.Set(int(i), 0, pixel.NewColor[pixel.RGB565BE](c.R, c.G, c.B))
			}
		}
		if k >= BatchLength {
			d.Tx(d.batchData.RawBuffer(), false)
		} else {
			d.Tx(d.batchData.Rescale(int(k), 1).RawBuffer(), false)
		}
		k -= BatchLength
		offset += BatchLength
	}
	return nil
}

// DrawFastVLine draws a vertical line faster than using SetPixel
func (d *Device) DrawFastVLine(x, y0, y1 int16, c color.RGBA) {
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	d.FillRectangle(x, y0, y1-y0+1, 1, c)
}

// DrawFastHLine draws a horizontal line faster than using SetPixel
func (d *Device) DrawFastHLine(x0, x1, y int16, c color.RGBA) {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	d.FillRectangle(x0, y, 1, x1-x0+1, c)
}

// FillScreen fills the screen with a given color
func (d *Device) FillScreen(c color.RGBA) {
	d.FillRectangle(0, 0, Width, Height, c)
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

// Size returns the current size of the display.
func (d *Device) Size() (w, h int16) {
	return Height, Width
}

// Backlight enables or disables the backlight
func (d *Device) Backlight(on bool) {
	d.blPin.Set(on)
}

// Set the sleep mode for this LCD panel. When sleeping, the panel uses a lot
// less power. The LCD won't display an image anymore, but the memory contents
// will be kept.
func (d *Device) Sleep(sleepEnabled bool) error {
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
	return nil
}

// Invert inverts the colors of the screen (pretty instant!)
func (d *Device) Invert(invert bool) {
	if invert {
		d.Command(INVON)
	} else {
		d.Command(INVOFF)
	}
}
