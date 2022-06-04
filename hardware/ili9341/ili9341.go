package ili9341

import (
	"errors"
	"math"
	"sync"
	"time"

	"github.com/marksaravi/devices-go/v1/colors/rgb565"
	"github.com/marksaravi/fonts-go/fonts"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/spi"
)

const (
	lcd_width                int  = 320 //LCD width
	lcd_height               int  = 240 //LCD height
	segment_width            int  = 16
	segment_height           int  = 12
	num_x_seg                int  = lcd_width / segment_width
	num_y_seg                int  = lcd_height / segment_height
	num_of_segments          int  = num_x_seg * num_y_seg
	bytes_per_segments       int  = segment_width * segment_height * 2
	row_address_order        byte = 1
	column_address_order     byte = 0
	row_col_exchange         byte = 1
	vertical_refresh_order   byte = 0
	rgb_bgr_order            byte = 0 // 0 RGB
	horizontal_refresh_order byte = 0
	memeory_access_control   byte = (row_address_order << 7) |
		(column_address_order << 6) |
		(row_col_exchange << 5) |
		(vertical_refresh_order << 4) |
		(rgb_bgr_order << 3) |
		(horizontal_refresh_order << 2) | 0b00000000
)

const (
	screen_left_padding int = 5
	screen_top_padding  int = 5
)

type device struct {
	conn              spi.Conn
	pinDC             gpio.PinOut // WriteDataByte/WriteCommand
	pinRST            gpio.PinOut // Reset
	mu                sync.Mutex
	segments          []byte
	isSegmentChanged  []bool
	changedSegment    chan int
	cursorX           int
	cursorY           int
	font              fonts.BitmapFont
	backgroundColor   rgb565.RGB565
	color             rgb565.RGB565
	letterSpacing     int
	lineHeight        int
	screenLeftPadding int
	screenTopPadding  int
}

func NewILI9341(
	spiConn spi.Conn,
	pinDC gpio.PinOut,
	pinRST gpio.PinOut,
) (*device, error) {
	d := &device{
		conn:              spiConn,
		pinDC:             pinDC,
		pinRST:            pinRST,
		segments:          make([]byte, num_of_segments*bytes_per_segments),
		isSegmentChanged:  make([]bool, num_of_segments),
		changedSegment:    make(chan int),
		cursorX:           screen_left_padding,
		cursorY:           screen_top_padding,
		font:              fonts.FreeSans24pt7b,
		color:             rgb565.BLACK,
		letterSpacing:     0,
		lineHeight:        32,
		screenLeftPadding: screen_top_padding,
		screenTopPadding:  screen_top_padding,
	}
	d.initLCD()
	d.startDeviceMemoryWriter()
	return d, nil
}

func (dev *device) WriteCommand(cmd byte) {
	dev.pinDC.Out(gpio.Low)
	dev.conn.Tx([]byte{cmd}, nil)
}

func (dev *device) WriteDataByte(data byte) (byte, error) {
	dev.pinDC.Out(gpio.High)
	res := make([]byte, 1)
	err := dev.conn.Tx([]byte{data}, res)
	return res[0], err
}

func (dev *device) refreshSegment(seg int) {
	start := seg * bytes_per_segments
	xseg := seg / num_x_seg
	yseg := seg % num_x_seg
	dev.SetWindow(xseg*segment_width, yseg*segment_height, (xseg+1)*segment_width-1, (yseg+1)*segment_height-1)
	dev.pinDC.Out(gpio.High)
	dev.conn.Tx(dev.segments[start:start+bytes_per_segments], nil)
}

func (dev *device) startDeviceMemoryWriter() {
	go func() {
		for {
			seg := <-dev.changedSegment
			dev.refreshSegment(seg)
		}
	}()
}

