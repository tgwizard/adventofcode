package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

const IDAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type P struct {
	X int
	Y int
}

type Point struct {
	P  P
	ID string
}

func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func Manhattan(a, b P) int {
	return Abs(b.X-a.X) + Abs(b.Y-a.Y)
}

type Grid struct {
	W int
	H int
	T [][]string
}

func (g *Grid) String() string {
	rowBuilders := make([]strings.Builder, g.H)
	for x := 0; x < g.W; x += 1 {
		for y := 0; y < g.H; y += 1 {
			rowBuilders[y].WriteString(g.T[x][y])
		}
	}

	rows := make([]string, g.H)
	for i, b := range rowBuilders {
		rows[i] = b.String()
	}

	return strings.Join(rows, "\n")
}

func MustAtoi(str string) int {
	ret, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return ret
}

func ParsePoints(data string) []Point {
	rows := strings.Split(data, "\n")
	points := make([]Point, len(rows))
	for i, row := range rows {
		parts := strings.Split(row, ", ")
		id := string(IDAlphabet[i])
		points[i] = Point{P{MustAtoi(parts[0]), MustAtoi(parts[1])}, id}
	}
	return points
}

func NewGrid(points []Point) *Grid {
	maxX, maxY := 0, 0

	for _, p := range points {
		if p.P.X > maxX {
			maxX = p.P.X
		}
		if p.P.Y > maxY {
			maxY = p.P.Y
		}
	}

	t := make([][]string, maxX+1)
	for i := range t {
		t[i] = make([]string, maxY+1)
		for j := range t[i] {
			t[i][j] = " "
		}
	}

	return &Grid{
		W: maxX + 1,
		H: maxY + 1,
		T: t,
	}
}

func MarkGrid(g *Grid, points []Point) {
	for x := 0; x < g.W; x += 1 {
		for y := 0; y < g.H; y += 1 {
			totalD := 0
			for _, p := range points {
				totalD += Manhattan(p.P, P{x, y})
			}

			if totalD < 10000 {
				g.T[x][y] = "#"
			}
		}
	}
}

func ComputeArea(g *Grid) int {
	area := 0

	for x := 0; x < g.W; x += 1 {
		for y := 0; y < g.H; y += 1 {
			if g.T[x][y] == "#" {
				area += 1
			}
		}
	}

	return area
}

func main() {
	data, err := ioutil.ReadFile("./day06/input.txt")
	if err != nil {
		log.Fatalf("error reading input.txt: %s", err)
	}

	points := ParsePoints(string(data))
	log.Printf("points: %v", points)

	grid := NewGrid(points)
	MarkGrid(grid, points)
	log.Printf("grid: \n%s", grid)

	area := ComputeArea(grid)
	log.Printf("area: %d", area)
}
