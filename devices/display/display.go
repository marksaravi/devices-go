package display

import (
	"github.com/marksaravi/devices-go/colors/rgb565"
)

type Rotation int

const (
	ROTATION_0   Rotation = 0
	ROTATION_90  Rotation = 1
	ROTATION_180 Rotation = 2
	ROTATION_270 Rotation = 3
)

const ()

type GenericDisplay interface {
	Update()
	ScreenWidth() float64
	ScreenHeight() float64

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

type RGB565Display interface {
	GenericDisplay
	SetBackgroundColor(rgb565.RGB565)
	SetColor(rgb565.RGB565)
}

type rgb565graphicHardware interface {
	Update()
	SetColor(color rgb565.RGB565)
	Pixel(x, y float64)
	ScreenWidth() float64
	ScreenHeight() float64
}

type rgb565hardware struct {
	rgb565Dev rgb565graphicHardware
	color     rgb565.RGB565
	bgColor   rgb565.RGB565
}

func NewRGB565Display(rgb565Dev rgb565graphicHardware) RGB565Display {
	return &rgb565hardware{
		rgb565Dev: rgb565Dev,
		color:     rgb565.WHITE,
		bgColor:   rgb565.BLACK,
	}
}

func (d *rgb565hardware) Update() {
	d.rgb565Dev.Update()
}

func (d *rgb565hardware) ScreenWidth() float64 {
	return d.rgb565Dev.ScreenWidth()
}

func (d *rgb565hardware) ScreenHeight() float64 {
	return d.rgb565Dev.ScreenHeight()
}

func (d *rgb565hardware) SetBackgroundColor(color rgb565.RGB565) {
	d.bgColor = color
}

func (d *rgb565hardware) SetColor(color rgb565.RGB565) {
	d.color = color
}

// Drawing methods
func (d *rgb565hardware) Clear() {
	d.rgb565Dev.SetColor(d.bgColor)
	clear(d.rgb565Dev)
}
func (d *rgb565hardware) Pixel(x, y float64) {
	d.rgb565Dev.SetColor(d.color)
	d.rgb565Dev.Pixel(x, y)
}
func (d *rgb565hardware) Line(x1, y1, x2, y2 float64) {
	d.rgb565Dev.SetColor(d.color)
	line(x1, y1, x2, y2, d.rgb565Dev)
}