func (dev *device) initLCD() {
	dev.reset()

	dev.WriteCommand(0x11) //Sleep out

	dev.WriteCommand(0xCF)
	dev.WriteDataByte(0x00)
	dev.WriteDataByte(0xC1)
	dev.WriteDataByte(0x30)
	dev.WriteCommand(0xED)
	dev.WriteDataByte(0x64)
	dev.WriteDataByte(0x03)
	dev.WriteDataByte(0x12)
	dev.WriteDataByte(0x81)
	dev.WriteCommand(0xE8)
	dev.WriteDataByte(0x85)
	dev.WriteDataByte(0x00)
	dev.WriteDataByte(0x79)
	dev.WriteCommand(0xCB)
	dev.WriteDataByte(0x39)
	dev.WriteDataByte(0x2C)
	dev.WriteDataByte(0x00)
	dev.WriteDataByte(0x34)
	dev.WriteDataByte(0x02)
	dev.WriteCommand(0xF7)
	dev.WriteDataByte(0x20)
	dev.WriteCommand(0xEA)
	dev.WriteDataByte(0x00)
	dev.WriteDataByte(0x00)
	dev.WriteCommand(0xC0)  //Power control
	dev.WriteDataByte(0x1D) //VRH[5:0]
	dev.WriteCommand(0xC1)  //Power control
	dev.WriteDataByte(0x12) //SAP[2:0];BT[3:0]
	dev.WriteCommand(0xC5)  //VCM control
	dev.WriteDataByte(0x33)
	dev.WriteDataByte(0x3F)
	dev.WriteCommand(0xC7) //VCM control
	dev.WriteDataByte(0x92)
	dev.WriteCommand(0x3A) // Memory Access Control
	dev.WriteDataByte(0x55)
	dev.WriteCommand(0x36) // Memory Access Control
	dev.WriteDataByte(memeory_access_control)
	dev.WriteCommand(0xB1)
	dev.WriteDataByte(0x00)
	dev.WriteDataByte(0x12)
	dev.WriteCommand(0xB6) // Display Function Control
	dev.WriteDataByte(0x0A)
	dev.WriteDataByte(0xA2)

	dev.WriteCommand(0x44)
	dev.WriteDataByte(0x02)

	dev.WriteCommand(0xF2) // 3Gamma Function Disable
	dev.WriteDataByte(0x00)
	dev.WriteCommand(0x26) //Gamma curve selected
	dev.WriteDataByte(0x01)
	dev.WriteCommand(0xE0) //Set Gamma
	dev.WriteDataByte(0x0F)
	dev.WriteDataByte(0x22)
	dev.WriteDataByte(0x1C)
	dev.WriteDataByte(0x1B)
	dev.WriteDataByte(0x08)
	dev.WriteDataByte(0x0F)
	dev.WriteDataByte(0x48)
	dev.WriteDataByte(0xB8)
	dev.WriteDataByte(0x34)
	dev.WriteDataByte(0x05)
	dev.WriteDataByte(0x0C)
	dev.WriteDataByte(0x09)
	dev.WriteDataByte(0x0F)
	dev.WriteDataByte(0x07)
	dev.WriteDataByte(0x00)
	dev.WriteCommand(0xE1) //Set Gamma
	dev.WriteDataByte(0x00)
	dev.WriteDataByte(0x23)
	dev.WriteDataByte(0x24)
	dev.WriteDataByte(0x07)
	dev.WriteDataByte(0x10)
	dev.WriteDataByte(0x07)
	dev.WriteDataByte(0x38)
	dev.WriteDataByte(0x47)
	dev.WriteDataByte(0x4B)
	dev.WriteDataByte(0x0A)
	dev.WriteDataByte(0x13)
	dev.WriteDataByte(0x06)
	dev.WriteDataByte(0x30)
	dev.WriteDataByte(0x38)
	dev.WriteDataByte(0x0F)
	dev.WriteCommand(0x29) //Display on
}

