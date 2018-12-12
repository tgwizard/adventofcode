// Find the sum of the frequencies.
package main

import (
	"errors"
	"io/ioutil"
	"log"
	"strings"
)

func ComputeChecksum(data string) (string, error) {
	ids := strings.Split(data, "\n")
	n := len(ids[0])
	for i := 0; i < n; i++ {
		m := map[string]int{}
		for _, id := range ids {
			var a, b string
			if i > 0 {
				a = id[0:i]
			} else {
				a = ""
			}

			if i < n-1 {
				b = id[i+1:]
			} else {
				b = ""
			}

			m[a+b] += 1
		}
		for match, c := range m {
			if c == 2 {
				log.Printf("match: %d", i)
				return match, nil
			}
		}
	}
	return "", errors.New("no match found")
}

func main() {
	data, err := ioutil.ReadFile("./day02/input.txt")
	if err != nil {
		log.Fatalf("error reading input.txt: %s", err)
	}

	match, err := ComputeChecksum(string(data))
	if err != nil {
		log.Fatalf("error computing: %s", err)
	}

	log.Printf("matching letters: %s", match)
}
