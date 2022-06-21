package rgb

import (
	"errors"

	"github.com/marksaravi/devices-go/colors/rgb565"
)

type RGB interface{}

func ToRGB565(color RGB) (rgb565.RGB565, error) {
	c, ok := color.(rgb565.RGB565)
	if !ok {
		return rgb565.BLACK, errors.New("rgb565 color type mistmatch")
	}
	return c, nil
}
