package display

import (
	"math"
	"testing"
)

func TestIsPointInsideArc(t *testing.T) {
	toRad := func(a float64) float64 {
		return math.Pi / 180 * a
	}
	const RADIUS float64 = 10
	XY := func(angle float64) (float64, float64) {
		x := RADIUS * math.Cos(toRad(angle))
		y := math.Sqrt(RADIUS*RADIUS - x*x)
		return x, y
	}

	sAngle := float64(0)
	eAngle := float64(30)
	x, y := XY(15)
	if !isInside(x, y, toRad(sAngle), toRad(eAngle)) {
		t.Errorf("(%f, %f) are not in %f, %f\n", x, y, sAngle, eAngle)
	}
}
