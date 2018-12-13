package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

type V struct {
	X int
	Y int
}

var (
	Up = V{0, -1}
	Down = V{0, 1}
	Left = V{-1, 0}
	Right = V{1, 0}
)

func (v V) Add(v2 V) V {
	return V{v.X + v2.X, v.Y + v2.Y}
}

type Cart struct {
	Pos V
	Dir V
	Turn int
}

type Game struct {
	B [][]rune
	Carts []*Cart
}

func (g *Game) String() string {
	cs := map[V]rune{}
	for _, c := range g.Carts {
		cs[c.Pos] = dirToCartMarker[c.Dir]
	}

	stringRows := make([]string, len(g.B))
	for y, row := range g.B {
		b := strings.Builder{}
		for x, r := range row {
			cm := cs[V{x, y}]
			if cm == 0 {
				b.WriteRune(r)
			} else {
				b.WriteRune(cm)
			}
		}
		stringRows[y] = b.String()
	}
	return strings.Join(stringRows, "\n")
}

func (g *Game) ContainsCart(v V) bool {
	for _, c := range g.Carts {
		if c.Pos == v {
			return true
		}
	}
	return false
}


var cartMarkerToDir = map[rune]V {
	'v': Down,
	'^': Up,
	'>': Right,
	'<': Left,
}

var dirToCartMarker = map[V]rune {
	Down: 'v',
	Up: '^',
	Right: '>',
	Left: '<',
}

var cartMarkerToTrack = map[rune]rune {
	'v': '|',
	'^': '|',
	'>': '-',
	'<': '-',
}

var crossRoadTurns = map[V][]V{
	Down: {Right, Down, Left},
	Up: {Left, Up, Right},
	Right: {Up, Right, Down},
	Left: {Down, Left, Up},
}


func ParseGame(data string) *Game {
	rows := strings.Split(data, "\n")

	game := &Game{B: make([][]rune, len(rows))}
	for y, rowS := range rows {
		row := make([]rune, len(rowS))
		for x, r := range rowS {
			t := cartMarkerToTrack[r]
			if t == 0 {
				row[x] = r
			} else {
				row[x] = t
				game.Carts = append(game.Carts, &Cart{V{x, y}, cartMarkerToDir[r], 0})
			}
		}
		game.B[y] = row
	}

	return game
}

func FindFirstCollission(game *Game) V {
	for {
		// Carts act by top-to-bottom, left-to-right
		sort.Slice(game.Carts, func(i, j int) bool {
			a, b := game.Carts[i], game.Carts[j]
			return (a.Pos.Y < b.Pos.Y) || (a.Pos.Y == b.Pos.Y && a.Pos.X < b.Pos.X)
		})

		for _, c := range game.Carts {
			newPos := c.Pos.Add(c.Dir)
			log.Printf("moving cart %v to %v", c.Pos, newPos)
			if game.ContainsCart(newPos) {
				return newPos
			}
			t := game.B[newPos.Y][newPos.X]

			var newDir V
			if t == '-' || t == '|' {
				// We're on a straight track, do nothing
				newDir = c.Dir
			} else if t == '/' {
				switch c.Dir {
				case Up:
					newDir = Right
				case Down:
					newDir = Left
				case Left:
					newDir = Down
				case Right:
					newDir = Up
				default:
					panic("unknown dir for /")
				}
			}  else if t == '\\' {
				switch c.Dir {
				case Up:
					newDir = Left
				case Down:
					newDir = Right
				case Left:
					newDir = Up
				case Right:
					newDir = Down
				default:
					panic("unknown dir for \\")
				}
			} else if t == '+' {
				newDirs := crossRoadTurns[c.Dir]
				newDir = newDirs[c.Turn % 3]
				c.Turn += 1
			} else {
				panic(fmt.Sprintf("unknown tile %c", t))
			}

			c.Pos = newPos
			c.Dir = newDir
		}

		log.Printf("game:\n%v\n", game)
	}
}

func main() {
	data, err := ioutil.ReadFile("./day13/input.txt")
	if err != nil {
		log.Fatalf("error reading input.txt: %s", err)
	}

	game := ParseGame(string(data))
	log.Printf("game:\n%v\n", game)

	collisionsPos := FindFirstCollission(game)
	log.Printf("first collission: %v", collisionsPos)
}
