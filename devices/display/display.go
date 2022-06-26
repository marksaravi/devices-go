package display

import (
	"fmt"
	"math"

	"github.com/marksaravi/devices-go/colors"
	"github.com/marksaravi/devices-go/utils"
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
	ThickArc(x, y, radius, startAngle, endAngle float64, width int, widthType WidthType)

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
	ok             bool
	xs, xe, ys, ye float64
}

func isInsideSector0(x, y, xs, ys, xe, ye float64) bool {
	return x <= xs && x > xe && y >= ys && y < ye
}

func isInsideSector1(x, y, xs, ys, xe, ye float64) bool {
	return x <= xs && x > xe && y <= ys && y > ye
}

func isInsideSector2(x, y, xs, ys, xe, ye float64) bool {
	return x >= xs && x < xe && y <= ys && y > ye
}

func isInsideSector3(x, y, xs, ys, xe, ye float64) bool {
	// fmt.Printf("%6.2f, %6.2f,%6.2f, %6.2f,%6.2f, %6.2f, %v\n", x, xs, xe, y, ys, ye, x <= xs && x > xe && y >= ys && y < ye)
	return x >= xs && x < xe && y >= ys && y < ye
}

func findArcSectors(startAngle, endAngle, radius float64) map[int][]arcSector {
	var sectorsmap map[int][]arcSector = map[int][]arcSector{
		0: make([]arcSector, 0),
		1: make([]arcSector, 0),
		2: make([]arcSector, 0),
		3: make([]arcSector, 0),
	}
	PI2 := math.Pi / 2
	from := math.Mod(startAngle, math.Pi*2)
	to := math.Mod(endAngle, math.Pi*2)
	fmt.Println("from: ", utils.ToDeg(from), " ,to: ", utils.ToDeg(to))
	if to < from {
		to += math.Pi * 2
	}
	angle := from
	for sec := 0; angle < to; sec++ {
		sector := arcSector{
			ok: false,
			xs: 0,
			ys: 0,
			xe: 0,
			ye: 0,
		}
		s := float64(sec) * PI2
		e := float64(sec+1) * PI2
		if e >= to {
			e = to
		}
		if angle >= s && angle < e {
			sector.ok = true
			sector.xs = radius * math.Cos(angle)
			sector.ys = radius * math.Sin(angle)
			sector.xe = radius * math.Cos(e)
			sector.ye = radius * math.Sin(e)
			angle = e
			sectorsmap[sec%4] = append(sectorsmap[sec%4], sector)
		}
	}
	showSectors(sectorsmap)
	return sectorsmap
}

func isInSector(sector int, fromAngle, toAngle float64) (xs, ys, xe, ye float64, ok bool) {
	insector := int(math.Floor(fromAngle * 2 / math.Pi))
	xs = math.Cos(fromAngle)
	ys = math.Sin(fromAngle)
	xe = math.Cos(toAngle)
	ye = math.Sin(toAngle)

	ok = insector == sector
	return
}

func (dev *rgbDevice) arcPutPixel(sector int, xc, yc, x, y float64, s arcSector) {
	tests := []func(x, y, xs, ys, xe, ye float64) bool{
		isInsideSector0,
		isInsideSector1,
		isInsideSector2,
		isInsideSector3,
	}
	if tests[sector](x, y, s.xs, s.ys, s.xe, s.ye) {
		dev.pixeldev.Pixel(int(math.Round(x+xc)), int(math.Round(y+yc)), dev.color)
	}
}

func showSectors(sectors map[int][]arcSector) {
	for sec := 0; sec < 4; sec++ {
		sectors := sectors[sec]
		for i := 0; i < len(sectors); i++ {
			s := sectors[i]
			fmt.Printf("%d(%v): xs: %5.3f, ys: %5.3f, xe: %5.3f, ye: %5.3f\n", sec, s.ok, s.xs, s.ys, s.xe, s.ye)
		}

	}
}

func (dev *rgbDevice) Arc(xc, yc, radius, startAngle, endAngle float64) {
	signs := [4][2]float64{{1, 1}, {-1, 1}, {-1, -1}, {1, -1}}
	iradius := math.Round(radius)
	sectormaps := findArcSectors(startAngle, endAngle, iradius)
	var iradius2 = iradius * iradius
	var l1 float64 = 0
	for l1 = 0; true; l1 += 1 {
		l2 := math.Sqrt(iradius2 - l1*l1)
		for sector := 0; sector < 4; sector++ {
			sectors := sectormaps[sector]
			for i := 0; i < len(sectors); i++ {
				if sectors[i].ok {
					dev.arcPutPixel(sector, xc, yc, signs[sector][0]*l1, signs[sector][1]*l2, sectors[i])
					dev.arcPutPixel(sector, xc, yc, signs[sector][0]*l2, signs[sector][1]*l1, sectors[i])
				}
			}
		}

		if l1 >= l2 {
			break
		}
	}
}

func (dev *rgbDevice) ThickArc(xc, yc, radius, startAngle, endAngle float64, width int, widthType WidthType) {
	rs := calcThicknessStart(radius, width, widthType)
	for dr := 0; dr < width; dr++ {
		dev.Arc(xc, yc, rs-float64(dr), startAngle, endAngle)
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
