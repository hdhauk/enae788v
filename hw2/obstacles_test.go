package main

import (
	"fmt"
	"os"
	"testing"
)

func TestReadObstacles(t *testing.T) {
	file, err := os.Open("obstacles.txt")
	if err != nil {
		t.Fatalf("could not open obstacles: %+v", err)
	}

	o, err := readObstacles(file)
	if err != nil {
		t.Errorf("failed to read obstacles.txt: %+v", err)
	}

	for _, v := range o {
		fmt.Printf("%+v\n", v)
	}

}
