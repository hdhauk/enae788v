package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Circle defines a ball in 2D space.
type Circle struct {
	X, Y, R float64
}

func (c Circle) String() string {
	return fmt.Sprintf("(%.2f, %.2f, r=%.2f)", c.X, c.Y, c.R)
}

// Problem defines a specific path planning problem with a given config space.
type Problem struct {
	Name            string
	Start           Point
	Goal            Circle `json:"goal_region"`
	Epsilon         float64
	Delta           float64
	AllowSmallSteps bool `json:"allow_steps_smaller_than_epsilon"`
}

// ConfigSpace should be renamed to workspace....
type ConfigSpace struct {
	XMin     float64 `json:"x_min"`
	XMax     float64 `json:"x_max"`
	YMin     float64 `json:"y_min"`
	YMax     float64 `json:"y_max"`
	VMax     float64 `json:"v_max"`
	VMin     float64 `json:"v_min"`
	WMax     float64 `json:"w_max"`
	WMin     float64 `json:"w_min"`
	AMax     float64 `json:"a_max"`
	AMin     float64 `json:"a_min"`
	GammaMax float64 `json:"gamma_max"`
	GammaMin float64 `json:"gamma_min"`
}

// Config is the go struct equivalent of the .json file describing the problems.
type Config struct {
	ObstaclesPath string      `json:"obstacles"`
	RobotPath     string      `json:"robot_path"`
	ConfigSpace   ConfigSpace `json:"config_space"`
	Problems      []Problem
}

// Robot is simply a set of points defining the edges of the robot.
type Robot []Point

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
		return nil, errors.Wrap(err, "could not read csv")
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

func readRobot(reader io.Reader) (Robot, error) {
	r := csv.NewReader(reader)
	records, err := r.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "could not read csv values")
	}

	var bot Robot
	for _, v := range records {
		x, err := strconv.ParseFloat(strings.TrimSpace(v[0]), 64)
		y, err := strconv.ParseFloat(strings.TrimSpace(v[1]), 64)
		if err != nil {
			return nil, errors.Wrap(err, "non-float value in csv")
		}
		bot = append(bot, Point{X: x, Y: y})
	}
	return bot, nil
}
