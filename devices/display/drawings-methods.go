package display

type genericGraphicHardware interface {
	Update()
	Pixel(x, y float64)
	ScreenWidth() float64
	ScreenHeight() float64
}

func clear(gh genericGraphicHardware) {
	for x := float64(0); x < gh.ScreenWidth(); x += 1 {
		for y := float64(0); y < gh.ScreenHeight(); y += 1 {
			gh.Pixel(x, y)
		}
	}
}

func line(x1, y1, x2, y2 float64, gh genericGraphicHardware) {
	a := (y2 - y1) / (x2 - x1)
	b := y1 - a*x1
	for x := x1; x < x2; x += 1 {
		y := a*x + b
		gh.Pixel(x, y)
	}
}
