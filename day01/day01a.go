// Find the sum of the frequencies.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func ComputeFrequency(data string) (int, error) {
	parts := strings.Split(data, "\n")
	sum := 0
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
	}
	return sum, nil
}

func main() {
	data, err := ioutil.ReadFile("./day01/input.txt")
	if err != nil {
		log.Fatalf("error reading input.txt: %s", err)
	}

	freq, err := ComputeFrequency(string(data))
	if err != nil {
		log.Fatalf("error computing frequency: %s", err)
	}

	log.Printf("frequency: %d", freq)
}
