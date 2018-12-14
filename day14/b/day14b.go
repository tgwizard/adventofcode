package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Score struct {
	Val int
	Next *Score
	Prev *Score
}

type Game struct {
	NumScores int
	LastScore *Score
	PlayerScores []*Score
}

func NewGame() *Game {
	a := &Score{Val: 3}
	b := &Score{Val:7 }
	a.Next = b
	a.Prev = b
	b.Next = a
	b.Prev = a
	return &Game{2, b, []*Score{a, b}}
}

func (g *Game) String() string {
	b := strings.Builder{}
	start := g.LastScore.Next
	b.WriteString(fmt.Sprintf("%d ", start.Val))
	for s := start.Next; s != start; s = s.Next {
		b.WriteString(fmt.Sprintf("%d ", s.Val))
	}
	return b.String()
}

func (g *Game) AddScore(s *Score) {
	s.Next = g.LastScore.Next
	g.LastScore.Next.Prev = s
	g.LastScore.Next = s
	s.Prev = g.LastScore
	g.LastScore = s

	g.NumScores += 1
}

func (g *Game) MovePlayer(playerIndex int, steps int) {
	for ; steps > 0; steps-- {
		g.PlayerScores[playerIndex] = g.PlayerScores[playerIndex].Next
	}
}

func Compute(targetSequence []int) int {
	game := NewGame()
	matchingScores := 0
	if game.LastScore.Next.Val == targetSequence[0] {
		matchingScores += 1
		if game.LastScore.Next.Next.Val == targetSequence[1] {
			matchingScores += 1
		}
	} else if game.LastScore.Next.Next.Val == targetSequence[0] {
		matchingScores += 1
	}

	addScore := func (val int) {
		if matchingScores < len(targetSequence) {
			if targetSequence[matchingScores] == val {
				matchingScores += 1
			} else {
				matchingScores = 0
				if targetSequence[matchingScores] == val {
					matchingScores += 1
				}
			}
		} else {
			// Just to make rewind easier
			matchingScores += 1
		}

		game.AddScore(&Score{Val: val})
	}

	for matchingScores < len(targetSequence) {
		//log.Printf("gg(%d): %s", matchingScores, game)
		sum := game.PlayerScores[0].Val + game.PlayerScores[1].Val

		if sum >= 100 {
			panic("unexpectedly large sum")
		}

		if sum >= 10 {
			addScore(sum / 10)
		}
		addScore(sum % 10)

		game.MovePlayer(0, game.PlayerScores[0].Val + 1)
		game.MovePlayer(1, game.PlayerScores[1].Val + 1)
	}

	return game.NumScores - matchingScores

}

func main() {
	d := [][]int{
		//{5, 1,5, 8, 9},
		//{0, 1, 2, 4, 5},
		//{9, 2, 5, 1, 0},
		//{5, 9, 4, 1, 4},

		// Our target:
		{2, 0, 9, 2, 3, 1},
	}

	for _, i := range d {
		st := time.Now()
		result := Compute(i)
		log.Printf("result for %d: %d (time taken: %s)", i, result, time.Since(st))
	}
}
