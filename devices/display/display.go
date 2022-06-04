package display

import (
	"github.com/marksaravi/devices-go/v1/colors/rgb565"
	"github.com/marksaravi/fonts-go/fonts"
)

type Rotation int

const (
	ROTATION_0   Rotation = 0
	ROTATION_90  Rotation = 1
	ROTATION_180 Rotation = 2
	ROTATION_270 Rotation = 3
)

const ()

type RGB565Display interface {
	Update()
	ScreenWidth() int
	ScreenHeight() int

	// Drawing methods
	Clear(backgroundColor rgb565.RGB565)
	Pixel(x, y float64, color rgb565.RGB565)
	Line(x1, y1, x2, y2 float64, color rgb565.RGB565)
	Rectangle(x1, y1, x2, y2 float64, color rgb565.RGB565)
	Circle(x, y, radius float64, color rgb565.RGB565)
	FillRectangle(x1, y1, x2, y2 float64, color rgb565.RGB565)
	FillCircle(x, y, radius float64, color rgb565.RGB565)

	// Printing methods
	MoveCursor(x, y int)
	SetFont(font fonts.BitmapFont)
	SetFontColor(color rgb565.RGB565)
	SetFontBackgroundColor(color rgb565.RGB565)
	SetLineHeight(height int)
	WriteChar(char byte) error
	Write(text string)
}
