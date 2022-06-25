package display

import (
	"fmt"
	"math"

	"github.com/marksaravi/devices-go/colors"
)

type WidthType int

const (
	DEG90  = math.Pi / 2
	DEG180 = DEG90 * 2
	DEG270 = DEG90 * 3
	DEG360 = DEG90 * 4
)
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

	Arc(x, y, radius, startAngle, endAngle float64)
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

func getSectors(startAngle, endAngle float64) (int, int) {
	s1 := int(startAngle * 2 / math.Pi)
	s2 := int(endAngle * 2 / math.Pi)
	return s1, s2
}

type arcSector struct {
	sector         int
	ok             bool
	xs, xe, ys, ye float64
}

func isInsideSector0(x, y, xs, ys, xe, ye float64) bool {
	return x <= xs && x >= xe && y >= ys && y <= ye
}

func isInsideSector1(x, y, xs, ys, xe, ye float64) bool {
	return x <= xs && x >= xe && y <= ys && y >= ye
}

func isInsideSector2(x, y, xs, ys, xe, ye float64) bool {
	return x >= xs && x <= xe && y <= ys && y >= ye
}

func isInsideSector3(x, y, xs, ys, xe, ye float64) bool {
	return x >= xs && x <= xe && y >= ys && y <= ye
}

func findArcSectors(startAngle, endAngle, radius float64) []arcSector {
	sectors := make([]arcSector, 4)
	for sector := 0; sector < 4; sector++ {
		sectors[sector].sector = sector
		sectors[sector].ok = true
		sxs, sys, sxe, sye, sok := isInSector(sector, startAngle)
		exs, eys, exe, eye, eok := isInSector(sector, endAngle)
		if sok && eok {
			sectors[sector].xs = sxs
			sectors[sector].ys = sys
			sectors[sector].xe = exs
			sectors[sector].ye = eys
		} else if sok {
			sectors[sector].xs = sxs
			sectors[sector].ys = sys
			sectors[sector].xe = sxe
			sectors[sector].ye = sye
		} else if eok {
			sectors[sector].xs = exs
			sectors[sector].ys = eys
			sectors[sector].xe = exe
			sectors[sector].ye = eye
		} else {
			sectors[sector].ok = false
		}
		sectors[sector].xs *= radius
		sectors[sector].ys *= radius
		sectors[sector].xe *= radius
		sectors[sector].ye *= radius
	}
	return sectors
}

func isInSector(sector int, angle float64) (xs, ys, xe, ye float64, ok bool) {
	insector := int(math.Floor(angle * 2 / math.Pi))
	xs = math.Cos(angle)
	ys = math.Sin(angle)
	sectorEndAngle := (float64(insector) + 1) * math.Pi / 2
	xe = math.Cos(sectorEndAngle)
	ye = math.Sin(sectorEndAngle)

	ok = insector == sector
	return
}

func (dev *rgbDevice) putpixel(sector int, xc, yc, x, y float64, s arcSector) {
	tests := []func(x, y, xs, ys, xe, ye float64) bool{
		isInsideSector0,
		isInsideSector1,
		isInsideSector2,
		isInsideSector3,
	}
	if tests[sector](x, y, s.xs, s.ys, s.xe, s.ye) {
		fmt.Println("putpixel ", sector, x+xc, y+yc)
		dev.pixeldev.Pixel(int(math.Round(x+xc)), int(math.Round(y+yc)), dev.color)
	}
}

func showSectors(sectors []arcSector) {
	for i := 0; i < 4; i++ {
		s := sectors[i]
		fmt.Printf("%d(%v): xs: %5.3f, ys: %5.3f, xe: %5.3f, ye: %5.3f\n", s.sector, s.ok, s.xs, s.ys, s.xe, s.ye)
	}
}

func (dev *rgbDevice) Arc(xc, yc, radius, startAngle, endAngle float64) {
	iradius := math.Round(radius)
	sectors := findArcSectors(startAngle, endAngle, iradius)
	showSectors(sectors)
	var iradius2 = iradius * iradius
	var l1 float64 = 0
	for l1 = 0; true; l1 += 1 {
		l2 := math.Sqrt(iradius2 - l1*l1)
		if sectors[0].ok {
			dev.putpixel(0, xc, yc, l1, l2, sectors[0])
			dev.putpixel(0, xc, yc, l2, l1, sectors[0])
		}

		if sectors[1].ok {
			dev.putpixel(1, xc, yc, -l1, l2, sectors[1])
			dev.putpixel(1, xc, yc, -l2, l1, sectors[1])
		}

		if sectors[2].ok {
			dev.putpixel(2, xc, yc, -l1, -l2, sectors[2])
			dev.putpixel(2, xc, yc, -l2, -l1, sectors[2])
		}
		if sectors[3].ok {
			dev.putpixel(3, xc, yc, l1, -l2, sectors[3])
			dev.putpixel(3, xc, yc, l2, -l1, sectors[3])

		}
		if l1 >= l2 {
			break
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
