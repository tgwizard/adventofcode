package main

import (
	"io/ioutil"
	"log"
	"regexp"
	"sort"
	"strings"
)

var lineRegex = regexp.MustCompile("^Step ([A-Z]) must be finished before step ([A-Z]) can begin.$")

type Node struct {
	ID       string
	Forward  map[string]*Node
	Backward map[string]*Node
}

func NewNode(id string) *Node {
	return &Node{ID: id, Forward: map[string]*Node{}, Backward: map[string]*Node{}}
}

type Graph struct {
	Nodes map[string]*Node
}

func ParseGraph(data string) *Graph {
	rows := strings.Split(data, "\n")

	graph := &Graph{Nodes: map[string]*Node{}}
	for _, row := range rows {
		groups := lineRegex.FindStringSubmatch(row)
		if len(groups) != 3 {
			panic("line regex mismatch")
		}

		a := groups[1]
		b := groups[2]

		nodeA := graph.Nodes[a]
		if nodeA == nil {
			nodeA = NewNode(a)
			graph.Nodes[a] = nodeA
		}

		nodeB := graph.Nodes[b]
		if nodeB == nil {
			nodeB = NewNode(b)
			graph.Nodes[b] = nodeB
		}

		nodeA.Forward[b] = nodeB
		nodeB.Backward[a] = nodeA
		delete(nodeA.Backward, a)
	}

	return graph
}

func TopSort(graph *Graph) []*Node {
	order := make([]*Node, 0, len(graph.Nodes))
	visited := map[string]bool{}
	var potentials []*Node
	for _, n := range graph.Nodes {
		if len(n.Backward) == 0 {
			potentials = append(potentials, n)
		}
	}

	for len(potentials) > 0 {
		sort.Slice(potentials, func(i, j int) bool {
			return potentials[i].ID > potentials[j].ID
		})
		x := potentials[len(potentials)-1]
		potentials = potentials[0 : len(potentials)-1]

		order = append(order, x)
		visited[x.ID] = true

		for _, p := range x.Forward {
			c := 0
			for _, pp := range p.Backward {
				if !visited[pp.ID] {
					c += 1
					break
				}
			}
			if c == 0 {
				potentials = append(potentials, p)
			}
		}
	}

	if len(order) != len(graph.Nodes) {
		panic("not all nodes visited")
	}

	return order
}

func ToIDs(ns []*Node) []string {
	ids := make([]string, len(ns))
	for i, n := range ns {
		ids[i] = n.ID
	}
	return ids
}

func main() {
	data, err := ioutil.ReadFile("./day07/input.txt")
	if err != nil {
		log.Fatalf("error reading input.txt: %s", err)
	}

	graph := ParseGraph(string(data))

	ns := TopSort(graph)
	ids := ToIDs(ns)
	log.Printf("result: %s", strings.Join(ids, ""))
}
