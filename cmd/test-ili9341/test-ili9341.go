package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/marksaravi/devices-go/colors"
	"github.com/marksaravi/devices-go/devices/display"
	"github.com/marksaravi/devices-go/hardware/ili9341"
	"github.com/marksaravi/devices-go/utils"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/sysfs"
)

func checkFatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Testing ILI9341...")
	host.Init()
	spiConn := createSPIConnection(0, 0)
	dataCommandSelect := createGpioOutPin("GPIO22")
	reset := createGpioOutPin("GPIO23")

	ili9341Dev, err := ili9341.NewILI9341(spiConn, dataCommandSelect, reset)
	var ili9341Display display.RGBDisplay
	ili9341Display = display.NewRGBDisplay(ili9341Dev)
	checkFatalErr(err)
	tests := []func(display.RGBDisplay){
		drawLines,
		drawArc,
		draThickwArc,
		drawCircle,
		drawFillCircle,
		drawThickCircle,
		drawRectangle,
		drawFillRectangle,
		drawThickRectangle,
	}
	for i := 0; i < len(tests); i++ {
		ili9341Display.SetBackgroundColor(colors.WHITE)
		ili9341Display.Clear()
		tests[i](ili9341Display)
		ili9341Display.Update()
		time.Sleep(1000 * time.Millisecond)
	}

	// testColors(ili9341Display)
	// time.Sleep(1000 * time.Millisecond)
	// testFonts(ili9341Display)
	// testShapes(ili9341Display)
}

func drawLines(ili9341Display display.RGBDisplay) {
	xmax := float64(ili9341Display.ScreenWidth() - 1)
	ymax := float64(ili9341Display.ScreenHeight() - 1)
	xc := xmax / 2
	yc := ymax / 2
	radius := ymax / 2
	sAngle := math.Pi / 180 * 0
	rAngle := 2 * math.Pi
	dAngle := math.Pi / 180 * 5

	ili9341Display.SetBackgroundColor(colors.WHITE)
	ili9341Display.Clear()
	ili9341Display.SetColor(colors.BLUE)
	for angle := sAngle; angle < sAngle+rAngle; angle += dAngle {
		x := math.Cos(angle) * radius
		y := math.Sin(angle) * radius
		ili9341Display.Line(xc, yc, xc+x, yc+y)
	}

}

func drawCircle(ili9341Display display.RGBDisplay) {
	const N int = 3
	xmax := float64(ili9341Display.ScreenWidth() - 1)
	ymax := float64(ili9341Display.ScreenHeight() - 1)
	xc := xmax / 2
	yc := ymax / 2
	radius := ymax / 2.1
	xyc := [N][]float64{{xc, yc, radius}, {xc, yc, radius * .75}, {xc, yc, radius * .45}}
	colorset := [N]colors.Color{colors.BLACK, colors.DARKBLUE, colors.DARKGREEN}
	for i := 0; i < N; i++ {
		ili9341Display.SetColor(colorset[i])
		ili9341Display.Circle(xyc[i][0], xyc[i][1], xyc[i][2])
	}
}

func drawFillCircle(ili9341Display display.RGBDisplay) {
	const N int = 3
	xyc := [N][]float64{{30, 30, 45}, {160, 120, 115}, {400, 400, 250}}
	colorset := [N]colors.Color{colors.BLACK, colors.DARKBLUE, colors.DARKGREEN}
	for i := 0; i < N; i++ {
		ili9341Display.SetColor(colorset[i])
		ili9341Display.FillCircle(xyc[i][0], xyc[i][1], xyc[i][2])
	}
}

