package display

import (
	"errors"
	"fmt"

	"github.com/marksaravi/fonts-go/fonts"
)

func (dev *rgbDevice) SetFont(font interface{}) error {
	dev.font = font
	if bitmapfont, ok := font.(fonts.BitmapFont); ok {
		dev.initBitmapFonts(bitmapfont)
		return nil
	}
	return errors.New("font format is not implemented")
}

func (dev *rgbDevice) SetLineWrapping(wrapping bool) {
	dev.lineWrapping = wrapping
}

func (dev *rgbDevice) SetLineHeight(height int) {
	dev.lineHeight = height
}

func (dev *rgbDevice) SetLetterSpacing(spacing int) {
	dev.letterSpacing = spacing
}

func (dev *rgbDevice) WriteChar(char byte) error {
	if char == '\n' {
		dev.nextLine()
		return nil
	}
	if char < ' ' || char > '~' {
		return errors.New("charCode code out of range")
	}

	switch dev.fontType {
	case BITMAP_FONT:
		dev.drawBitmapChar(char)
	default:
		return errors.New("font is not defined")
	}
	dev.moveCursorForward()
	return nil
}

func (dev *rgbDevice) Write(text string) {
	for i := 0; i < len(text); i++ {
		dev.WriteChar(text[i])
	}
}

func (dev *rgbDevice) MoveCursor(col, row int) {
	dev.cursorX = col
	dev.cursorY = row
}

func (dev *rgbDevice) nextLine() {
	dev.cursorX = 0
	dev.cursorY += 1
}

func (dev *rgbDevice) moveCursorForward() {
	dev.cursorX += 1
	if dev.lineWrapping && dev.calcCharX() >= dev.ScreenWidth() {
		dev.cursorX = 0
		dev.cursorY += 1
	}
}

func (dev *rgbDevice) calcCharX() int {
	return dev.textLeftPadding + dev.cursorX*(dev.charWidth+dev.letterSpacing)
}

func (dev *rgbDevice) calcCharY() int {
	return dev.textTopPadding + dev.lineHeight*dev.cursorY
}

func (dev *rgbDevice) charXY() (int, int) {
	x := dev.calcCharX()
	y := dev.calcCharY()

	return x, y
}

func (dev *rgbDevice) drawBitmapChar(char byte) error {
	glyph := dev.bitmapFont.Glyphs[char-0x20]
	charx, chary := dev.charXY()
	fmt.Println("x,y: ", charx, chary)
	for h := 0; h < glyph.Height; h++ {
		for w := 0; w < glyph.Width; w++ {
			bitIndex := h*glyph.Width + w
			shift := byte(bitIndex) % 8
			d := dev.bitmapFont.Bitmap[glyph.BitmapOffset+bitIndex/8]
			mask := byte(0b10000000) >> shift
			bit := d & mask
			color := dev.bgColor
			if bit != 0 {
				color = dev.color
			}
			fmt.Println("w,h: ", chary, glyph.YOffset, w, h)
			dev.pixeldev.Pixel(charx+w+glyph.XOffset, chary+h+glyph.YOffset+dev.lineHeight, color)
		}
	}
	return nil
}

func (dev *rgbDevice) initBitmapFonts(bitmapfont fonts.BitmapFont) {
	dev.fontType = BITMAP_FONT
	dev.bitmapFont = bitmapfont

	glyph := dev.bitmapFont.Glyphs[0x20]
	dev.charWidth = glyph.XAdvance
	fmt.Println("charWidth: ", dev.charWidth)
	if dev.lineHeight == 0 {
		dev.lineHeight = glyph.YOffset
	}
}
