package display

import (
	"math"
	"testing"
)

func TestIsPointInsideArc(t *testing.T) {
}

func TestGetSectors(t *testing.T) {
	toRad := func(a float64) float64 {
		return math.Pi / 180 * a
	}

	angles := [][2]int{{0, 15}, {55, 200}, {330, 359}, {330, 0}, {330, 15}}
	want := [][2]int{
		{0, 0},
		{0, 2},
		{3, 3},
		{3, 0},
		{3, 0},
	}
	for i := 0; i < len(angles); i++ {
		sa := toRad(float64(angles[i][0]))
		ea := toRad(float64(angles[i][1]))
		s1, s2 := getSectors(sa, ea)
		if s1 != want[i][0] || s2 != want[i][1] {
			t.Errorf("%d,%d sectors are not %d,%d", angles[i][0], angles[i][1], s1, s2)
		}
	}
}
