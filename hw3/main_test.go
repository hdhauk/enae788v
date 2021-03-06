package main

import (
	"math"
	"testing"
)

func TestRobotPointToGlobal(t *testing.T) {

	var tests = []struct {
		name        string
		robotCenter Point
		robotPoint  Point
		exp         Point
	}{
		{"0 degree", Point{10, 10, 0}, Point{1, 1, 0}, Point{11, 11, 0}},
		{"90 degree", Point{10, 10, math.Pi / 2}, Point{1, 1, 0}, Point{9, 11, 0}},
		{"-180 degree", Point{10, 10, -math.Pi}, Point{1, 1, 0}, Point{9, 9, 0}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := robotPointGlobal(tc.robotCenter, tc.robotPoint)
			equals(t, tc.exp.X, got.X)
			equals(t, tc.exp.Y, got.Y)
		})
	}

}

func TestPointsAlongPath(t *testing.T) {

	var tests = []struct {
		name    string
		start   Point
		end     Point
		epsilon float64
		exp     []Point
	}{
		{"vertical_line", Point{0, 0, 0}, Point{0, 5, 0}, 1.0,
			[]Point{
				Point{0, 0, math.Pi / 2},
				Point{0, 1, math.Pi / 2},
				Point{0, 2, math.Pi / 2},
				Point{0, 3, math.Pi / 2},
				Point{0, 4, math.Pi / 2},
				Point{0, 5, math.Pi / 2},
			},
		},
		{"horizontal_line", Point{0, 0, 0}, Point{2.5, 0, 0}, 0.5,
			[]Point{
				Point{0, 0, 0},
				Point{0.5, 0, 0},
				Point{1, 0, 0},
				Point{1.5, 0, 0},
				Point{2.0, 0, 0},
				Point{2.5, 0, 0},
			},
		},
		{"horizontal_line_different_angles", Point{0, 0, math.Pi / 2}, Point{2.5, 0, -math.Pi / 2}, 0.5,
			[]Point{
				Point{0, 0, 0},
				Point{0.5, 0, 0},
				Point{1, 0, 0},
				Point{1.5, 0, 0},
				Point{2.0, 0, 0},
				Point{2.5, 0, 0},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := getPointsAlongPath(tc.start, tc.end, tc.epsilon)
			equals(t, len(tc.exp), len(got))
			for i := 0; i < len(got); i++ {
				equals(t, tc.exp[i].X, got[i].X)
				equals(t, tc.exp[i].Y, got[i].Y)
			}
		})
	}

}

func TestAngleDiff(t *testing.T) {
	var tests = []struct {
		name     string
		start    float64
		end      float64
		expAngle float64
		expDir   float64
	}{
		{"90 degree", 0, math.Pi / 2, math.Pi / 2, 1},
		{"-90 degree", 0, -math.Pi / 2, math.Pi / 2, -1},
		{"-90 degree, 45 to neg45", math.Pi / 4, -math.Pi / 4, math.Pi / 2, -1},
		{"-90 degree, 45 to 45", -math.Pi / 4, math.Pi / 4, math.Pi / 2, 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			angle, dir := angleDiff(tc.start, tc.end)
			equals(t, tc.expAngle, angle)
			equals(t, tc.expDir, dir)

		})
	}
}
