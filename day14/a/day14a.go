package main

import (
	"fmt"
	"log"
	"strings"
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

func Compute(cutoff int) string {
	game := NewGame()
	target := cutoff + 10
	for game.NumScores < target {
		//log.Printf("gg: %s", game)
		sum := game.PlayerScores[0].Val + game.PlayerScores[1].Val

		if sum >= 100 {
			panic("unexpectedly large sum")
		}

		if sum >= 10 {
			game.AddScore(&Score{Val: sum / 10})
		}
		game.AddScore(&Score{Val: sum % 10})

		game.MovePlayer(0, game.PlayerScores[0].Val + 1)
		game.MovePlayer(1, game.PlayerScores[1].Val + 1)
	}

	var end *Score
	if game.NumScores == target {
		end = game.LastScore
	} else if game.NumScores == target + 1 {
		end = game.LastScore.Prev
	} else {
		panic("unexpected end")
	}

	result := make([]string, 10)
	for i := 9; i >= 0; i-- {
		result[i] = fmt.Sprintf("%d", end.Val)
		end = end.Prev
	}

	log.Printf("gg: %s", game)
	return strings.Join(result, "")
}

func main() {
	d := []int{
		//9,
		//5,
		//18,
		//2018,
		//
		//// Our input:
		209231,
	}

	for _, i := range d {
		result := Compute(i)
		log.Printf("result for %d: %s", i, result)
	}
}
