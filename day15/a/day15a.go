package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

const (
	Elf rune = 'E'
	Goblin rune = 'G'
)

func EnemyKind(kind rune) rune {
	switch kind {
	case Elf:
		return Goblin
	case Goblin:
		return Elf
	default:
		panic("unknown kind")
	}
}

type V struct {
	X int
	Y int
}

func (v V) Add(w V) V {
	return V{v.X + w.X, v.Y + w.Y}
}

var (
	Empty = V{0, 0} // This works because we always have a # in the top left corner.
	Up = V{0, -1}
	Left = V{-1, 0}
	Right= V{1, 0}
	Down = V{0, 1}
)

// AdjacentTiles return adjacent tiles, in reading order.
func AdjacentTiles(v V) []V {
	return []V{
		v.Add(Up),
		v.Add(Left),
		v.Add(Right),
		v.Add(Down),
	}
}

type Game struct {
	B [][]rune
	Monsters []*Monster
}

type Monster struct {
	Kind rune
	Pos V
	Dead bool
	HP int
	Power int
}

func (g *Game) String() string {
	mm := g.MonsterPosMap()
	b := strings.Builder{}
	for y, row := range g.B {
		for x, r := range row {
			m := mm[V{x, y}]
			if m != nil && !m.Dead {
				b.WriteRune(m.Kind)
			} else {
				b.WriteRune(r)
			}
		}

		var rowM []*Monster
		for _, m := range g.Monsters {
			if !m.Dead && m.Pos.Y == y{
				rowM = append(rowM, m)
			}
		}
		sort.Slice(rowM, func(i, j int) bool {
			return rowM[i].Pos.X <rowM[j].Pos.X
		})

		b.WriteString("   ")
		for _, m := range rowM {
			b.WriteString(fmt.Sprintf("%c(%d) ", m.Kind, m.HP))
		}
		b.WriteRune('\n')
	}
	return b.String()
}

func (g *Game) MonsterPosMap() map[V]*Monster {
	r := map[V]*Monster{}
	for _, m := range g.Monsters {
		r[m.Pos] = m
	}
	return r
}

func (g *Game) AdjacentEmptyTiles(v V) []V {
	adjTiles := AdjacentTiles(v)
	var tiles []V
	monsterPos := map[V]bool{}
	for _, m := range g.Monsters {
		if m.Dead {
			continue
		}
		monsterPos[m.Pos] = true
	}
	for _, t := range adjTiles {
		if g.B[t.Y][t.X] == '.' && !monsterPos[V{t.X, t.Y}] {
			tiles = append(tiles, t)
		}
	}
	return tiles
}

func ParseGame(data string) *Game {
	rows := strings.Split(data, "\n")
	game := &Game{B: make([][]rune, len(rows))}
	for y, rowS := range rows {
		row := make([]rune, len(rowS))

		for x, r := range rowS {
			if r == Elf || r == Goblin {
				m := &Monster{r ,V{x, y}, false, 200, 3}
				game.Monsters = append(game.Monsters, m)
				row[x] = '.'
			} else {
				row[x] = r
			}
		}

		game.B[y] = row
	}

	return game
}

func PlayGame(game *Game) (int, int) {
	for round := 0; ; round += 1 {
		monsterActed := false
		var monsters []*Monster
		aliveMonsters := map[V]*Monster{}

		for _, m := range game.Monsters {
			if m.Dead {
				continue
			}
			monsters = append(monsters, m)
			aliveMonsters[m.Pos] = m
		}

		sort.Slice(monsters, func(i, j int) bool {
			a, b := monsters[i], monsters[j]
			return (a.Pos.Y < b.Pos.Y) || (a.Pos.Y == b.Pos.Y && a.Pos.X < b.Pos.X)
		})

		combatEnded := func() bool {
			x := map[rune]int{}
			for _, m := range aliveMonsters {
				x[m.Kind] += 1
			}
			return x[Goblin] == 0 || x[Elf] == 0
		}

		totalHitPoints := func() int {
			hp := 0
			for _, m := range aliveMonsters {
				hp += m.HP
			}
			return hp
		}

		attack := func (m *Monster, enemies []*Monster) {
			sort.Slice(enemies, func(i, j int) bool {
				a, b := enemies[i], enemies[j]
				x := []V{a.Pos, b.Pos}
				SortInReadingOrder(x)

				return (a.HP < b.HP) || (a.HP == b.HP && a.Pos == x[0])
			})

			target := enemies[0]
			target.HP -= m.Power
			if target.HP <= 0 {
				target.Dead = true
				delete(aliveMonsters, target.Pos)
			}
		}


		for _, m := range monsters {
			if m.Dead {
				continue
			}


			// Attack?
			enemies := AdjacentEnemies(m, aliveMonsters)
			if len(enemies) > 0 {
				attack(m, enemies)
				if combatEnded() {
					return round, totalHitPoints()
				}

				monsterActed = true
				continue
			}

			// Move
			targets := game.FindTargets(m, aliveMonsters)
			if len(targets) == 0 {
				continue
			}

			path := game.FindShortestBestPath(m, targets)

			if len(path) > 0 {
				delete(aliveMonsters, m.Pos)
				m.Pos = path[0]
				aliveMonsters[m.Pos] = m

				monsterActed = true
			}

			// Attack?
			enemies = AdjacentEnemies(m, aliveMonsters)
			if len(enemies) > 0 {
				attack(m, enemies)
				if combatEnded() {
					return round, totalHitPoints()
				}

				monsterActed = true
			}
		}

		log.Printf("game after round %d:\n%s", round, game)

		if !monsterActed {
			panic(fmt.Sprintf("no monster acted during round %d", round))
		}
	}
	return 0, 0
}

