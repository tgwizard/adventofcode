package main

import (
	"fmt"
	"log"
	"strings"
)

type PlacedMarble struct {
	Score int
	Left  *PlacedMarble
	Right *PlacedMarble
}

type Board struct {
	CurrentMarble *PlacedMarble
}

func NewBoard() *Board {
	z := &PlacedMarble{Score: 0}
	z.Left = z
	z.Right = z
	b := &Board{
		CurrentMarble: z,
	}
	return b
}

func (b *Board) Act(m int) int {
	score := 0
	if m%23 == 0 {
		score += m

		r := b.CurrentMarble.Left.Left.Left.Left.Left.Left.Left
		score += r.Score

		// Remove it
		r.Left.Right = r.Right
		r.Right.Left = r.Left

		b.CurrentMarble = r.Right

	} else {
		left := b.CurrentMarble.Right
		right := left.Right
		current := &PlacedMarble{m, left, right}
		left.Right = current
		right.Left = current
		b.CurrentMarble = current
	}
	return score
}

func (b *Board) String() string {
	x := b.CurrentMarble
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%d ", x.Score))
	for n := x.Right; n != x; n = n.Right {
		builder.WriteString(fmt.Sprintf("%d ", n.Score))
	}
	return builder.String()
}

func Play(numPlayers int, numMarbles int) int {
	board := NewBoard()
	playerScores := make([]int, numPlayers)

	for m := 1; m <= numMarbles; m += 1 {
		playerI := (m - 1) % numPlayers
		playerScores[playerI] += board.Act(m)
		//log.Print(board)
	}

	max := 0
	for _, s := range playerScores {
		if s > max {
			max = s
		}
	}

	return max
}

func main() {
	games := []struct {
		Players int
		Marbles int
	}{
		{10, 1618},
		{13, 7999},
		{17, 1104},
		{21, 6111},
		{30, 5807},

		// Our input
		{473, 70904},

		// B
		{473, 70904 * 100},
	}

	for _, game := range games {
		score := Play(game.Players, game.Marbles)
		log.Printf("game (%d, %d), score: %d", game.Players, game.Marbles, score)
	}
}
