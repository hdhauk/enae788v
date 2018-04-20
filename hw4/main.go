package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/ungerik/go3d/float64/vec2"
)

func main() {
	configPath := flag.String("c", "problems.json", "config file")
	pIndex := flag.Int("p", 0, "which problem in config file to solve (0-indexed)")
	flag.Parse()

	// Read in config.
	configFile, err := os.Open(*configPath)
	defer configFile.Close()
	if err != nil {
		log.Fatalf("could not open config file: %v\n", err)
	}

	config, err := parseConfig(configFile)
	if err != nil {
		log.Fatalf("could not parse config file: %v\n", err)
	}

	// Sanity check problem number.
	if len(config.Problems)-1 < *pIndex || *pIndex < 0 {
		log.Fatalln("invalid problem number")
	}

	// Read in obstacles.
	obstacleFile, err := os.Open(config.ObstaclesPath)
	if err != nil {
		log.Fatalf("could not open obstacle file: %v\n", err)
	}
	defer obstacleFile.Close()
	obstacles, err := readObstacles(obstacleFile)
	if err != nil {
		log.Fatalf("could not read obstacles from file: %v\n", err)
	}

	// Read in robot.
	robotFile, err := os.Open(config.RobotPath)
	if err != nil {
		log.Fatalf("could not open robot file: %v\n", err)
	}
	defer robotFile.Close()
	robot, err := readRobot(robotFile)
	if err != nil {
		log.Fatalf("could not read robot from file: %v\n", err)
	}

	// Solve problem.
	p := config.Problems[*pIndex]
	safe := getSafeFunc(obstacles, config.ConfigSpace, robot)
	seed := time.Now().UnixNano()
	path, tree, err := RRT(obstacles, p, config.ConfigSpace, safe, seed)
	if err != nil {
		log.Fatalf("RRT failed during execution: %v\n", err)
	}

	// printing
	fmt.Printf("start=[%.4f,%.4f] goal=[%.4f,%.4f,%.4f]\n\n", p.Start.X, p.Start.Y, p.Goal.X, p.Goal.Y, p.Goal.R)

	fmt.Println("START_PATH")
	for _, v := range path {
		fmt.Printf("%.4f, %.4f, %.4f, %.4f, %.4f\n", v.head.X, v.head.Y, v.tail.X, v.tail.Y, v.head.Theta)
	}
	fmt.Println("END_PATH")
	fmt.Println()

	fmt.Println("START_TREE")
	for _, v := range tree {
		fmt.Printf("%.4f, %.4f, %.4f, %.4f\n", v.head.X, v.head.Y, v.tail.X, v.tail.Y)
	}
	fmt.Println("END_TREE")
}

func getSafeFunc(obstacles []Circle, cSpace ConfigSpace, bot Robot) SafeFunc {
	illegalPoint := func(p Point) bool {
		inConfigSpace := (cSpace.XMin < p.X && p.X < cSpace.XMax) && (cSpace.YMin < p.Y && p.Y < cSpace.YMax)
		if !inConfigSpace {
			return true
		}

		for _, circle := range obstacles {
			if near(newVertex(p.X, p.Y, 0, nil), circle) {
				return true
			}
		}
		return false
	}

	return func(v, w *Vertex) bool {
		// Get waypoints with epsilon spacing between v and w.
		pointsAlongPath := getPointsAlongPath(Point{v.X, v.Y, 0}, Point{w.X, w.Y, 0}, 0.5)

		// collitions := make(chan bool)
		for _, waypoint := range pointsAlongPath {
			// Check all points of robot.
			// TODO: Do this can be done concurrently.
			for _, robotPoint := range bot {
				globalPoint := robotPointGlobal(waypoint, robotPoint)
				if illegalPoint(globalPoint) {
					return false
				}
			}
		}
		return true
	}
}

func robotPointGlobal(base, offset Point) Point {
	b := vec2.T{offset.X, offset.Y}
	b.Rotate(base.Theta)
	a := vec2.T{base.X, base.Y}
	c := a.Add(&b)

	return Point{c[0], c[1], 0}
}

func getPointsAlongPath(start, end Point, epsilon float64) []Point {
	a := vec2.T{start.X, start.Y}
	b := vec2.T{end.X, end.Y}

	a2b := b.Sub(&a)
	a2bNorm := a2b.Normalized()

	distance := a2b.Length()
	numPoints := math.Floor(distance / epsilon)

	diffTheta, dir := angleDiff(start.Theta, end.Theta)
	deltaTheta := diffTheta / (numPoints - 1)

	pointsAlong := []Point{Point{a[0], a[1], a2b.Angle()}}
	for i := 0.0; i < numPoints; i++ {
		offset := a2bNorm.Scaled(epsilon)
		waypoint := a.Add(&offset)
		theta := start.Theta + deltaTheta*dir*i
		pointsAlong = append(pointsAlong, Point{waypoint[0], waypoint[1], theta})
	}

	return pointsAlong
}

// angleDiff returns the absolute value of the difference in angle, and if that is in
// positive or negative direction.
func angleDiff(start, end float64) (float64, float64) {
	mod := func(a, n int) int {
		return a - int(math.Floor(float64(a/n)))*n
	}

	// use degrees
	angleDiff := int(math.Round(end*180/math.Pi - start*180/math.Pi))
	angleDiff = mod(angleDiff+180, 360) - 180
	angleDiff64 := float64(angleDiff)
	angleDiffRad := angleDiff64 * math.Pi / 180.0

	if angleDiffRad < 0 {
		return math.Abs(angleDiffRad), -1.0

	}
	return math.Abs(angleDiffRad), 1.0
}
