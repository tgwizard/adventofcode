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

func ParseData(data string) *Unit {
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

	polymer := ParseData(string(data))

	log.Printf("len: %d", polymer.Len())
	log.Printf("polymer: %s", polymer)

	e := Eliminate(polymer)
	log.Printf("len: %d", e.Len())
	log.Printf("eliminated polymer: %s", e)
}
