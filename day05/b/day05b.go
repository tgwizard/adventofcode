package main

import (
	"io/ioutil"
	"log"
	"strings"
)

type Unit struct {
	R    string
	Next *Unit
	Prev *Unit
}

func (u *Unit) Len() int {
	if u == nil {
		return 0
	}
	return 1 + u.Next.Len()
}

func (u *Unit) String() string {
	b := &strings.Builder{}
	u.stringIntoBuilder(b)
	return b.String()
}

func (u *Unit) stringIntoBuilder(b *strings.Builder) {
	if u == nil {
		return
	}
	b.WriteString(u.R)
	u.Next.stringIntoBuilder(b)
}

func (u *Unit) Filter(predicate func(*Unit) bool) *Unit {
	if u == nil {
		return nil
	}
	if !predicate(u) {
		x := u.Next.Filter(predicate)
		if x != nil {
			x.Prev = nil
		}
		return x
	}
	u.Next = u.Next.Filter(predicate)
	if u.Next != nil {
		u.Next.Prev = u
	}
	return u
}

func UniqueRunes(data string) []string {
	x := map[string]bool{}
	for _, r := range data {
		x[strings.ToLower(string(r))] = true
	}

	var result []string
	for r := range x {
		result = append(result, r)
	}
	return result
}

func Iterate(data string) *Unit {
	runes := UniqueRunes(data)

	var bestPolymer *Unit
	bestPolymerLen := len(data)
	for _, r := range runes {
		p := ParsePolymer(data)
		p = p.Filter(func(u *Unit) bool {
			return strings.ToLower(u.R) != r
		})

		p = Eliminate(p)
		log.Printf("removing %s, result %d: %s", r, p.Len(), p)

		l := p.Len()
		if l < bestPolymerLen {
			bestPolymerLen = l
			bestPolymer = p
		}
	}

	return bestPolymer
}

func ParsePolymer(data string) *Unit {
	var head, prev *Unit
	for _, r := range data {
		u := &Unit{R: string(r), Prev: prev}
		if head == nil {
			head = u
		}
		if prev != nil {
			prev.Next = u
		}
		prev = u
	}
	return head
}

func Eliminate(u *Unit) *Unit {
	for u != nil && u.Next != nil {
		a, b := u.R, u.Next.R
		if strings.ToLower(a) != strings.ToLower(b) {
			// Not same type, skip.
			u = u.Next
			continue
		}
		if a == b {
			// Same type, same polarity.
			u = u.Next
			continue
		}

		// u and u.Next should go away from the chain.
		if u.Prev != nil {
			u.Prev.Next = u.Next.Next
		}
		if u.Next.Next != nil {
			u.Next.Next.Prev = u.Prev
		}

		if u.Prev != nil {
			// Now we go back one, check again
			u = u.Prev
		} else if u.Next.Next != nil {
			// We were at the head - go forth 2.
			u = u.Next.Next
		} else {
			panic("alsdjfaslkdjf")
		}

	}

	if u != nil {
		// Wind back to the head.
		for u.Prev != nil {
			u = u.Prev
		}
	}
	return u
}

func main() {
	data, err := ioutil.ReadFile("./day05/input.txt")
	if err != nil {
		log.Fatalf("error reading input.txt: %s", err)
	}

	polymer := Iterate(string(data))

	log.Printf("len: %d", polymer.Len())
	log.Printf("eliminated best polymer: %s", polymer)
}
