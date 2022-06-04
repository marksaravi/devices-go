package rgb565

type RGB565 uint16

const (
	BLACK  RGB565 = 0b0
	GREEN  RGB565 = 0x003f
	BLUE   RGB565 = 0x1f << 6
	RED    RGB565 = 0x1f << 11
	WHITE  RGB565 = RED | GREEN | BLUE
	YELLOW RGB565 = RED | GREEN
)
