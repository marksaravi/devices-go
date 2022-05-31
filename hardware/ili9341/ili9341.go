package ili9341

import "periph.io/x/conn/v3/spi"

type device struct {
	conn spi.Conn
}

func NewILI9341(conn spi.Conn) *device {
	return &device{
		conn: conn,
	}
}
