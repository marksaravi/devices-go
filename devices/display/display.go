package display

import (
	"github.com/marksaravi/devices-go/colors/rgb"
)

type pixelDevice interface {
	Update()
	Pixel(x, y float64, color rgb.RGB)
	ScreenWidth() float64
	ScreenHeight() float64
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
}

func (d *rgbDevice) SetColor(color rgb.RGB) {
	d.color = color
}

// Drawing methods
func (d *rgbDevice) Clear() {
	for x := float64(0); x < d.pixeldev.ScreenWidth(); x += 1 {
		for y := float64(0); y < d.pixeldev.ScreenHeight(); y += 1 {
			d.pixeldev.Pixel(x, y, d.bgColor)
		}
	}
}

func (d *rgbDevice) Pixel(x, y float64) {
	d.pixeldev.Pixel(x, y, d.color)
}

func (d *rgbDevice) Line(x1, y1, x2, y2 float64) {
}
