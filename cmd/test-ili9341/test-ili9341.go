package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/marksaravi/devices-go/colors/rgb565"
	"github.com/marksaravi/devices-go/devices/display"
	"github.com/marksaravi/devices-go/hardware/ili9341"
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
	// time.Sleep(1000 * time.Millisecond)
	testLines(ili9341Display)
	// time.Sleep(1000 * time.Millisecond)
	// testColors(ili9341Display)
	// time.Sleep(1000 * time.Millisecond)
	// testFonts(ili9341Display)
	// testShapes(ili9341Display)
	time.Sleep(1000 * time.Millisecond)
}

func testLines(ili9341Display display.RGBDisplay) {
	ili9341Display.SetBackgroundColor(rgb565.WHITE)
	ili9341Display.Clear()
	ili9341Display.SetColor(rgb565.BLUE)
	xmax := float64(ili9341Display.ScreenWidth() - 1)
	ymax := float64(ili9341Display.ScreenHeight() - 1)
	xc := xmax / 2
	yc := ymax / 2
	radius := ymax / 2

	sAngle := math.Pi / 180 * 0
	rAngle := 2 * math.Pi
	dAngle := math.Pi / 180 * 5
	for angle := sAngle; angle < sAngle+rAngle; angle += dAngle {
		x := math.Cos(angle) * radius
		y := math.Sin(angle) * radius
		ili9341Display.Line(xc, yc, xc+x, yc+y) // error
		// fmt.Println(angle, x, y)
	}
	ili9341Display.Update()
}

// func testFonts(ili9341Display display.RGB565Display) {
// 	ili9341Display.SetBackgroundColor(rgb565.WHITE)
// 	ili9341Display.Clear()
// 	ili9341Display.MoveCursor(5, 5)
// 	ili9341Display.SetColor(rgb565.BLUE)
// 	ili9341Display.SetFont(fonts.Org_01)
// 	ili9341Display.SetLineHeight(40)
// 	ili9341Display.SetFont(fonts.FreeSans24pt7b)
// 	ili9341Display.Write("Hello Mark!")
// 	ili9341Display.SetFont(fonts.FreeMono18pt7b)
// 	ili9341Display.Write("Hello Mark!\n")
// 	ili9341Display.Write("0123456789")
// 	ili9341Display.Update()
// }

// func testColors(ili9341Display display.RGB565Display) {
// 	ili9341Display.SetBackgroundColor(rgb565.BLACK)
// 	ili9341Display.Clear()
// 	colors := []rgb565.RGB565{rgb565.WHITE, rgb565.YELLOW, rgb565.GREEN, rgb565.BLUE, rgb565.RED}
// 	xmax := float64(ili9341Display.ScreenWidth() - 1)
// 	const height = 20
// 	const margin = 10
// 	for color := 0; color < len(colors); color++ {
// 		ys := float64(color * (height + margin))
// 		ili9341Display.SetColor(colors[color])
// 		ili9341Display.FillRectangle(0, ys, xmax, ys+height)
// 	}
// 	ili9341Display.Update()
// }

// func testShapes(ili9341Display display.RGB565Display) {
// 	ili9341Display.SetBackgroundColor(rgb565.BLUE)
// 	ili9341Display.Clear()
// 	ili9341Display.SetColor(rgb565.YELLOW)
// 	ili9341Display.Circle(50, 50, 30)
// 	ili9341Display.SetColor(rgb565.GREEN)
// 	ili9341Display.FillCircle(100, 100, 30)
// 	ili9341Display.Arc(120, 120, 118, -math.Pi/4, math.Pi/4, 40)
// 	// ili9341Display.SetColor(rgb565.RED)
// 	// ili9341Display.FillRectangle(50, 150, 220, 180)
// 	ili9341Display.Update()
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
