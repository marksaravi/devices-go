package main

import (
	"fmt"
	"log"
	"time"

	"github.com/marksaravi/devices-go/v1/colors/rgb565"
	"github.com/marksaravi/devices-go/v1/devices/display"
	"github.com/marksaravi/devices-go/v1/hardware/ili9341"
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
	// testAllScreen(display)
	// time.Sleep(1000 * time.Millisecond)
	// testLines(display)
	// time.Sleep(1000 * time.Millisecond)
	// testColors(display)
	// time.Sleep(1000 * time.Millisecond)
	// color := display.RED
	// for color < display.WHITE {
	// 	fmt.Println(color)
	// 	color = testColorsPallet(display, color)
	// 	time.Sleep(1000 * time.Millisecond)
	// }
	// time.Sleep(1000 * time.Millisecond)
	testFonts(display)
	time.Sleep(1000 * time.Millisecond)
}

// func testAllScreen(display display.RGB565Display) {
// 	display.Clear(display.WHITE)

// 	for x := 0; x < display.NumOfXSegments(); x++ {
// 		for y := 0; y < display.NumOfYSegments(); y++ {
// 			state := "dir"
// 			color := display.GREEN
// 			if x == 0 && y == 0 {
// 				state = "first"
// 				color = display.BLUE
// 			}
// 			if x == display.NumOfXSegments()-1 && y == display.NumOfYSegments()-1 {
// 				state = "last"
// 				color = display.RED
// 			}
// 			display.ShowSegment(x, y, color, state)
// 		}
// 	}
// 	display.Update()
// }

// func testLines(display display.RGB565Display) {
// 	display.Clear(display.BLACK)
// 	xmax := float64(display.ScreenWidth() - 1)
// 	ymax := float64(display.ScreenHeight() - 1)
// 	display.Line(0, 0, xmax, ymax, display.GREEN)
// 	display.Line(0, ymax, xmax, 0, display.GREEN)
// 	display.Line(0, 0, xmax, 0, display.YELLOW)
// 	display.Line(xmax, 0, xmax, ymax, display.YELLOW)
// 	display.Line(xmax, ymax, 0, ymax, display.YELLOW)
// 	display.Line(0, ymax, 0, 0, display.YELLOW)
// 	display.Update()
// }

func testFonts(display display.RGB565Display) {
	display.Clear(rgb565.WHITE)
	display.MoveCursor(5, 5)
	display.SetFontBackgroundColor(rgb565.WHITE)
	display.SetFontColor(rgb565.BLUE)
	display.SetFont(fonts.Org_01)
	display.SetLineHeight(40)
	// display.SetFont(fonts.Sans24)
	// display.Write("Hello Mark!")
	// display.SetFont(fonts.Font24)
	// display.Write("Hello Mark!")
	// display.SetFontColor(display.YELLOW)
	// display.SetFontBackgroundColor(display.BLACK)
	// for c := ' '; c <= '~'; c++ {
	// 	display.WriteChar(byte(c))
	// }

	display.Write("Hello Mark!\n")
	display.Write("0123456789")
	// display.SetFont(fonts.Font24)
	// display.Write("Hello Mark!")
	fmt.Println()
	display.Update()
}

// func testColors(display display.RGB565Display) {
// 	display.Clear(display.BLACK)
// 	colors := []display.RGB565{display.WHITE, display.YELLOW, display.GREEN, display.BLUE, display.RED}
// 	xmax := float64(display.ScreenWidth() - 1)
// 	const height = 20
// 	const margin = 10
// 	for color := 0; color < len(colors); color++ {
// 		ys := float64(color * (height + margin))
// 		display.FillRectangle(0, ys, xmax, ys+height, colors[color])
// 	}
// 	display.Update()
// }

// func testColorsPallet(display display.RGB565Display, color display.RGB565) display.RGB565 {
// 	display.Clear(display.BLACK)
// 	height := 3
// 	xmax := float64(display.ScreenWidth() - 1)
// 	n := display.ScreenHeight() / height
// 	c := color
// 	for i := uint16(0); i < uint16(n); i++ {
// 		y := float64(i * uint16(height))
// 		display.FillRectangle(0, y, xmax, y+float64(height), c)
// 		c += 1
// 	}
// 	display.Update()
// 	return c
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