func (dev *device) reset() {
	const SLEEP_MS = 120
	dev.pinRST.Out(gpio.High)
	time.Sleep(time.Millisecond * SLEEP_MS)
	dev.pinRST.Out(gpio.Low)
	time.Sleep(time.Millisecond * SLEEP_MS)
	dev.pinRST.Out(gpio.High)
	time.Sleep(time.Millisecond * SLEEP_MS)
}

func (dev *device) Update() {
	for seg := 0; seg < num_of_segments; seg++ {
		if dev.isSegmentChanged[seg] {
			dev.changedSegment <- seg
			dev.isSegmentChanged[seg] = false
		}
	}
}

func (dev *device) pixel(x, y float64, color rgb565.RGB565) {
	xi := int(math.Round(x))
	yi := int(math.Round(y))
	if xi < 0 || yi < 0 || xi >= lcd_width || yi >= lcd_height {
		return
	}
	xseg := xi / segment_width
	yseg := yi / segment_height
	seg := xseg*num_x_seg + yseg
	xoffs := xi % segment_width
	yoffs := yi % segment_height
	i := seg*bytes_per_segments + (yoffs*segment_width+xoffs)*2
	dev.mu.Lock()
	dev.segments[i] = byte(color)
	dev.segments[i+1] = byte(color >> 8)
	dev.isSegmentChanged[seg] = true
	dev.mu.Unlock()
}

func (dev *device) Pixel(x, y float64) {
	dev.pixel(x, y, dev.color)
}

func (dev *device) Circle(x, y, radius float64) {
	dangle := math.Pi / 180

	for angle := float64(0); angle < math.Pi*2; angle += dangle {
		dev.Pixel(x+radius*math.Cos(angle), y+radius*math.Sin(angle))
	}
}

func (dev *device) Rectangle(x1, y1, x2, y2 float64) {
	dev.Line(x1, y1, x2, y1)
	dev.Line(x2, y1, x2, y2)
	dev.Line(x2, y2, x1, y2)
	dev.Line(x1, y2, x1, y1)
}

func (dev *device) FillRectangle(x1, y1, x2, y2 float64) {
	xs := x1
	xe := x2
	if x2 < x1 {
		xs = x2
		xe = x1
	}
	for x := xs; x <= xe; x++ {
		dev.verLine(x, y1, y2)
	}
}

func (dev *device) FillCircle(x, y, radius float64) {
	dangle := math.Pi / 180

	for r := radius; r > 0; r -= 1 {
		for angle := float64(0); angle < math.Pi*2; angle += dangle {
			dev.Pixel(x+r*math.Cos(angle), y+r*math.Sin(angle))
		}
	}
}

func (dev *device) Line(x1, y1, x2, y2 float64) {
	xs := x1
	ys := y1
	xe := x2
	ye := y2
	if x2 < x1 {
		xs = x2
		ys = y2
		xe = x1
		ye = y1
	}

	if xs == xe {
		dev.verLine(xs, ys, ye)
	} else {
		dev.sloppedLine(xs, ys, xe, ye)
	}
}

func (dev *device) sloppedLine(xs, ys, xe, ye float64) {
	a := (ye - ys) / (xe - xs)
	b := (ys) - a*(xs)
	var x, y float64
	for x = xs; x <= xe; x += 1 {
		y = a*float64(x) + b
		dev.Pixel(x, y)
	}
}

func (dev *device) verLine(x, y1, y2 float64) {
	var y float64 = 0
	ys := y1
	ye := y2
	if y2 < y1 {
		ys = y2
		ye = y1
	}
	for y = ys; y <= ye; y += 1 {
		dev.Pixel(x, y)
	}
}