func drawThickCircle(ili9341Display display.RGBDisplay) {
	const N int = 3
	xmax := float64(ili9341Display.ScreenWidth() - 1)
	ymax := float64(ili9341Display.ScreenHeight() - 1)
	xc := xmax / 2
	yc := ymax / 2
	radius := ymax / 2.1
	xyc := [N][]float64{{xc, yc, radius}, {xc, yc, radius * .75}, {xc, yc, radius * .45}}
	colorset := [N]colors.Color{colors.ROYALBLUE, colors.SILVER, colors.MEDIUMSPRINGGREEN}
	widhTypes := [N]display.WidthType{display.INNER_WIDTH, display.CENTER_WIDTH, display.OUTER_WIDTH}
	const width = 10
	for i := 0; i < N; i++ {
		ili9341Display.SetColor(colorset[i])
		ili9341Display.ThickCircle(xyc[i][0], xyc[i][1], xyc[i][2], width, widhTypes[i])
		ili9341Display.SetColor(colors.RED)
		ili9341Display.Circle(xyc[i][0], xyc[i][1], xyc[i][2])
	}
}

func drawArc(ili9341Display display.RGBDisplay) {
	const N int = 12
	colorset := [N]colors.Color{
		colors.RED,
		colors.GREEN,
		colors.BLUE,
		colors.BLACK,
		colors.RED,
		colors.GREEN,
		colors.BLUE,
		colors.BLACK,
		colors.RED,
		colors.GREEN,
		colors.BLUE,
		colors.BLACK,
	}

	xyc := [N][]float64{
		{160, 120, 50, utils.ToRad(0), utils.ToRad(90)},
		{160, 120, 55, utils.ToRad(90), utils.ToRad(180)},
		{160, 120, 60, utils.ToRad(180), utils.ToRad(270)},
		{160, 120, 65, utils.ToRad(270), utils.ToRad(360)},
		{160, 120, 70, utils.ToRad(15), utils.ToRad(45)},
		{160, 120, 75, utils.ToRad(45), utils.ToRad(15)},
		{160, 120, 80, utils.ToRad(105), utils.ToRad(135)},
		{160, 120, 85, utils.ToRad(135), utils.ToRad(105)},
		{160, 120, 90, utils.ToRad(195), utils.ToRad(225)},
		{160, 120, 95, utils.ToRad(225), utils.ToRad(195)},
		{160, 120, 100, utils.ToRad(285), utils.ToRad(315)},
		{160, 120, 105, utils.ToRad(315), utils.ToRad(285)},
	}
	for i := 0; i < N; i++ {
		ili9341Display.SetColor(colorset[i])
		ili9341Display.Arc(xyc[i][0], xyc[i][1], xyc[i][2], xyc[i][3], xyc[i][4])
	}
	ili9341Display.SetColor(colors.RED)
	ili9341Display.Line(160, 0, 160, 239)
	ili9341Display.Line(0, 120, 319, 120)
}

func draThickwArc(ili9341Display display.RGBDisplay) {
	const N int = 3
	colorset := [N]colors.Color{
		colors.CYAN,
		colors.GREEN,
		colors.LIGHTBLUE,
	}

	widhTypes := [N]display.WidthType{display.OUTER_WIDTH, display.CENTER_WIDTH, display.INNER_WIDTH}
	xyc := [N][]float64{
		{160, 120, 70, utils.ToRad(45), utils.ToRad(175)},
		{160, 120, 90, utils.ToRad(15), utils.ToRad(300)},
		{160, 120, 115, utils.ToRad(300), utils.ToRad(15)},
	}

	for i := 0; i < N; i++ {
		ili9341Display.SetColor(colorset[i])
		ili9341Display.ThickArc(xyc[i][0], xyc[i][1], xyc[i][2], xyc[i][3], xyc[i][4], 10, widhTypes[i])
		ili9341Display.SetColor(colors.RED)
		ili9341Display.Arc(xyc[i][0], xyc[i][1], xyc[i][2], xyc[i][3], xyc[i][4])
	}
}

func drawRectangle(ili9341Display display.RGBDisplay) {
	const N int = 2
	xy := [N][]float64{{10, 10, 100, 100}, {50, 50, 200, 200}}
	colorset := [N]colors.Color{colors.BLUE, colors.GREEN}
	for i := 0; i < 2; i++ {
		ili9341Display.SetColor(colorset[i])
		ili9341Display.Rectangle(xy[i][0], xy[i][1], xy[i][2], xy[i][3])
	}

}

