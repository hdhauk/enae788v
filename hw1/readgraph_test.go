package main

import (
	"testing"
)

func TestReadProblems(t *testing.T) {
	problems, err := readProblems("problems/problems.txt")
	if err != nil {
		t.Error(err)
	}

	if len(problems) != 5 {
		t.Errorf("expected 5 problems, got %d", len(problems))
	}

	p1 := problems[0]
	equals(t, p1.id, 1)
	equals(t, p1.startID, 1)
	equals(t, p1.goalID, 10)
	equals(t, len(p1.vertices), 100)

}

func TestReadVertices(t *testing.T) {
	vertices, err := readVertices("./problems/nodes_1.txt")
	if err != nil {
		t.Error(err)
	}

	if len(vertices) != 100 {
		t.Errorf("expected 100 vertices, got %d", len(vertices))
	}

}

func TestReadEdges(t *testing.T) {
	edges, err := readEdges("./problems/edges_1.txt")
	if err != nil {
		t.Error(err)
	}

	if len(edges) != 1000 {
		t.Errorf("expected 1000 edges, got %d", len(edges))
	}
}
