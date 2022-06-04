package display

import "github.com/marksaravi/fonts-go/fonts"

type Rotation int
type RGB565 uint16

const (
	BLACK  RGB565 = 0b0
	GREEN  RGB565 = 0x003f
	BLUE   RGB565 = 0x1f << 6
	RED    RGB565 = 0x1f << 11
	WHITE  RGB565 = RED | GREEN | BLUE
	YELLOW RGB565 = RED | GREEN
)

const (
	ROTATION_0   Rotation = 0
	ROTATION_90  Rotation = 1
	ROTATION_180 Rotation = 2
	ROTATION_270 Rotation = 3
)

const ()

type Display interface {
	Update()
	ScreenWidth() int
	ScreenHeight() int

	// Drawing methods
	Clear(backgroundColor RGB565)
	Pixel(x, y float64, color RGB565)
	Line(x1, y1, x2, y2 float64, color RGB565)
	Rectangle(x1, y1, x2, y2 float64, color RGB565)
	Circle(x, y, radius float64, color RGB565)
	FillRectangle(x1, y1, x2, y2 float64, color RGB565)
	FillCircle(x, y, radius float64, color RGB565)

	// Printing methods
	MoveCursor(x, y int)
	SetFont(font fonts.BitmapFont)
	SetFontColor(color RGB565)
	SetFontBackgroundColor(color RGB565)
	SetLineHeight(height int)
	WriteChar(char byte) error
	Write(text string)
}
