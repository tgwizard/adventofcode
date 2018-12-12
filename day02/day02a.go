// Find number of exactly 2 and exactly 3, return as product
package main

import (
	"io/ioutil"
	"log"
	"strings"
)

func ComputeChecksum(data string) (int, error) {
	ids := strings.Split(data, "\n")
	num2 := 0
	num3 := 0
	for _, id := range ids {
		counts := map[rune]int{}
		for _, r := range id {
			counts[r] += 1
		}

		done2, done3 := false, false

		for _, c := range counts {
			if c == 2 && !done2 {
				num2 += 1
				done2 = true
			} else if c == 3 && !done3 {
				num3 += 1
				done3 = true
			}
		}
		log.Printf("w: %s, counts: %v, %v, %v", id, counts, done2, done3)
	}
	return num2 * num3, nil
}

func main() {
	data, err := ioutil.ReadFile("./day02/input.txt")
	if err != nil {
		log.Fatalf("error reading input.txt: %s", err)
	}

	cs, err := ComputeChecksum(string(data))
	if err != nil {
		log.Fatalf("error computing: %s", err)
	}

	log.Printf("checksum: %d", cs)
}
