package display

import "github.com/marksaravi/devices-go/colors/rgb"

type pixelDisplay struct {
	pixeldev pixelDevice
	color    rgb.RGB
	bgColor  rgb.RGB
}

func (d *pixelDisplay) Pixel(x, y float64) {
	d.pixeldev.SetColor(d.color)
	d.pixeldev.Pixel(x, y)
}

func (d *pixelDisplay) Clear() {
	d.pixeldev.SetBackgroundColor(d.bgColor)
	for x := float64(0); x < d.pixeldev.ScreenWidth(); x += 1 {
		for y := float64(0); y < d.pixeldev.ScreenHeight(); y += 1 {
			d.pixeldev.Pixel(x, y)
		}
	}
}

func (d *pixelDisplay) Line(x1, y1, x2, y2 float64) {
	a := (y2 - y1) / (x2 - x1)
	b := y1 - a*x1
	for x := x1; x < x2; x += 1 {
		y := a*x + b
		d.pixeldev.Pixel(x, y)
	}
}
