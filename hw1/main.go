package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func main() {
	searchTreePath := flag.String("tree", "search_tree.txt", "path for search tree file")
	shortestPathPath := flag.String("path", "output_path.txt", "path for shortest path path")
	problemSet := flag.Int("problem", 1, "number identifier for problem set in provided problem file")
	dijkstra := flag.Bool("dijk", false, "set this flag to not set heuristic return 0, effectively rendering the algorithm equal to Dijkstra")
	silent := flag.Bool("silent", false, "turn for file outputs, will still report length of shortest path")
	help := flag.Bool("h", false, "show help")
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	if len((os.Args)) < 2 {
		fmt.Println("please provide a problem file")
		printHelp()
		return
	}

	log.SetFlags(log.Ltime | log.Lshortfile)

	problemsPath := os.Args[len(os.Args)-1]

	problems, err := readProblems(problemsPath)
	if err != nil {
		log.Fatal(err)
	}

	p1 := problems[*problemSet-1]
	h := cartesianDistance
	if *dijkstra {
		h = func(u, goal *Vertex) float64 { return 0.0 }
	}
	results, err := aStar(p1.vertices, p1.startID, p1.goalID, h)
	if err != nil {
		red := color.New(color.FgRed).FprintfFunc()
		red(os.Stderr, "%s", err)
		if *silent {
			return
		}
		if err := writeSearchTree(*searchTreePath, results.searchTree); err != nil {
			log.Fatal(err)
		}
		return
	}

	green := color.New(color.FgGreen).PrintfFunc()
	green("Found shortest path with distance %.3f\n", results.pathCost)

	if *silent {
		return
	}
	if err := writeSearchTree(*searchTreePath, results.searchTree); err != nil {
		log.Fatal(err)
	}

	if err := writePath(*shortestPathPath, results.path); err != nil {
		log.Fatal(err)
	}
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

func printHelp() {
	fmt.Print(`-dijk              set this flag to not set heuristic return 0, effectively rendering the algorithm equal to Dijkstra
-h                 show help
-path     string   path for shortest path path (default "output_path.txt")
-problem  int      number identifier for problem set in provided problem file (default 1)
-tree     string   path for search tree file (default "search_tree.txt")
-silent            turn for file outputs, will still report length of shortest path`)
}
