package main

import (
	"math"
	"math/rand"
	"testing"
)

func TestRandomSample(t *testing.T) {

	rand.Seed(69)
	cSpace := ConfigSpace{0, 100, 0, 100}

	for i := 0; i < 1000000; i++ {
		p := randomSample(cSpace)
		assert(t, p.X > 0.0 || p.X < 100.0, "point outside of range")
		assert(t, p.Y > 0.0 || p.Y < 100.0, "point outside of range")
	}
}

func TestSmallDistanceAlong(t *testing.T) {
	// Pure y-direction
	a, b := newVertex(0, 0, nil), newVertex(0, 10, nil)
	c := smallDistanceAlong(a, b, 1)
	equals(t, newVertex(0, 1, a), c)

	// Pure x-direction
	a, b = newVertex(0, 0, nil), newVertex(100, 0, nil)
	c = smallDistanceAlong(a, b, 50)
	equals(t, newVertex(50, 0, a), c)

	// Pure 45-direction
	a, b = newVertex(0, 0, nil), newVertex(100, 100, nil)
	c = smallDistanceAlong(a, b, 50)
	equals(t, 50.0, math.Sqrt(math.Pow(c.X, 2)+math.Pow(c.Y, 2)))
	assert(t, c.X == c.Y, "should be same length")
}

func TestNear(t *testing.T) {
	p := newVertex(10, 10, nil)

	c := Circle{10, 15, 6}
	assert(t, near(p, c), "should be near")

	c = Circle{10, 15, 4.9999}
	assert(t, !near(p, c), "should not be near")

}