func (dev *device) ShowSegment(segX, segY int, color rgb565.RGB565, state string) {
	xOffs := float64(segX * segment_width)
	yOffs := float64(segY * segment_height)
	x1 := xOffs
	y1 := yOffs
	x2 := xOffs + float64(segment_width-1)
	y2 := yOffs + float64(segment_height-1)
	dev.Rectangle(x1, y1, x2, y2)
	xm := x1 + (x2-x1)/2
	ym := y1 + (y2-y1)/2

	if state == "dir" {
		dev.Line(x1+3, ym, x2-3, ym)
		dev.Line(x2-3, ym, x2-6, ym-3)
		dev.Line(x2-3, ym, x2-6, ym+3)
	}
	if state == "first" {
		dev.FillCircle(xm, ym, 3)
	}
	if state == "last" {
		const offset = 3
		dev.FillRectangle(x1+offset, y1+offset, x2-offset, y2-offset)
	}
}

func (dev *device) Clear() {
	for y := 0; y < lcd_height; y++ {
		for x := 0; x < lcd_width; x++ {
			dev.pixel(float64(x), float64(y), dev.backgroundColor)
		}
	}
	dev.Update()
}

func (dev *device) SetWindow(xStart, yStart, xEnd, yEnd int) {
	dev.WriteCommand(0x2a)
	dev.WriteDataByte(byte(xStart >> 8))
	dev.WriteDataByte(byte(xStart & 0xff))
	dev.WriteDataByte(byte(xEnd >> 8))
	dev.WriteDataByte(byte(xEnd & 0xff))

	dev.WriteCommand(0x2b)
	dev.WriteDataByte(byte(yStart >> 8))
	dev.WriteDataByte(byte(yStart & 0xff))
	dev.WriteDataByte(byte(yEnd >> 8))
	dev.WriteDataByte(byte(yEnd & 0xff))
	dev.WriteCommand(0x2C)
}

func (dev *device) ScreenWidth() int {
	return lcd_width
}

func (dev *device) ScreenHeight() int {
	return lcd_height
}

func (dev *device) NumOfXSegments() int {
	return num_x_seg
}

func (dev *device) NumOfYSegments() int {
	return num_y_seg
}

func (dev *device) MoveCursor(x, y int) {
	dev.cursorX = x
	dev.cursorY = y
}

func (dev *device) SetFont(font fonts.BitmapFont) {
	dev.font = font
}

func (dev *device) SetBackgroundColor(color rgb565.RGB565) {
	dev.backgroundColor = color
}

func (dev *device) SetColor(color rgb565.RGB565) {
	dev.color = color
}

func (dev *device) SetLineHeight(height int) {
	dev.lineHeight = height
}

func (dev *device) WriteChar(char byte) error {
	if char == '\n' {
		dev.nextLine()
		return nil
	}
	if char < ' ' || char > '~' {
		return errors.New("charCode code out of range")
	}

	dev.drawBitmapChar(char)
	return nil
}

func (dev *device) drawBitmapChar(char byte) error {
	glyph := dev.font.Glyphs[char-0x20]
	for h := 0; h < glyph.Height; h++ {
		for w := 0; w < glyph.Width; w++ {
			bitIndex := h*glyph.Width + w
			shift := byte(bitIndex) % 8
			d := dev.font.Bitmap[glyph.BitmapOffset+bitIndex/8]
			mask := byte(0b10000000) >> shift
			bit := d & mask
			color := dev.backgroundColor
			if bit != 0 {
				color = dev.color
			}
			x := float64(dev.cursorX + w + glyph.XOffset)
			y := float64(dev.cursorY + dev.lineHeight + h + glyph.YOffset)
			dev.pixel(x, y, color)
		}
	}
	xforward := glyph.XAdvance
	dev.cursorX += xforward
	if dev.cursorX+xforward >= dev.ScreenWidth() {
		dev.nextLine()
	}

	return nil
}

func (dev *device) nextLine() {
	dev.cursorX = dev.screenLeftPadding
	dev.cursorY += dev.lineHeight
}

func (dev *device) Write(text string) {
	for i := 0; i < len(text); i++ {
		dev.WriteChar(text[i])
	}
}
