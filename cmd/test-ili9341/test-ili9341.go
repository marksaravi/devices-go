package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/marksaravi/devices-go/colors/rgb565"
	"github.com/marksaravi/devices-go/devices/display"
	"github.com/marksaravi/devices-go/hardware/ili9341"
	"github.com/marksaravi/fonts-go/fonts"
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
	var display display.RGB565Display
	var err error
	display, err = ili9341.NewILI9341(spiConn, dataCommandSelect, reset)
	checkFatalErr(err)
	// time.Sleep(1000 * time.Millisecond)
	// testLines(display)
	// time.Sleep(1000 * time.Millisecond)
	// testColors(display)
	// time.Sleep(1000 * time.Millisecond)
	testFonts(display)
	time.Sleep(1000 * time.Millisecond)
	testShapes(display)
	time.Sleep(1000 * time.Millisecond)
}

func testLines(display display.RGB565Display) {
	display.SetBackgroundColor(rgb565.BLACK)
	display.Clear()
	xmax := float64(display.ScreenWidth() - 1)
	ymax := float64(display.ScreenHeight() - 1)
	display.SetColor(rgb565.GREEN)
	display.Line(0, 0, xmax, ymax)
	display.Line(0, ymax, xmax, 0)
	display.SetColor(rgb565.YELLOW)
	display.Line(0, 0, xmax, 0)
	display.Line(xmax, 0, xmax, ymax)
	display.Line(xmax, ymax, 0, ymax)
	display.Line(0, ymax, 0, 0)
	display.Update()
}

func testFonts(display display.RGB565Display) {
	display.SetBackgroundColor(rgb565.WHITE)
	display.Clear()
	display.MoveCursor(5, 5)
	display.SetColor(rgb565.BLUE)
	display.SetFont(fonts.Org_01)
	display.SetLineHeight(40)
	display.SetFont(fonts.FreeSans24pt7b)
	display.Write("Hello Mark!")
	display.SetFont(fonts.FreeMono18pt7b)
	display.Write("Hello Mark!\n")
	display.Write("0123456789")
	display.Update()
}

func testColors(display display.RGB565Display) {
	display.SetBackgroundColor(rgb565.BLACK)
	display.Clear()
	colors := []rgb565.RGB565{rgb565.WHITE, rgb565.YELLOW, rgb565.GREEN, rgb565.BLUE, rgb565.RED}
	xmax := float64(display.ScreenWidth() - 1)
	const height = 20
	const margin = 10
	for color := 0; color < len(colors); color++ {
		ys := float64(color * (height + margin))
		display.SetColor(colors[color])
		display.FillRectangle(0, ys, xmax, ys+height)
	}
	display.Update()
}

func testShapes(display display.RGB565Display) {
	display.SetBackgroundColor(rgb565.BLUE)
	display.Clear()
	display.SetColor(rgb565.YELLOW)
	display.Circle(50, 50, 30)
	display.SetColor(rgb565.GREEN)
	display.FillCircle(100, 100, 30)
	display.Arc(120, 120, 118, -math.Pi/4, math.Pi/4, 40)
	// display.SetColor(rgb565.RED)
	// display.FillRectangle(50, 150, 220, 180)
	display.Update()
}

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
