package main

import (
	"fmt"
	"log"
	"strings"
)

type V struct {
	X int
	Y int
}

type VW struct {
	X int
	Y int
	W int
}

type Grid [300][300]int

func (g *Grid) String() string {
	rows := make([]string, 300)
	for y := 0; y < 300; y += 1 {
		b := strings.Builder{}
		for x := 0; x < 300; x += 1 {
			b.WriteString(fmt.Sprintf("%3d", g[x][y]))
		}
		rows[y] = b.String()
	}

	return strings.Join(rows, "\n")
}

func ComputeIterate(serialNumber int) {
	grids := make([]*Grid, 300)
	grids[0] = NewGrid(serialNumber)
	//log.Printf("initial grid:\n%s", grids[0])

	maxPower := 0
	maxV := V{}
	maxSize := 0

	for i := 1; i < 300; i += 1 {
		grids[i] = &Grid{}
		size := i + 1
		power, v := Compute(grids, i)
		log.Printf("size: %d, power: %d, v: %v", size, power, v)

		if power > maxPower {
			maxPower = power
			maxV = v
			maxSize = size
		}

		//log.Printf("grid:\n%s", grids[i])
	}

	log.Printf("maxSize: %d, maxPower: %d, maxV: %v", maxSize, maxPower, maxV)
}

func NewGrid(serialNumber int) *Grid {
	grid := &Grid{}
	// Set initial power levels.
	for x := 0; x < 300; x += 1 {
		for y := 0; y < 300; y += 1 {
			grid[x][y] = CalcPowerLevel(serialNumber, x+1, y+1)
		}
	}
	return grid
}

func Compute(grids []*Grid, i int) (int, V) {
	size := i + 1
	grid := grids[i]
	pGrid := grids[i-1]
	zGrid := grids[0]
	var qGrid *Grid
	if i-2 >= 0 {
		qGrid = grids[i-2]
	}

	maxPower := 0
	maxV := V{}

	for x := 0; x <= 300-size; x += 1 {
		for y := 0; y <= 300-size; y += 1 {
			power := 0
			if size <= 2 || size == 300 {
				power = SquarePower(pGrid, size, x, y)
			} else {
				power = pGrid[x][y] + pGrid[x+1][y+1] - qGrid[x+1][y+1] + zGrid[x][y+size-1] + zGrid[x+size-1][y]
			}

			grid[x][y] = power
			if power > maxPower {
				maxPower = power
				maxV = V{x + 1, y + 1}
			}
		}
	}

	return maxPower, maxV
}

func SquarePower(grid *Grid, squareSize, x, y int) int {
	if x+squareSize > 300 || y+squareSize > 300 {
		return 0
	}

	power := 0
	for i := 0; i < squareSize; i += 1 {
		for j := 0; j < squareSize; j += 1 {
			power += grid[x+i][y+j]
		}
	}

	return power
}

func CalcPowerLevel(serialNumber, x, y int) int {
	rackID := x + 10
	powerLevel := rackID * y
	powerLevel += serialNumber
	powerLevel *= rackID
	powerLevel = (powerLevel / 100) % 10
	powerLevel -= 5
	return powerLevel
}

func main() {
	serialNumbers := []int{
		//18,
		//42,
		//
		//// Our serial number
		5719,
	}

	for _, serialNumber := range serialNumbers {
		ComputeIterate(serialNumber)
	}
}
