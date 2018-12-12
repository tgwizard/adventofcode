package day10

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var lineRe = regexp.MustCompile(`^position=< *(-?[0-9]+), *(-?[0-9]+)> velocity=< *(-?[0-9]+), *(-?[0-9]+)>$`)

func MustAtoi(str string) int {
	ret, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return ret
}

func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type V struct {
	X int
	Y int
}

func (v *V) IsAdjacent(v2 *V) bool {
	xdiff := Abs(v.X - v2.X)
	ydiff := Abs(v.Y - v2.Y)

	return (xdiff == 0 && ydiff == 1) || (xdiff == 1 && ydiff == 0)
}

func (v *V) Distance(v2 *V) int {
	xdiff := Abs(v.X - v2.X)
	ydiff := Abs(v.Y - v2.Y)

	return int(math.Sqrt(float64(xdiff*xdiff + ydiff*ydiff)))
}

type Grid struct {
	Points     []*V
	Velocities []*V
}

func (g *Grid) String() string {
	var min V = *g.Points[0]
	var max V

	for _, p := range g.Points {
		if p.X < min.X {
			min.X = p.X
		}
		if p.X > max.X {
			max.X = p.X
		}
		if p.Y < min.Y {
			min.Y = p.Y
		}
		if p.Y > max.Y {
			max.Y = p.Y
		}
	}

	min.X -= 1
	min.Y -= 1
	max.X += 2
	max.Y += 2

	numRows := max.Y - min.Y
	numCols := max.X - min.X

	log.Printf("%d, %d, %v, %v", numRows, numCols, min, max)

	rows := make([]string, numRows)
	for y := 0; y < numRows; y += 1 {
		b := strings.Builder{}
		for x := 0; x < numCols; x += 1 {
			hasPoint := false
			for _, p := range g.Points {
				if p.X == (min.X+x) && p.Y == (min.Y+y) {
					hasPoint = true
					break
				}
			}

			if hasPoint {
				b.WriteString("#")
			} else {
				b.WriteString(".")
			}
		}
		rows[y] = b.String()
	}

	return strings.Join(rows, "\n")
}

func ParseGrid(data string) *Grid {
	rows := strings.Split(data, "\n")
	grid := &Grid{
		make([]*V, len(rows)),
		make([]*V, len(rows)),
	}

	for i, row := range rows {
		groups := lineRe.FindStringSubmatch(row)
		if len(groups) != 5 {
			panic(fmt.Sprintf("row mismatch: %v", groups))
		}
		grid.Points[i] = &V{MustAtoi(groups[1]), MustAtoi(groups[2])}
		grid.Velocities[i] = &V{MustAtoi(groups[3]), MustAtoi(groups[4])}
	}

	return grid
}

func PrintIterations(grid *Grid, start int, step int, numSeconds int) {
	seconds := start
	prevDist := 20000000000

	StepGridIteration(grid, start)

	for i := 0; true; i += 1 {
		d := MaxDistance(grid)
		a := NumWithNeighbours(grid)
		log.Printf(">> iteration=%d, seconds=%d, maxDist=%d, numAdj=%d", i, seconds, d, a)

		if d < 100 {
			log.Printf("grid:\n%s", grid)
			step = 1
		}
		StepGridIteration(grid, step)
		seconds += step

		if d > prevDist {
			return
		}
		prevDist = d
	}
}

func StepGridIteration(grid *Grid, step int) {
	for i, p := range grid.Points {
		v := grid.Velocities[i]
		p.X = p.X + v.X*step
		p.Y = p.Y + v.Y*step
	}
}

func MaxDistance(grid *Grid) int {
	max := 0
	for i, p1 := range grid.Points {
		for _, p2 := range grid.Points[i+1:] {
			d := p1.Distance(p2)
			if d > max {
				max = d
			}
		}
	}
	return max
}

func NumWithNeighbours(grid *Grid) int {
	num := 0
	for i, p1 := range grid.Points {
		for _, p2 := range grid.Points[i+1:] {
			if p1.IsAdjacent(p2) {
				num += 1
			}
		}
	}
	return num
}

func main() {
	data, err := ioutil.ReadFile("./day10/input.txt")
	if err != nil {
		log.Fatalf("error reading input.txt: %s", err)
	}

	grid := ParseGrid(string(data))
	log.Printf("parsed grid")
	PrintIterations(grid, 10000, 1, 50000)
}
