package display

import (
	"math"

	"github.com/marksaravi/devices-go/colors"
)

type WidthType int

const (
	INNER_WIDTH  WidthType = 0
	OUTER_WIDTH  WidthType = 1
	CENTER_WIDTH WidthType = 2
)

type pixelDevice interface {
	Update()
	Pixel(x, y int, color colors.Color)
	ScreenWidth() int
	ScreenHeight() int
}

type RGBDisplay interface {
	Update()
	ScreenWidth() int
	ScreenHeight() int

	// Color
	SetBackgroundColor(colors.Color)
	SetColor(colors.Color)

	// Drawing methods
	Clear()
	Pixel(x, y float64)
	Line(x1, y1, x2, y2 float64)

	Rectangle(x1, y1, x2, y2 float64)
	FillRectangle(x1, y1, x2, y2 float64)
	ThickRectangle(x1, y1, x2, y2 float64, width int, widthType WidthType)

	// Arc(x, y, radius, startAngle, endAngle, width float64)
	Circle(x, y, radius float64)
	ThickCircle(x, y, radius float64, width int, widthType WidthType)
	FillCircle(x, y, radius float64)

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
	color    colors.Color
	bgColor  colors.Color
}

func NewRGBDisplay(pixeldev pixelDevice) RGBDisplay {
	return &rgbDevice{
		pixeldev: pixeldev,
	}
}

func (d *rgbDevice) Update() {
	d.pixeldev.Update()
}

func (d *rgbDevice) ScreenWidth() int {
	return d.pixeldev.ScreenWidth()
}

func (d *rgbDevice) ScreenHeight() int {
	return d.pixeldev.ScreenHeight()
}

func (d *rgbDevice) SetBackgroundColor(color colors.Color) {
	d.bgColor = color
}

func (d *rgbDevice) SetColor(color colors.Color) {
	d.color = color
}

// Drawing methods
func (d *rgbDevice) Clear() {
	for x := 0; x < d.pixeldev.ScreenWidth(); x += 1 {
		for y := 0; y < d.pixeldev.ScreenHeight(); y += 1 {
			d.pixeldev.Pixel(x, y, d.bgColor)
		}
	}
}

func (d *rgbDevice) Pixel(x, y float64) {
	d.pixeldev.Pixel(int(math.Round(x)), int(math.Round(y)), d.color)
}

func (d *rgbDevice) Line(x1, y1, x2, y2 float64) {
	// Bresenham's line algorithm https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
	xs := int(math.Round(x1))
	ys := int(math.Round(y1))
	xe := int(math.Round(x2))
	ye := int(math.Round(y2))
	dx := int(math.Abs(x2 - x1))
	// sx := xs < xe ? 1 : -1
	sx := -1
	if xs < xe {
		sx = 1
	}
	dy := -int(math.Abs(y2 - y1))
	// sy := ys < ye ? 1 : -1
	sy := -1
	if ys < ye {
		sy = 1
	}
	err := dx + dy

	for true {
		d.pixeldev.Pixel(xs, ys, d.color)
		if xs == xe && ys == ye {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			if xs == xe {
				break
			}
			err = err + dy
			xs = xs + sx
		}
		if e2 <= dx {
			if ys == ye {
				break
			}
			err = err + dx
			ys = ys + sy
		}
	}
}

func (dev *rgbDevice) Circle(x, y, radius float64) {
	// Midpoint circle algorithm https://en.wikipedia.org/wiki/Midpoint_circle_algorithm
	putpixels := func(xc, yc, dr, d float64) {
		dev.Pixel(xc+d, yc+dr)
		dev.Pixel(xc+d, yc-dr)
		dev.Pixel(xc+dr, yc+d)
		dev.Pixel(xc+dr, yc-d)

		dev.Pixel(xc-d, yc+dr)
		dev.Pixel(xc-d, yc-dr)
		dev.Pixel(xc-dr, yc+d)
		dev.Pixel(xc-dr, yc-d)
	}

	var dy float64 = radius
	for dx := float64(0); dx < dy; dx += 1 {
		dy = math.Sqrt(radius*radius - dx*dx)
		putpixels(x, y, dx, dy)
	}
}

func (dev *rgbDevice) FillCircle(x, y, radius float64) {
	// Midpoint circle algorithm https://en.wikipedia.org/wiki/Midpoint_circle_algorithm
	putpixels := func(xc, yc, dr, d float64) {
		dev.Line(xc+d, yc+dr, xc-d, yc+dr)
		dev.Line(xc+d, yc-dr, xc-d, yc-dr)

		dev.Line(xc+dr, yc+d, xc-dr, yc+d)
		dev.Line(xc+dr, yc-d, xc-dr, yc-d)
	}
	for dr := float64(0); dr <= math.Ceil(radius*0.707); dr += 1 {
		d := math.Sqrt(radius*radius - dr*dr)
		putpixels(x, y, dr, d)
	}
}

func calcThicknessStart(mid float64, width int, widthType WidthType) float64 {
	from := mid
	switch widthType {
	case OUTER_WIDTH:
		from = mid + float64(width)
	case CENTER_WIDTH:
		from = mid + float64(width)/2
	}
	return from
}

func (dev *rgbDevice) ThickCircle(x, y, radius float64, width int, widthType WidthType) {
	rs := calcThicknessStart(radius, width, widthType)
	for dr := 0; dr < width; dr++ {
		dev.Circle(x, y, rs-float64(dr))
	}
}

func (dev *rgbDevice) Rectangle(x1, y1, x2, y2 float64) {
	dev.Line(x1, y1, x2, y1)
	dev.Line(x2, y1, x2, y2)
	dev.Line(x2, y2, x1, y2)
	dev.Line(x1, y2, x1, y1)
}

func (dev *rgbDevice) FillRectangle(x1, y1, x2, y2 float64) {
	l := math.Round(y2 - y1)
	dy := float64(1)
	if l < 0 {
		dy = -1
	}

	for y := float64(0); y != l; y += dy {
		dev.Line(x1, y1+y, x2, y1+y)
	}
}

func (dev *rgbDevice) ThickRectangle(x1, y1, x2, y2 float64, width int, widthType WidthType) {
	xs := x1
	xe := x2
	if x2 < x1 {
		xs = x2
		xe = x1
	}
	ys := y1
	ye := y2
	if y2 < y1 {
		ys = y2
		ye = y1
	}
	s := calcThicknessStart(0, width, widthType)
	for dxy := float64(0); dxy < float64(width); dxy++ {
		dev.Rectangle(xs-s+dxy, ys-s+dxy, xe+s-dxy, ye+s-dxy)
	}
}
