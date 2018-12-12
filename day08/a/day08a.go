package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func MustAtoi(str string) int {
	ret, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return ret
}

type Node struct {
	Meta     []int
	Children []*Node
}

func (n *Node) String() string {
	if n == nil {
		return "nil"
	}

	cs := make([]string, len(n.Children))
	for i, c := range n.Children {
		cs[i] = c.String()
	}

	return fmt.Sprintf("node(%v): [%s]", n.Meta, strings.Join(cs, ", "))
}

func (n *Node) Visit(visitor func(*Node)) {
	if n == nil {
		return
	}

	visitor(n)
	for _, c := range n.Children {
		c.Visit(visitor)
	}
}

func (n *Node) MetaSum() int {
	sum := 0
	n.Visit(func(node *Node) {
		for _, m := range node.Meta {
			sum += m
		}
	})
	return sum
}

func ParseTree(data string) *Node {
	parts := strings.Split(data, " ")
	elements := make([]int, len(parts))
	for i, part := range parts {
		elements[i] = MustAtoi(part)
	}

	node, elements := ParseNode(elements)
	if len(elements) != 0 {
		panic("elements remaining after parsing full tree")
	}

	return node
}

func ParseNode(elements []int) (*Node, []int) {
	if len(elements) < 2 {
		panic("too few elements for node")
	}
	numChildren, numMeta := elements[0], elements[1]
	elements = elements[2:]
	children := make([]*Node, numChildren)
	for i := 0; i < numChildren; i += 1 {
		child, newElements := ParseNode(elements)
		elements = newElements
		children[i] = child
	}

	meta := elements[:numMeta]
	elements = elements[numMeta:]

	return &Node{meta, children}, elements
}

func main() {
	data, err := ioutil.ReadFile("./day08/input.txt")
	if err != nil {
		log.Fatalf("error reading input.txt: %s", err)
	}

	tree := ParseTree(string(data))
	log.Printf("tree: %s", tree)

	metaSum := tree.MetaSum()
	log.Printf("meta sum: %d", metaSum)
}