func AdjacentEnemies(m *Monster, aliveMonsters map[V]*Monster) []*Monster {
	enemyKind := EnemyKind(m.Kind)
	var enemies []*Monster

	for _, t := range AdjacentTiles(m.Pos) {
		om := aliveMonsters[t]
		if om != nil && om.Kind == enemyKind {
			enemies = append(enemies, om)
		}
	}

	return enemies
}


func SortInReadingOrder(vs []V) {
	sort.Slice(vs, func(i, j int) bool {
		a, b := vs[i], vs[j]
		return a.Y < b.Y || (a.Y == b.Y && a.X < b.X)
	})
}

func (g *Game) FindTargets(m *Monster, aliveMonsters map[V]*Monster) []V {
	enemyKind := EnemyKind(m.Kind)

	var enemies []*Monster
	for _, om := range aliveMonsters {
		if om.Kind == enemyKind {
			enemies = append(enemies, om)
		}
	}

	targetMap := map[V]bool{}
	for _, e := range enemies {
		et := g.AdjacentEmptyTiles(e.Pos)
		for _, t := range et {
			targetMap[t] = true
		}
	}

	var targets []V
	for t := range targetMap {
		targets = append(targets, t)
	}

	return targets
}

func (g *Game) FindShortestBestPath(m *Monster, targets []V) []V {
	var bestPath []V
	for _, target := range targets {
		sp := g.FindShortestPath(m.Pos, target)
		if len(sp) == 0 {
			continue // Not reachable.
		}

		if len(bestPath) == 0 {
			bestPath = sp
			continue
		}

		spX, bpX := sp[len(sp)-1], bestPath[len(bestPath)-1]

		if len(sp) < len(bestPath) {
			bestPath = sp
		} else if len(sp) == len(bestPath) && spX != bpX {
			x := []V{spX, bpX}

			SortInReadingOrder(x)
			if x[0] == spX {
				bestPath = sp
			}
		}
	}

	return bestPath
}

func (g *Game) FindShortestPath(start V, target V) []V {
	visited := map[V]bool{start: true}
	q := g.AdjacentEmptyTiles(start)
	parents := map[V]V{}
	for _, v := range q {
		parents[v] = start
	}
	for len(q) != 0 {
		v := q[0]
		q = q[1:]
		if visited[v] {
			continue
		}
		visited[v] = true

		if v == target {
			reversePath := []V{target}
			for parent := parents[v]; parent != start; parent = parents[parent] {
				reversePath = append(reversePath, parent)
			}

			path := make([]V, len(reversePath))
			for i, vv := range reversePath {
				path[len(reversePath) - i - 1] = vv
			}
			return path
		}

		adj := g.AdjacentEmptyTiles(v)
		for _, vv := range adj {
			if parents[vv] != Empty {
				continue
			}
			parents[vv] = v
			q = append(q, vv)
		}
	}


	return nil
}

func main() {
	data, err := ioutil.ReadFile("./day15/input.txt")
	if err != nil {
		log.Fatalf("error reading input.txt: %s", err)
	}

	game := ParseGame(string(data))
	log.Printf("game:\n%s", game)

	rounds, totalHitPoints := PlayGame(game)
	log.Printf("final game:\n%s", game)
	outcome := rounds * totalHitPoints
	log.Printf("Combat ends after %d full rounds\n", rounds)
	log.Printf("XX win with %d total hit points left\n", totalHitPoints)
	log.Printf("Outcome: %d * %d = %d\n", rounds, totalHitPoints, outcome)
}