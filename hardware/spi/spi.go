package spi

type SPI interface {
	Tx(w, r []byte) error
}
