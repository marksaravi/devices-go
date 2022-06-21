package rgb

import (
	"errors"

	"github.com/marksaravi/devices-go/colors/rgb565"
)

type RGB interface{}

func ToRGB565(rgb RGB) (rgb565.RGB565, error) {
	switch v := rgb.(type) {
	case rgb565.RGB565:
		return v, nil
	default:
		return rgb565.BLACK, errors.New("rgb565 color type mistmatch")
	}
}
