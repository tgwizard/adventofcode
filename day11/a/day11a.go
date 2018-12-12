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

type Grid [300][300]int

func (g Grid) String() string {
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

func Compute(serialNumber int) V {
	grid := Grid{}

	// Set initial power levels.
	for x := 0; x < 300; x += 1 {
		for y := 0; y < 300; y += 1 {
			grid[x][y] = CalcPowerLevel(serialNumber, x+1, y+1)
		}
	}

	// log.Printf("grid:\n%s", grid)
	//log.Printf("xxxx %d", grid[2][4])

	maxPower := 0
	maxV := V{}

	for x := 0; x < 297; x += 1 {
		for y := 0; y < 297; y += 1 {
			p := SquarePower(&grid, x, y)
			if p > maxPower {
				maxPower = p
				maxV = V{x + 1, y + 1}
			}
		}
	}

	log.Printf("POWER: %d, %v", maxPower, maxV)

	return maxV
}

func SquarePower(grid *Grid, x, y int) int {
	if x+3 >= 300 || y+3 >= 300 {
		return 0
	}

	power := 0
	for i := 0; i < 3; i += 1 {
		for j := 0; j < 3; j += 1 {
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
		18,
		42,

		// Our serial number
		5719,
	}

	for _, serialNumber := range serialNumbers {
		result := Compute(serialNumber)
		log.Printf("result for %d: %v", serialNumber, result)
	}
}
