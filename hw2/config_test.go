package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func TestParseConfig(t *testing.T) {
	want := Config{
		ObstaclesPath: "obstacles.txt",
		ConfigSpace:   ConfigSpace{0, 100.0, 0, 100.0},
		Problems: []Problem{
			Problem{
				Start:   Point{75.0, 85.0},
				Goal:    Circle{100.0, 0.0, 20.0},
				Epsilon: 10,
			},
		},
	}

	s := strings.NewReader(`{
		"obstacles": "obstacles.txt",
		"config_space": {
			"x_min": 0,
			"x_max": 100,
			"y_min": 0,
			"y_max": 100
		},
	
		"problems": [{
			"start": {
                "x": 75,
                "y": 85
            },
            "goal_region": {
                "x": 100,
                "y": 0,
                "r": 20
            },
            "epsilon": 10
		}]
	}`)

	c, err := parseConfig(s)
	ok(t, err)
	equals(t, &want, c)
	assert(t, len(c.Problems) == 1, "expected to have 1 problem")

}

func TestReadObstacles(t *testing.T) {
	file, err := os.Open("obstacles.txt")
	if err != nil {
		t.Fatalf("could not open obstacles: %+v", err)
	}

	_, err = readObstacles(file)
	ok(t, err)

}
