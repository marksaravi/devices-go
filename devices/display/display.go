package display

import (
	"math"

	"github.com/marksaravi/devices-go/colors/rgb"
)

type pixelDevice interface {
	Update()
	Pixel(x, y int, color rgb.RGB)
	ScreenWidth() int
	ScreenHeight() int
}

type RGBDisplay interface {
	Update()
	ScreenWidth() int
	ScreenHeight() int

	// Color
	SetBackgroundColor(rgb.RGB)
	SetColor(rgb.RGB)

	// Drawing methods
	Clear()
	Pixel(x, y float64)
	Line(x1, y1, x2, y2 float64)
	// Rectangle(x1, y1, x2, y2 float64)
	// Arc(x, y, radius, startAngle, endAngle, width float64)
	Circle(x, y, radius float64)
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

func (d *rgbDevice) ScreenWidth() int {
	return d.pixeldev.ScreenWidth()
}

func (d *rgbDevice) ScreenHeight() int {
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
		dev.pixeldev.Pixel(int(math.Round(xc+d)), int(math.Round(yc+dr)), dev.color)
		dev.pixeldev.Pixel(int(math.Round(xc+d)), int(math.Round(yc-dr)), dev.color)
		dev.pixeldev.Pixel(int(math.Round(xc+dr)), int(math.Round(yc+d)), dev.color)
		dev.pixeldev.Pixel(int(math.Round(xc+dr)), int(math.Round(yc-d)), dev.color)

		dev.pixeldev.Pixel(int(math.Round(xc-d)), int(math.Round(yc+dr)), dev.color)
		dev.pixeldev.Pixel(int(math.Round(xc-d)), int(math.Round(yc-dr)), dev.color)
		dev.pixeldev.Pixel(int(math.Round(xc-dr)), int(math.Round(yc+d)), dev.color)
		dev.pixeldev.Pixel(int(math.Round(xc-dr)), int(math.Round(yc-d)), dev.color)
	}
	for dr := float64(0); dr <= radius*0.7; dr += 1 {
		// r*r = dx*dx+dy+dy
		d := math.Sqrt(radius*radius - dr*dr)
		// fmt.Println(dx, dy)
		putpixels(x, y, dr, d)
	}
}
