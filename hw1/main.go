package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	problemsPath := os.Args[1]

	problems, err := readProblems(problemsPath)
	if err != nil {
		log.Fatal(err)
	}

	p1 := problems[0]
	var h heuristic
	h = func(u, goal *Vertex) float64 { return 0.0 }
	path, err := aStar(p1.vertices, p1.startID, p1.goalID, h)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range path {
		fmt.Printf("%3d: ID: %2d\n", k, v.id)
	}

}
