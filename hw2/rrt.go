package main

import (
	"encoding/csv"
	"io"
	"log"
	"strconv"

	"github.com/pkg/errors"
)

func main() {

}

type obstacle struct {
	x, y, radius int
}

func readObstacles(reader io.Reader) ([]obstacle, error) {
	r := csv.NewReader(reader)
	d, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var obstacles []obstacle
	for _, v := range d {
		x, err := strconv.Atoi(v[0])
		y, err := strconv.Atoi(v[1])
		r, err := strconv.Atoi(v[2])
		if err != nil {
			return nil, errors.Wrap(err, "non-integer value")
		}

		obstacles = append(obstacles, obstacle{x, y, r})
	}

	return obstacles, nil
}
