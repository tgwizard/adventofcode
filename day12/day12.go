package day12

import (
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Game struct {
	IndexZero int
	State     string
	Rules     map[string]string
}

func ParseGame(data string) *Game {
	rows := strings.Split(data, "\n")
	isRow := rows[0]
	ruleRows := rows[2:]

	initialState := strings.TrimPrefix(isRow, "initial state: ")

	rules := map[string]string{}
	for _, row := range ruleRows {
		key := row[:5]
		result := row[9:]
		if _, isOK := rules[key]; isOK {
			panic("rule already exists")
		}
		rules[key] = result
	}

	return &Game{len(initialState), EmptyS(len(initialState)) + initialState, rules}
}

func EmptyS(n int) string {
	b := strings.Builder{}
	for i := 0; i < n; i += 1 {
		b.WriteString(".")
	}
	return b.String()
}

func PlayGame(game *Game, iterations int) {
	log.Printf("initial state:    %v\n", game.State)

	start := time.Now()
	tx := start
	prevScore := 0

	compareDiff := 10000

	for i := 0; i < iterations; i += 1 {
		if game.State[0:2] != ".." {
			panic("assumption: first two states are empty")
		}

		if i%compareDiff == 0 {
			elapsedTotal := time.Since(start)
			elapsedX := time.Since(tx)
			tx = time.Now()
			score := Score(game)
			remainingIterations := iterations - i
			expectedScore := score + remainingIterations*(score-prevScore)/compareDiff
			log.Printf("iteration: %d (len: %d, score: %d, speed: %s, elapsed: %s), expeted score: %d", i, len(game.State), score, elapsedX, elapsedTotal, expectedScore)
			log.Printf("state after %4d: %v\n", i, game.State)

			prevScore = score
		}

		b := strings.Builder{}
		b.Grow(len(game.State) + 2)
		b.WriteString("..")
		hasWrittenLife := false
		indexZero := game.IndexZero + 2
		for x := 0; x < len(game.State)+2; x += 1 {
			s := BuildS(game, x)
			r := game.Rules[s]
			if r == "" {
				r = "."
			}

			if x >= len(game.State) && r == "." {
				continue
			}

			if r == "." && !hasWrittenLife {
				indexZero -= 1
				continue
			}

			if r == "#" && !hasWrittenLife {
				hasWrittenLife = true
			}

			b.WriteString(r)
		}

		game.State = b.String()
		game.IndexZero = indexZero

		//log.Printf("state after %4d: %v\n", i, game.State)
	}

	finalScore := Score(game)
	log.Printf("state after %4d (zero: %d): %v\n", iterations, game.IndexZero, game.State)
	log.Printf("final score after %d iterations: %d", iterations, finalScore)
}

func Score(game *Game) int {
	sum := 0
	for x := 0; x < len(game.State); x += 1 {
		if game.State[x] == '#' {
			index := x - game.IndexZero
			sum += index
		}
	}
	return sum
}

func BuildS(game *Game, x int) string {
	get := func(i int) string {
		if i >= 0 && i < len(game.State) {
			return string(game.State[i])
		} else {
			return "."
		}
	}

	return get(x-2) + get(x-1) + get(x) + get(x+1) + get(x+2)
}

func main() {
	data, err := ioutil.ReadFile("./day12/input.txt")
	if err != nil {
		log.Fatalf("error reading input.txt: %s", err)
	}

	game := ParseGame(string(data))
	log.Printf("game: %v\n", game)

	PlayGame(game, 50000000000)
}
