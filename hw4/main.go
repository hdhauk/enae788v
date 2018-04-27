package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"

	"github.com/pkg/profile"

	"github.com/ungerik/go3d/float64/vec2"
)

func main() {
	defer profile.Start(profile.ProfilePath(".")).Stop()

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

	// create output file for delivery
	f, err := os.OpenFile(fmt.Sprintf("problem%d_state.csv", *pIndex), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Println(err)
	}

	printPath(path, f)
	printTree(tree)

}

func printPath(path []*PathPoint, w io.WriteCloser) {
	fmt.Println("START_PATH")
	var head, tail *PathPoint
	for i := len(path) - 1; i > 0; i-- {
		head = path[i]
		tail = path[i-1]
		fmt.Printf("%.4f, %.4f, %.4f, %.4f, %.4f\n",
			head.x, head.y, tail.x, tail.y, head.θ)

		// 			    t_i,  x_i, y_i, θ_i, v_i, w_i, a_i, γ_i
	}
	fmt.Println("END_PATH")
	fmt.Println()

	// create csv file
	for i := 0; i < len(path); i++ {
		head = path[i]
		fmt.Fprintf(w, "%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f\n",
			float64(i)*timestep, head.x, head.y, head.θ, head.v, head.v, head.a, head.γ)

	}
}

func printTree(tree []*Edge) {
	fmt.Println("START_TREE")
	for _, v := range tree {
		fmt.Printf("%.4f, %.4f, %.4f, %.4f\n", v.head.X, v.head.Y, v.tail.X, v.tail.Y)
	}
	fmt.Println("END_TREE")
}

func getSafeFunc(obstacles []Circle, cSpace ConfigSpace, bot Robot) SafeFunc {
	legalPoint := func(p *PathPoint) bool {
		inConfigSpace := (cSpace.XMin < p.x && p.x < cSpace.XMax) && (cSpace.YMin < p.y && p.y < cSpace.YMax)
		legalVelocities := (cSpace.VMin < p.v && p.v < cSpace.VMax) && (cSpace.WMin < p.w && p.w < cSpace.WMax)
		if !inConfigSpace || !legalVelocities {
			return false
		}

		for _, circle := range obstacles {
			if near(newVertex(p.x, p.y, 0, 0, 0, nil), circle) {
				return false
			}
		}
		return true
	}

	return func(p *PathPoint) bool {
		for _, robotPoint := range bot {
			globalPoint := robotPointGlobal(p, &robotPoint)
			if !legalPoint(globalPoint) {
				return false
			}
		}
		return true
	}
}

func robotPointGlobal(base, offset *PathPoint) *PathPoint {
	b := vec2.T{offset.x, offset.y}
	b.Rotate(base.θ)
	a := vec2.T{base.x, base.y}
	c := a.Add(&b)

	return &PathPoint{x: c[0], y: c[1], θ: 0, v: 0, w: 0}
}

// deprecated in hw4
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
