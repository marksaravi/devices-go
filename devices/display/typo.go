package display

import (
	"errors"

	"github.com/marksaravi/fonts-go/fonts"
)

func (dev *rgbDevice) SetFont(font interface{}) error {
	dev.font = font
	if bitmapfont, ok := font.(fonts.BitmapFont); ok {
		dev.fontType = BITMAP_FONT
		dev.bitmapFont = bitmapfont
		return nil
	}
	return errors.New("font format is not implemented")
}

func (dev *rgbDevice) writeChar(char byte) error {
	if char < ' ' || char > '~' {
		return errors.New("charCode code out of range")
	}

	switch dev.fontType {
	case BITMAP_FONT:
		dev.drawBitmapChar(char)
	default:
		return errors.New("font is not defined")
	}
	return nil
}

func (dev *rgbDevice) Write(text string) {
	for i := 0; i < len(text); i++ {
		dev.writeChar(text[i])
	}
}

func (dev *rgbDevice) MoveCursor(x, y int) {
	dev.cursorX = x
	dev.cursorY = y
}

func (dev *rgbDevice) GetTextArea(text string) (x1, y1, x2, y2 int) {
	x1 = 0
	y1 = 0
	x2 = 0
	y2 = 0
	switch dev.fontType {
	case BITMAP_FONT:
		x1, y1, x2, y2 = dev.getBitmapFontTextArea(text)
	}
	return
}

func (dev *rgbDevice) drawBitmapChar(char byte) {
	glyph := dev.bitmapFont.Glyphs[char-0x20]
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
			x := dev.cursorX + w + glyph.XOffset
			y := dev.cursorY + h + glyph.YOffset
			dev.pixeldev.Pixel(x, y, color)
		}
	}
	dev.cursorX += glyph.XAdvance
}

func (dev *rgbDevice) getBitmapFontTextArea(text string) (int, int, int, int) {
	bytes := []byte(text)
	ymax := 0
	ymin := 0
	x := 0
	for i := 0; i < len(bytes); i++ {
		glyph := dev.bitmapFont.Glyphs[bytes[i]-0x20]
		x += glyph.XAdvance
		y := glyph.YOffset + glyph.Height
		if y > ymax {
			ymax = y
		}
		if glyph.YOffset < ymin {
			ymin = glyph.YOffset
		}
	}
	return 0, ymin, x, ymax
}
