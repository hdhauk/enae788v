package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
)

func main() {
	problemsPath := os.Args[1]

	problems, err := readProblems(problemsPath)
	if err != nil {
		log.Fatal(err)
	}

	p1 := problems[0]
	// var h heuristic
	// h = func(u, goal *Vertex) float64 { return 0.0 }
	// results, err := aStar(p1.vertices, p1.startID, p1.goalID, h)
	results, err := aStar(p1.vertices, p1.startID, p1.goalID, cartesianDistance)
	if err != nil {
		log.Fatal(err)
	}

	if err := writeSearchTree("search_tree.txt", results.searchTree); err != nil {
		log.Fatal(err)
	}

	if err := writePath("output_path.txt", results.path); err != nil {
		log.Fatal(err)
	}

	// {
	// 	file, err := os.OpenFile("search_tree.txt", os.O_RDWR, 0666)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer file.Close()

	// 	for i := 1; i < len(results.searchTree); i++ {
	// 		a := results.searchTree[i]
	// 		fmt.Fprintf(file, "%d, %f, %f, %d, %f, %f\n", a.id, a.x, a.y, a.parent.id, a.parent.x, a.parent.y)
	// 		// fmt.Printf("%3d: ID: %2d\n", k, v.id)
	// 	}
	// }
	// {
	// 	file, err := os.OpenFile("output_path.txt", os.O_RDWR, 0666)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer file.Close()

	// 	for k, v := range results.path {
	// 		fmt.Fprintf(file, "%d, %f, %f\n", v.id, v.x, v.y)
	// 		fmt.Printf("%3d: ID: %2d\n", k, v.id)
	// 	}
	// }

}

func writeSearchTree(path string, tree []*Vertex) error {

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "could not create file")
	}
	defer file.Close()

	for i := 1; i < len(tree); i++ {
		v := tree[i]
		fmt.Fprintf(file, "%d, %f, %f, %d, %f, %f\n", v.id, v.x, v.y, v.parent.id, v.parent.x, v.parent.y)
	}

	return nil

}

func writePath(filePath string, path []*Vertex) error {

	file, err := os.Create(filePath)
	if err != nil {
		return errors.Wrap(err, "could not create file")
	}
	defer file.Close()

	for _, v := range path {
		fmt.Fprintf(file, "%d, %f, %f\n", v.id, v.x, v.y)
	}

	return nil
}
