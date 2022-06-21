package ili9341

import (
	"math"
	"sync"
	"time"

	"github.com/marksaravi/devices-go/colors/rgb565"
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

type device struct {
	conn             spi.Conn
	pinDC            gpio.PinOut // WriteDataByte/writeCommand
	pinRST           gpio.PinOut // Reset
	mu               sync.Mutex
	segments         []byte
	isSegmentChanged []bool
	changedSegment   chan int
	color            rgb565.RGB565
}

func NewILI9341(
	spiConn spi.Conn,
	pinDC gpio.PinOut,
	pinRST gpio.PinOut,
) (*device, error) {
	d := &device{
		conn:             spiConn,
		pinDC:            pinDC,
		pinRST:           pinRST,
		segments:         make([]byte, num_of_segments*bytes_per_segments),
		isSegmentChanged: make([]bool, num_of_segments),
		changedSegment:   make(chan int),
		color:            rgb565.BLACK,
	}
	d.initLCD()
	d.startDeviceMemoryWriter()
	return d, nil
}

func (dev *device) Update() {
	for seg := 0; seg < num_of_segments; seg++ {
		if dev.isSegmentChanged[seg] {
			dev.changedSegment <- seg
			dev.isSegmentChanged[seg] = false
		}
	}
}

func (dev *device) Pixel(x, y float64) {
	dev.pixel(x, y, dev.color)
}

func (dev *device) SetColor(color rgb565.RGB565) {
	dev.color = color
}

func (dev *device) ScreenWidth() float64 {
	return float64(lcd_width)
}

func (dev *device) ScreenHeight() float64 {
	return float64(lcd_height)
}

func (dev *device) writeCommand(cmd byte) {
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
	dev.setWindow(xseg*segment_width, yseg*segment_height, (xseg+1)*segment_width-1, (yseg+1)*segment_height-1)
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

	dev.writeCommand(0x11) //Sleep out

	dev.writeCommand(0xCF)
	dev.WriteDataByte(0x00)
	dev.WriteDataByte(0xC1)
	dev.WriteDataByte(0x30)
	dev.writeCommand(0xED)
	dev.WriteDataByte(0x64)
	dev.WriteDataByte(0x03)
	dev.WriteDataByte(0x12)
	dev.WriteDataByte(0x81)
	dev.writeCommand(0xE8)
	dev.WriteDataByte(0x85)
	dev.WriteDataByte(0x00)
	dev.WriteDataByte(0x79)
	dev.writeCommand(0xCB)
	dev.WriteDataByte(0x39)
	dev.WriteDataByte(0x2C)
	dev.WriteDataByte(0x00)
	dev.WriteDataByte(0x34)
	dev.WriteDataByte(0x02)
	dev.writeCommand(0xF7)
	dev.WriteDataByte(0x20)
	dev.writeCommand(0xEA)
	dev.WriteDataByte(0x00)
	dev.WriteDataByte(0x00)
	dev.writeCommand(0xC0)  //Power control
	dev.WriteDataByte(0x1D) //VRH[5:0]
	dev.writeCommand(0xC1)  //Power control
	dev.WriteDataByte(0x12) //SAP[2:0];BT[3:0]
	dev.writeCommand(0xC5)  //VCM control
	dev.WriteDataByte(0x33)
	dev.WriteDataByte(0x3F)
	dev.writeCommand(0xC7) //VCM control
	dev.WriteDataByte(0x92)
	dev.writeCommand(0x3A) // Memory Access Control
	dev.WriteDataByte(0x55)
	dev.writeCommand(0x36) // Memory Access Control
	dev.WriteDataByte(memeory_access_control)
	dev.writeCommand(0xB1)
	dev.WriteDataByte(0x00)
	dev.WriteDataByte(0x12)
	dev.writeCommand(0xB6) // Display Function Control
	dev.WriteDataByte(0x0A)
	dev.WriteDataByte(0xA2)

	dev.writeCommand(0x44)
	dev.WriteDataByte(0x02)

	dev.writeCommand(0xF2) // 3Gamma Function Disable
	dev.WriteDataByte(0x00)
	dev.writeCommand(0x26) //Gamma curve selected
	dev.WriteDataByte(0x01)
	dev.writeCommand(0xE0) //Set Gamma
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
	dev.writeCommand(0xE1) //Set Gamma
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
	dev.writeCommand(0x29) //Display on
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

func (dev *device) setWindow(xStart, yStart, xEnd, yEnd int) {
	dev.writeCommand(0x2a)
	dev.WriteDataByte(byte(xStart >> 8))
	dev.WriteDataByte(byte(xStart & 0xff))
	dev.WriteDataByte(byte(xEnd >> 8))
	dev.WriteDataByte(byte(xEnd & 0xff))

	dev.writeCommand(0x2b)
	dev.WriteDataByte(byte(yStart >> 8))
	dev.WriteDataByte(byte(yStart & 0xff))
	dev.WriteDataByte(byte(yEnd >> 8))
	dev.WriteDataByte(byte(yEnd & 0xff))
	dev.writeCommand(0x2C)
}
