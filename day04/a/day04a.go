// find the guard with most total sleep time, then take guard ID x most slept minute
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var lineRegex = regexp.MustCompile("^\\[1518-[0-9]{2}-[0-9]{2} [0-9]{2}:([0-9]{2})] (.+)$")
var guardRegex = regexp.MustCompile("^Guard #([0-9]+) begins shift$")

func MustAtoi(str string) int {
	ret, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return ret
}

func GetSortedData(data string) []string {
	parts := strings.Split(data, "\n")
	sort.Strings(parts)
	return parts
}

type Period struct {
	Start int
	End   int
}

func (p Period) Length() int {
	return p.End - p.Start
}

type Guard struct {
	ID    int
	Times []Period
}

func (g Guard) TotalSleepTime() int {
	sum := 0
	for _, p := range g.Times {
		sum += p.Length()
	}
	return sum
}

func ParseData(data string) (int, error) {
	events := GetSortedData(data)

	guards := map[int]*Guard{}

	var currentGuard *Guard
	var currentPeriod *Period
	for i, s := range events {
		lineGroups := lineRegex.FindStringSubmatch(s)
		if len(lineGroups) != 3 {
			return 0, fmt.Errorf("invalid line: %s, %d", s, len(lineGroups))
		}
		t := MustAtoi(lineGroups[1])
		a := lineGroups[2]

		gaGroups := guardRegex.FindStringSubmatch(a)
		if len(gaGroups) == 2 {
			if currentPeriod != nil {
				currentPeriod.End = 60
				currentGuard.Times = append(currentGuard.Times, *currentPeriod)
				currentPeriod = nil
			}

			// New guard
			guardID := MustAtoi(gaGroups[1])
			currentGuard = guards[guardID]

			if currentGuard == nil {
				// Guard we haven't seen before
				currentGuard = &Guard{
					ID:    guardID,
					Times: []Period{},
				}
				guards[guardID] = currentGuard
			}
		} else {
			if a == "falls asleep" {
				if currentPeriod != nil {
					return 0, fmt.Errorf("non-nil current period when falling asleep: %d, %s", i, s)
				}
				currentPeriod = &Period{Start: t}
			} else if a == "wakes up" {
				if currentPeriod == nil {
					return 0, fmt.Errorf("nil current period when waking up: %d, %s", i, s)
				}
				currentPeriod.End = t - 1
				currentGuard.Times = append(currentGuard.Times, *currentPeriod)
				currentPeriod = nil
			}
		}
	}

	if currentPeriod != nil {
		currentPeriod.End = 60
		currentGuard.Times = append(currentGuard.Times, *currentPeriod)
		currentPeriod = nil
	}

	// Find guard that slept the most
	var maxGuard *Guard
	maxGuardSleep := 0
	for _, guard := range guards {
		//log.Printf("guard: %v\n", guard)

		sleepTime := guard.TotalSleepTime()
		if sleepTime > maxGuardSleep {
			maxGuard = guard
			maxGuardSleep = sleepTime
		}
	}

	// Find minute they slept the most
	sort.Slice(maxGuard.Times, func(i, j int) bool {
		if maxGuard.Times[i].Start < maxGuard.Times[j].Start {
			return true
		}
		if maxGuard.Times[i].Start == maxGuard.Times[j].Start {
			return maxGuard.Times[i].End <= maxGuard.Times[j].End
		}
		return false
	})

	log.Printf("guard: %v\n", maxGuard)

	maxMin, maxCount := 0, 0
	for minute := 0; minute < 60; minute += 1 {
		counts := 0
		for _, p := range maxGuard.Times {
			if p.Start <= minute && p.End >= minute {
				counts += 1
			}
		}
		if counts > maxCount {
			maxMin = minute
			maxCount = counts
		}
	}

	log.Printf("m: %d, c: %d", maxMin, maxCount)

	return maxMin * maxGuard.ID, nil
}

func main() {
	data, err := ioutil.ReadFile("./day04/input.txt")
	if err != nil {
		log.Fatalf("error reading input.txt: %s", err)
	}

	result, err := ParseData(string(data))
	if err != nil {
		log.Fatalf("error parsing data: %s", err)
	}

	log.Printf("guard ID x minute: %d", result)
}
