package main

import (
	"fmt"
	"testing"
)

/*
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
	equals(t, newVertex(0, 1, a.Theta, a), c)

	// Pure x-direction
	a, b = newVertex(0, 0, 0, nil), newVertex(100, 0, 0, nil)
	c = smallDistanceAlong(a, b, 50, false)
	equals(t, newVertex(50, 0, a.Theta, a), c)

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
*/

func TestEuler(t *testing.T) {
	init := Point{1, 0, 0, 0, 0}
	next := euler(init, 0.01, 1, 1)
	fmt.Println(next)
}

func TestForwardSim(t *testing.T) {
	init := Point{1, 0, 0, 0, 0}
	end := Point{10, 0, 0, 3, 0}
	start := &Vertex{Point: init}
	goal := &Vertex{Point: end}

	cSpace := &ConfigSpace{
		XMin: 0, XMax: 100,
		YMin: 0, YMax: 100,
		VMin: -5, VMax: 5,
		WMin: -5, WMax: 5,
		AMin: -5, AMax: 5,
		GammaMin: -5, GammaMax: 5,
	}
	safe := func(p *Point) bool { return true }
	vert, edge, ok := forwardSim(start, goal, 1, 0.5, safe, cSpace)
	fmt.Println(vert)
	fmt.Println(edge)
	fmt.Println(ok)
}
