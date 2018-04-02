package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/pkg/errors"
)

type Circle struct {
	X, Y, R float64
}

func (c Circle) String() string {
	return fmt.Sprintf("(%.2f, %.2f, r=%.2f)", c.X, c.Y, c.R)
}

type Point struct {
	X, Y float64
}

type Problem struct {
	Start           Point
	Goal            Circle `json:"goal_region"`
	Epsilon         float64
	AllowSmallSteps bool `json:"allow_steps_smaller_than_epsilon"`
}
type ConfigSpace struct {
	XMin float64 `json:"x_min"`
	XMax float64 `json:"x_max"`
	YMin float64 `json:"y_min"`
	YMax float64 `json:"y_max"`
}

type Config struct {
	ObstaclesPath string      `json:"obstacles"`
	ConfigSpace   ConfigSpace `json:"config_space"`
	Problems      []Problem
}

func parseConfig(reader io.Reader) (*Config, error) {
	var c Config
	if err := json.NewDecoder(reader).Decode(&c); err != nil {
		return nil, errors.Wrap(err, "could not decode JSON")
	}
	return &c, nil
}

func readObstacles(reader io.Reader) ([]Circle, error) {
	r := csv.NewReader(reader)
	d, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var obstacles []Circle
	for _, v := range d {
		x, err := strconv.ParseFloat(v[0], 64)
		y, err := strconv.ParseFloat(v[1], 64)
		r, err := strconv.ParseFloat(v[2], 64)
		if err != nil {
			return nil, errors.Wrap(err, "non-integer value")
		}

		obstacles = append(obstacles, Circle{x, y, r})
	}

	return obstacles, nil
}
