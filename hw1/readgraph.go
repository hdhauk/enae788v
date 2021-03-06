package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

type problem struct {
	id       int
	vertices map[int]*Vertex
	startID  int
	goalID   int
}

func readProblems(filePath string) ([]*problem, error) {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read file")
	}

	dir, err := filepath.Abs(filepath.Dir(filePath))
	if err != nil {
		return nil, errors.Wrap(err, "could not determine file path")
	}

	chunks := splitByEmptyNewline(string(dat))

	problemIDRegex := regexp.MustCompile(`Problem\s(\d+):`)
	vertexFileRegex := regexp.MustCompile(`node\sfile:\s*(.*)`)
	edgeFileRegex := regexp.MustCompile(`edge\sfile:\s*(.*)`)
	startIDRegex := regexp.MustCompile(`start\snode ID:\s*(\d*)`)
	goalIDRegex := regexp.MustCompile(`goal\snode ID:\s*(\d*)`)

	var problems []*problem
	for _, chunk := range chunks {
		matches := problemIDRegex.FindAllStringSubmatch(chunk, 1)
		if len(matches) < 1 {
			continue
		}
		id, err := strconv.Atoi(matches[0][1])
		if err != nil {
			log.Println("problem id must be an integer")
		}

		matches = vertexFileRegex.FindAllStringSubmatch(chunk, 1)
		if len(matches) < 1 {
			log.Println("could not find node file path")
			continue
		}
		nodeFilePath := matches[0][1]
		vertices, err := readVertices(dir + "/" + nodeFilePath)
		if err != nil {
			log.Printf("failed to read nodes %v\n", err)
			continue
		}

		matches = edgeFileRegex.FindAllStringSubmatch(chunk, 1)
		if len(matches) < 1 {
			log.Println("could not find node file path")
			continue
		}
		edgeFilePath := matches[0][1]
		edges, err := readEdges(dir + "/" + edgeFilePath)
		if err != nil {
			log.Printf("failed to read edges %v\n", err)
			continue
		}
		addNeighborsToVertices(vertices, edges)

		matches = startIDRegex.FindAllStringSubmatch(chunk, 1)
		if len(matches) < 1 {
			log.Println("could not find start id")
			continue
		}
		startID, err := strconv.Atoi(matches[0][1])
		if err != nil {
			log.Println("could not parse start ID")
		}

		matches = goalIDRegex.FindAllStringSubmatch(chunk, 1)
		if len(matches) < 1 {
			log.Println("could not find goal id")
			continue
		}
		goalID, err := strconv.Atoi(matches[0][1])
		if err != nil {
			log.Println("could not parse goal ID")
		}

		p := &problem{
			id:       id,
			startID:  startID,
			goalID:   goalID,
			vertices: vertices,
		}
		problems = append(problems, p)
	}

	if len(problems) < 1 {
		return nil, errors.Errorf("no problems in file %s", filePath)
	}

	return problems, nil
}

func readVertices(filePath string) (map[int]*Vertex, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open file %s", filePath)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan() // First line is just number of elements...doesn't want this

	nodeRegex := regexp.MustCompile(`^(\d+),\s*(\d+.\d+),\s*(\d+.\d+)`)

	vertices := make(map[int]*Vertex)
	for scanner.Scan() {
		matches := nodeRegex.FindAllStringSubmatch(scanner.Text(), 1)
		if len(matches) < 1 {
			log.Println("failed to parse line")
			continue
		}
		id, err := strconv.Atoi(matches[0][1])
		if err != nil {
			log.Println("failed to parse line")
			continue
		}
		x, err := strconv.ParseFloat(matches[0][2], 64)
		if err != nil {
			log.Println("failed to parse line")
			continue
		}
		y, err := strconv.ParseFloat(matches[0][3], 64)
		if err != nil {
			log.Println("failed to parse line")
			continue
		}
		v := &Vertex{
			id:        id,
			x:         x,
			y:         y,
			neighbors: make(map[int]float64),
		}
		vertices[id] = v
	}

	if len(vertices) < 1 {
		return nil, errors.New("no vertices in file")
	}
	return vertices, nil
}

type edge struct {
	tail, head int
	distance   float64
}

func addNeighborsToVertices(vertices map[int]*Vertex, edges []edge) {
	for _, v := range vertices {
		for _, e := range edges {
			if e.tail == v.id {
				v.neighbors[e.head] = e.distance
			}
		}
	}
}

func readEdges(filePath string) ([]edge, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open edge file %s", filePath)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan() // First line is just number of elements...doesn't want this

	edgeRegex := regexp.MustCompile(`^(\d+),\s*(\d+),\s*(\d+.\d+)`)

	var edges []edge
	for scanner.Scan() {
		matches := edgeRegex.FindAllStringSubmatch(scanner.Text(), 1)
		if len(matches) < 1 {
			log.Println("failed to parse line")
			continue
		}
		start, err := strconv.Atoi(matches[0][1])
		if err != nil {
			log.Println("failed to parse line")
			continue
		}
		end, err := strconv.Atoi(matches[0][2])
		if err != nil {
			log.Println("failed to parse line")
			continue
		}
		dist, err := strconv.ParseFloat(matches[0][3], 64)
		if err != nil {
			log.Println("failed to parse line")
			continue
		}
		edges = append(edges, edge{start, end, dist})
	}

	if len(edges) < 1 {
		return nil, errors.New("no edges in file")
	}

	return edges, nil
}

func splitByEmptyNewline(str string) []string {
	strNormalized := regexp.
		MustCompile("\r\n").
		ReplaceAllString(str, "\n")

	return regexp.
		MustCompile(`\n\s*\n`).
		Split(strNormalized, -1)

}
