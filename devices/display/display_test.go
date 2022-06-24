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

func TestSector(t *testing.T) {
	points := [][2]float64{{1, 0}, {0.7, 0.7}, {0, 1}, {-0.7, 0.7}, {-1, 0}, {-0.7, -0.7}, {0, -1}, {0.7, -0.7}}
	want := []int{0, 0, 1, 1, 2, 2, 3, 3}
	for i := 0; i < len(points); i++ {
		x := points[i][0]
		y := points[i][1]
		got := getSector(x, y)
		if got != int(want[i]) {
			t.Errorf("point (%f,%f) is expected to be in sector %d but was in %d\n", x, y, want[i], got)
		}
	}
}
