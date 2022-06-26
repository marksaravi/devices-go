package display

import (
	"errors"

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

func (dev *rgbDevice) SetLineHeight(height int) {
	dev.lineHeight = height
}

func (dev *rgbDevice) WriteChar(char byte) error {
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
		dev.WriteChar(text[i])
	}
}

func (dev *rgbDevice) MoveCursor(x, y int) {
	dev.cursorX = x
	dev.cursorY = y
}

func (dev *rgbDevice) drawBitmapChar(char byte) error {
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
			dev.pixeldev.Pixel(dev.cursorX+w+glyph.XOffset, dev.cursorY+h+glyph.YOffset+dev.lineHeight, color)
		}
	}
	dev.cursorX += glyph.XAdvance
	return nil
}

func (dev *rgbDevice) initBitmapFonts(bitmapfont fonts.BitmapFont) {
	dev.fontType = BITMAP_FONT
	dev.bitmapFont = bitmapfont

	glyph := dev.bitmapFont.Glyphs[0x20]
	dev.charAdvanceX = 0
	if dev.lineHeight == 0 {
		dev.lineHeight = glyph.YOffset
	}
}