func drawFillRectangle(ili9341Display display.RGBDisplay) {
	const N int = 2
	xy := [N][]float64{{100, 100, 10, 10}, {50, 50, 200, 200}}
	colors := [N]colors.Color{colors.BLUE, colors.GREEN}
	for i := 0; i < 2; i++ {
		ili9341Display.SetColor(colors[i])
		ili9341Display.FillRectangle(xy[i][0], xy[i][1], xy[i][2], xy[i][3])
	}

}

func drawThickRectangle(ili9341Display display.RGBDisplay) {
	const N int = 3
	xy := [N][]float64{{100, 100, 10, 10}, {50, 50, 200, 200}, {100, 100, 300, 220}}
	colorset := [N]colors.Color{colors.ROYALBLUE, colors.NAVY, colors.FORESTGREEN}
	widhTypes := [N]display.WidthType{display.INNER_WIDTH, display.CENTER_WIDTH, display.OUTER_WIDTH}
	const width = 10
	for i := 0; i < N; i++ {
		ili9341Display.SetColor(colorset[i])
		ili9341Display.ThickRectangle(xy[i][0], xy[i][1], xy[i][2], xy[i][3], width, widhTypes[i])
		ili9341Display.SetColor(colors.RED)
		ili9341Display.Rectangle(xy[i][0], xy[i][1], xy[i][2], xy[i][3])
	}

}

// func testFonts(ili9341Display display.RGB565Display) {
// 	ili9341Display.SetBackgroundColor(colors.WHITE)
// 	ili9341Display.Clear()
// 	ili9341Display.MoveCursor(5, 5)
// 	ili9341Display.SetColor(colors.BLUE)
// 	ili9341Display.SetFont(fonts.Org_01)
// 	ili9341Display.SetLineHeight(40)
// 	ili9341Display.SetFont(fonts.FreeSans24pt7b)
// 	ili9341Display.Write("Hello Mark!")
// 	ili9341Display.SetFont(fonts.FreeMono18pt7b)
// 	ili9341Display.Write("Hello Mark!\n")
// 	ili9341Display.Write("0123456789")
//
// }

// func testColors(ili9341Display display.RGB565Display) {
// 	ili9341Display.SetBackgroundColor(colors.BLACK)
// 	ili9341Display.Clear()
// 	colors := []colors.Colors{colors.WHITE, colors.YELLOW, colors.GREEN, colors.BLUE, colors.RED}
// 	xmax := float64(ili9341Display.ScreenWidth() - 1)
// 	const height = 20
// 	const margin = 10
// 	for color := 0; color < len(colors); color++ {
// 		ys := float64(color * (height + margin))
// 		ili9341Display.SetColor(colors[color])
// 		ili9341Display.FillRectangle(0, ys, xmax, ys+height)
// 	}
//
// }

// func testShapes(ili9341Display display.RGB565Display) {
// 	ili9341Display.SetBackgroundColor(colors.BLUE)
// 	ili9341Display.Clear()
// 	ili9341Display.SetColor(colors.YELLOW)
// 	ili9341Display.Circle(50, 50, 30)
// 	ili9341Display.SetColor(colors.GREEN)
// 	ili9341Display.FillCircle(100, 100, 30)
// 	ili9341Display.Arc(120, 120, 118, -math.Pi/4, math.Pi/4, 40)
// 	// ili9341Display.SetColor(colors.RED)
// 	// ili9341Display.FillRectangle(50, 150, 220, 180)
//
// }

func createGpioOutPin(gpioPinNum string) gpio.PinOut {
	var pin gpio.PinOut = gpioreg.ByName(gpioPinNum)
	if pin == nil {
		checkFatalErr(fmt.Errorf("failed to create GPIO pin %s", gpioPinNum))
	}
	pin.Out(gpio.Low)
	return pin
}

func createSPIConnection(busNumber int, chipSelect int) spi.Conn {
	spibus, _ := sysfs.NewSPI(
		busNumber,
		chipSelect,
	)
	spiConn, err := spibus.Connect(
		physic.Frequency(12)*physic.MegaHertz,
		spi.Mode3,
		8,
	)
	checkFatalErr(err)
	return spiConn
}
