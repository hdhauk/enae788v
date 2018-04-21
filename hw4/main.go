package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"

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
	// seed := time.Now().UnixNano()
	seed := int64(11)
	path, tree, err := RRT(obstacles, p, &config.ConfigSpace, safe, seed)
	if err != nil {
		log.Fatalf("RRT failed during execution: %v\n", err)
	}

	// printing
	fmt.Printf("start=[%.4f,%.4f] goal=[%.4f,%.4f,%.4f]\n\n", p.Start.X, p.Start.Y, p.Goal.X, p.Goal.Y, p.Goal.R)
	_, _ = path, tree
	printPath(path)
	printTree(tree)

}

func printPath(path []*Point) {
	// for i, p := range path {
	// 	fmt.Printf("#%d (%.1f,%.1f)\n", i, p.X, p.Y)

	// }

	fmt.Println("START_PATH")
	var head, tail *Point
	for i := len(path) - 1; i > 0; i-- {
		head = path[i]
		tail = path[i-1]
		fmt.Printf("%.4f, %.4f, %.4f, %.4f, %.4f\n",
			head.X, head.Y, tail.X, tail.Y, head.Theta)
	}

	fmt.Println("END_PATH")
	fmt.Println()
}

func printTree(tree []*Edge) {
	fmt.Println("START_TREE")
	for _, v := range tree {
		fmt.Printf("%.4f, %.4f, %.4f, %.4f\n", v.head.X, v.head.Y, v.tail.X, v.tail.Y)
	}
	fmt.Println("END_TREE")
}

func getSafeFunc(obstacles []Circle, cSpace ConfigSpace, bot Robot) SafeFunc {
	// fmt.Printf("xmin:%.2f, xmax:%.2f\n", cSpace.XMin, cSpace.XMax)
	// fmt.Printf("ymin:%.2f, ymax:%.2f\n", cSpace.YMin, cSpace.YMax)
	// fmt.Printf("vmin:%.2f, vmax:%.2f\n", cSpace.VMin, cSpace.VMax)
	// fmt.Printf("wmin:%.2f, wmax:%.2f\n", cSpace.WMin, cSpace.WMax)
	legalPoint := func(p *Point) bool {
		inConfigSpace := (cSpace.XMin < p.X && p.X < cSpace.XMax) && (cSpace.YMin < p.Y && p.Y < cSpace.YMax)
		legalVelocities := (cSpace.VMin < p.V && p.V < cSpace.VMax) && (cSpace.WMin < p.W && p.W < cSpace.WMax)
		if !inConfigSpace || !legalVelocities {
			return false
		}

		for _, circle := range obstacles {
			if near(newVertex(p.X, p.Y, 0, 0, 0, nil), circle) {
				return false
			}
		}
		return true
	}

	return func(p *Point) bool {
		for _, robotPoint := range bot {
			globalPoint := robotPointGlobal(p, &robotPoint)
			if !legalPoint(globalPoint) {
				return false
			}
		}
		return true
	}
}

func robotPointGlobal(base, offset *Point) *Point {
	b := vec2.T{offset.X, offset.Y}
	b.Rotate(base.Theta)
	a := vec2.T{base.X, base.Y}
	c := a.Add(&b)

	return &Point{c[0], c[1], 0, 0, 0}
}

// deprecated
func getPointsAlongPath(start, end Point, epsilon float64) []Point {
	a := vec2.T{start.X, start.Y}
	b := vec2.T{end.X, end.Y}

	a2b := b.Sub(&a)
	a2bNorm := a2b.Normalized()

	distance := a2b.Length()
	numPoints := math.Floor(distance / epsilon)

	diffTheta, dir := angleDiff(start.Theta, end.Theta)
	deltaTheta := diffTheta / (numPoints - 1)

	// no longer valid...
	pointsAlong := []Point{Point{a[0], a[1], a2b.Angle(), 0, 0}}
	for i := 0.0; i < numPoints; i++ {
		offset := a2bNorm.Scaled(epsilon)
		waypoint := a.Add(&offset)
		theta := start.Theta + deltaTheta*dir*i
		pointsAlong = append(pointsAlong, Point{waypoint[0], waypoint[1], theta, 0, 0})
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
