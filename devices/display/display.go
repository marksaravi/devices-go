package display

import (
	"github.com/marksaravi/devices-go/colors/rgb"
)

type Rotation int

const (
	ROTATION_0   Rotation = 0
	ROTATION_90  Rotation = 1
	ROTATION_180 Rotation = 2
	ROTATION_270 Rotation = 3
)

const ()

type pixelDevice interface {
	Update()
	Pixel(x, y float64)
	ScreenWidth() float64
	ScreenHeight() float64
	SetBackgroundColor(rgb.RGB)
	SetColor(rgb.RGB)
}

type RGBDisplay interface {
	Update()
	ScreenWidth() float64
	ScreenHeight() float64

	// Color
	SetBackgroundColor(rgb.RGB)
	SetColor(rgb.RGB)

	// Drawing methods
	Clear()
	Pixel(x, y float64)
	Line(x1, y1, x2, y2 float64)
	// Rectangle(x1, y1, x2, y2 float64)
	// Arc(x, y, radius, startAngle, endAngle, width float64)
	// Circle(x, y, radius float64)
	// FillRectangle(x1, y1, x2, y2 float64)
	// FillCircle(x, y, radius float64)

	// Printing methods
	// MoveCursor(x, y int)
	// SetFont(font fonts.BitmapFont)
	// SetLineHeight(height int)
	// SetLetterSpacing(spacing int)
	// WriteChar(char byte) error
	// Write(text string)
}

type rgbDevice struct {
	pixeldev pixelDevice
	color    rgb.RGB
	bgColor  rgb.RGB
}

func NewRGBDisplay(pixeldev pixelDevice) RGBDisplay {
	return &rgbDevice{
		pixeldev: pixeldev,
	}
}

func (d *rgbDevice) Update() {
	d.pixeldev.Update()
}

func (d *rgbDevice) ScreenWidth() float64 {
	return d.pixeldev.ScreenWidth()
}

func (d *rgbDevice) ScreenHeight() float64 {
	return d.pixeldev.ScreenHeight()
}

func (d *rgbDevice) SetBackgroundColor(color rgb.RGB) {
	d.bgColor = color
	d.pixeldev.SetBackgroundColor(color)
}

func (d *rgbDevice) SetColor(color rgb.RGB) {
	d.color = color
	d.pixeldev.SetColor(color)
}

// Drawing methods
func (d *rgbDevice) Clear() {
	d.pixeldev.SetColor(d.bgColor)

}

func (d *rgbDevice) Pixel(x, y float64) {
	d.pixeldev.SetColor(d.color)
	d.pixeldev.Pixel(x, y)
}

func (d *rgbDevice) Line(x1, y1, x2, y2 float64) {
	d.pixeldev.SetColor(d.color)
}
