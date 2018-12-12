// Find how many squares are overlapped by 2 or more rectangles
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var rectRegex = regexp.MustCompile("^#([0-9]+) @ ([0-9]+),([0-9]+): ([0-9]+)x([0-9]+)$")

type Rect struct {
	ID int
	X  int
	Y  int
	W  int
	H  int
}

func (r Rect) Right() int {
	return r.X + r.W
}

func (r Rect) Bottom() int {
	return r.Y + r.H
}

func MustAtoi(str string) int {
	ret, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return ret
}

func ParseRects(data string) ([]Rect, error) {
	parts := strings.Split(data, "\n")
	rects := make([]Rect, len(parts))
	for i, part := range parts {
		groups := rectRegex.FindStringSubmatch(part)
		if len(groups) != 6 {
			return nil, fmt.Errorf("invalid number of groups: %s, %v", part, groups)
		}

		rects[i] = Rect{
			MustAtoi(groups[1]),
			MustAtoi(groups[2]),
			MustAtoi(groups[3]),
			MustAtoi(groups[4]),
			MustAtoi(groups[5]),
		}
	}
	return rects, nil
}

func ComputeOverlap(rects []Rect) (int, error) {
	maxX, maxY := 0, 0
	for _, rect := range rects {
		if rect.Right() > maxX {
			maxX = rect.Right()
		}
		if rect.Bottom() > maxY {
			maxY = rect.Bottom()
		}
	}

	fabric := make([][]int, maxX)
	for i := range fabric {
		fabric[i] = make([]int, maxY)
	}

	for _, rect := range rects {
		for x := rect.X; x < rect.Right(); x += 1 {
			for y := rect.Y; y < rect.Bottom(); y += 1 {
				fabric[x][y] += 1
			}
		}
	}

	overlaps := 0
	for x := 0; x < maxX; x += 1 {
		for y := 0; y < maxY; y += 1 {
			if fabric[x][y] > 1 {
				overlaps += 1
			}
		}
	}

	return overlaps, nil
}

func main() {
	data, err := ioutil.ReadFile("./day03/input.txt")
	if err != nil {
		log.Fatalf("error reading input.txt: %s", err)
	}

	rects, err := ParseRects(string(data))
	if err != nil {
		log.Fatalf("error parsing rects: %s", err)
	}

	overlap, err := ComputeOverlap(rects)
	if err != nil {
		log.Fatalf("error computing overlap: %s", err)
	}

	log.Printf("overlap: %d", overlap)
}
