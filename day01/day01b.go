// Find the first computed frequency arrived at twice.
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func ComputeFrequencyReachedTwice(data string) (int, error) {
	parts := strings.Split(data, "\n")
	sum := 0
	reaches := map[int]bool{
		0: true,
	}
	for i := 0; i < 1000; i++ {
		for _, part := range parts {
			operator := part[0]
			operand, err := strconv.Atoi(part[1:])
			if err != nil {
				return 0, err
			}
			switch operator {
			case '+':
				sum += operand
			case '-':
				sum -= operand
			default:
				return 0, fmt.Errorf("invalid operator: %c", operator)
			}

			wasReached := reaches[sum]
			log.Printf("%d: %d, %d\n", i, sum, wasReached)
			if wasReached {
				log.Println("yeas!")
				return sum, nil
			}
			reaches[sum] = true
		}
	}

	return 0, errors.New("no frequency reached twice")
}

func main() {
	data, err := ioutil.ReadFile("./day01/input.txt")
	if err != nil {
		log.Fatalf("error reading input.txt: %s", err)
	}

	freq, err := ComputeFrequencyReachedTwice(string(data))
	if err != nil {
		log.Fatalf("error computing frequency reaches: %s", err)
	}

	log.Printf("frequency reached first twice: %d", freq)
}
