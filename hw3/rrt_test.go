package main

import (
	"math"
	"math/rand"
	"testing"
)

func TestRandomSample(t *testing.T) {

	rand.Seed(69)
	cSpace := ConfigSpace{0, 100, 0, 100}

	for i := 0; i < 1000; i++ {
		p := randomSample(cSpace)
		assert(t, p.X > 0.0 && p.X < 100.0, "point outside of range")
		assert(t, p.Y > 0.0 && p.Y < 100.0, "point outside of range")
	}
}

func TestSmallDistanceAlong(t *testing.T) {
	// Pure y-direction
	a, b := newVertex(0, 0, 0, nil), newVertex(0, 10, 0, nil)
	c := smallDistanceAlong(a, b, 1, false)
	equals(t, newVertex(0, 1, math.Pi/2, a), c)

	// Pure x-direction
	a, b = newVertex(0, 0, 0, nil), newVertex(100, 0, 0, nil)
	c = smallDistanceAlong(a, b, 50, false)
	equals(t, newVertex(50, 0, 0, a), c)

	// Pure 45-direction
	a, b = newVertex(0, 0, 0, nil), newVertex(100, 100, 0, nil)
	c = smallDistanceAlong(a, b, 50, false)
	equals(t, 50.0, math.Sqrt(math.Pow(c.X, 2)+math.Pow(c.Y, 2)))
	assert(t, c.X == c.Y, "should be same length")
}

func TestNear(t *testing.T) {
	p := newVertex(10, 10, 0, nil)

	c := Circle{10, 15, 6}
	assert(t, near(p, c), "should be near")

	c = Circle{10, 15, 4.9999}
	assert(t, !near(p, c), "should not be near")

}

// func TestSafeFunc(t *testing.T) {
// 	obstacles := []Circle{
// 		Circle{5, 5, 3},
// 		Circle{63, 25, 8},
// 		Circle{53, 25, 8},
// 	}
// 	cSpace := ConfigSpace{0, 100, 0, 100}

// 	safe := getSafeFunc(obstacles, cSpace, )

// 	var tests = []struct {
// 		name string
// 		a    *Vertex
// 		b    *Vertex
// 		exp  bool
// 	}{
// 		{"through obstacle, b on perimeter", newVertex(4, 9, nil), newVertex(8, 5, nil), false},
// 		{"b on perimeter", newVertex(4, 9, nil), newVertex(5, 8, nil), true},
// 		{"b just inside obstacle", newVertex(4, 9, nil), newVertex(5, 7.99, nil), false},
// 		{"a and b on each side", newVertex(1, 5, nil), newVertex(9, 5, nil), false},
// 		{"tangent on top", newVertex(1, 8, nil), newVertex(9, 8, nil), true},
// 		{"barely not tangent", newVertex(1, 8, nil), newVertex(9, 8.1, nil), true},
// 		{"above directly toward center", newVertex(5, 11, nil), newVertex(5, 9, nil), true},
// 		{"above directly away from center", newVertex(5, 9, nil), newVertex(5, 11, nil), true},
// 		{"to the right directly toward center", newVertex(11, 5, nil), newVertex(9, 5, nil), true},
// 		{"to the right almost directly toward center", newVertex(11, 5, nil), newVertex(9, 4.999, nil), true},
// 		{"to the right angled down", newVertex(12, 5, nil), newVertex(9, 3, nil), true},
// 		{"45 degrees down toward center", newVertex(10, 10, nil), newVertex(10, 10, nil), true},
// 		{"special case", newVertex(11, 5, nil), newVertex(10.1, 4.9, nil), true},
// 		{"over angled downward", newVertex(59.1185, 34.2603, nil), newVertex(67.4492, 16.078, nil), false},
// 		{"over angled upward", newVertex(52.00, 16.60, nil), newVertex(51.15, 36.58, nil), false},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			equals(t, tc.exp, safe(tc.a, tc.b))
// 		})
// 	}

// }
