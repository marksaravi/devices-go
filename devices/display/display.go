package display

import (
	"github.com/marksaravi/devices-go/colors/rgb565"
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

type GenericDisplay interface {
	Update()
	ScreenWidth() int
	ScreenHeight() int

	// Drawing methods
	Clear()
	Pixel(x, y float64)
	Line(x1, y1, x2, y2 float64)
	Rectangle(x1, y1, x2, y2 float64)
	Arc(x, y, radius, startAngle, endAngle, width float64)
	Circle(x, y, radius float64)
	FillRectangle(x1, y1, x2, y2 float64)
	FillCircle(x, y, radius float64)

	// Printing methods
	MoveCursor(x, y int)
	SetFont(font fonts.BitmapFont)
	SetLineHeight(height int)
	SetLetterSpacing(spacing int)
	WriteChar(char byte) error
	Write(text string)
}

type RGB565Display interface {
	GenericDisplay
	SetBackgroundColor(rgb565.RGB565)
	SetColor(rgb565.RGB565)
}
