package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type Registers [4]int
type OpVals [4]int

func CopyRegisters(r *Registers) *Registers {
	n := &Registers{}
	copy(n[:], r[:])
	return n
}

type Op func(*OpVals, *Registers)

func OpAddr(opVals *OpVals, regs *Registers) {
	regs[opVals[3]] = regs[opVals[1]] + regs[opVals[2]]
}

func OpAddi(opVals *OpVals, regs *Registers) {
	regs[opVals[3]] = regs[opVals[1]] + opVals[2]
}

func OpMulr(opVals *OpVals, regs *Registers) {
	regs[opVals[3]] = regs[opVals[1]] * regs[opVals[2]]
}

func OpMuli(opVals *OpVals, regs *Registers) {
	regs[opVals[3]] = regs[opVals[1]] * opVals[2]
}

func OpBanr(opVals *OpVals, regs *Registers) {
	regs[opVals[3]] = regs[opVals[1]] & regs[opVals[2]]
}

func OpBani(opVals *OpVals, regs *Registers) {
	regs[opVals[3]] = regs[opVals[1]] & opVals[2]
}

func OpBorr(opVals *OpVals, regs *Registers) {
	regs[opVals[3]] = regs[opVals[1]] | regs[opVals[2]]
}

func OpBori(opVals *OpVals, regs *Registers) {
	regs[opVals[3]] = regs[opVals[1]] | opVals[2]
}

func OpSetr(opVals *OpVals, regs *Registers) {
	regs[opVals[3]] = regs[opVals[1]]
}

func OpSeti(opVals *OpVals, regs *Registers) {
	regs[opVals[3]] = opVals[1]
}

func OpGtir(opVals *OpVals, regs *Registers) {
	val := 0
	if opVals[1] > regs[opVals[2]] {
		val = 1
	}
	regs[opVals[3]] = val
}

func OpGtri(opVals *OpVals, regs *Registers) {
	val := 0
	if regs[opVals[1]] > opVals[2] {
		val = 1
	}
	regs[opVals[3]] = val
}

func OpGtrr(opVals *OpVals, regs *Registers) {
	val := 0
	if regs[opVals[1]] > regs[opVals[2]] {
		val = 1
	}
	regs[opVals[3]] = val
}

func OpEqir(opVals *OpVals, regs *Registers) {
	val := 0
	if opVals[1] == regs[opVals[2]] {
		val = 1
	}
	regs[opVals[3]] = val
}

func OpEqri(opVals *OpVals, regs *Registers) {
	val := 0
	if regs[opVals[1]] == opVals[2] {
		val = 1
	}
	regs[opVals[3]] = val
}

func OpEqrr(opVals *OpVals, regs *Registers) {
	val := 0
	if regs[opVals[1]] == regs[opVals[2]] {
		val = 1
	}
	regs[opVals[3]] = val
}


var ops = []Op{
	OpAddr,
	OpAddi,
	OpMulr,
	OpMuli,
	OpBanr,
	OpBani,
	OpBorr,
	OpBori,
	OpSetr,
	OpSeti,
	OpGtir,
	OpGtri,
	OpGtrr,
	OpEqir,
	OpEqri,
	OpEqrr,
}

func MustAtoi(str string) int {
	ret, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return ret
}

func ParseInts(s string) [4] int{
	s = strings.Replace(s, ",", "", -1)
	p := strings.Split(s, " ")
	if len(p) != 4 {
		panic(fmt.Sprintf("not 4: %v", p))
	}

	return [4]int{
		MustAtoi(p[0]),
		MustAtoi(p[1]),
		MustAtoi(p[2]),
		MustAtoi(p[3]),
	}
}

func ComputeOpCodeMapping(data string) map[int]Op {
	rows := strings.Split(data, "\n")

	opIndexSuccess := map[int]map[int]bool{}
	for opIndex, _ := range ops {
		opIndexSuccess[opIndex] = map[int]bool{}
	}
	for i := 0; i < len(rows); {
		b, o, a := rows[i], rows[i+1], rows[i+2]
		b = strings.TrimSuffix(strings.TrimPrefix(b, "Before: ["), "]")
		a = strings.TrimSuffix(strings.TrimPrefix(a, "After:  ["), "]")

		registers := Registers(ParseInts(b))
		opVals := OpVals(ParseInts(o))
		expected := Registers(ParseInts(a))

		log.Printf("checking: %v -> %v -> %v", registers, opVals, expected)

		opCode := opVals[0]

		for opIndex, op := range ops {
			r := CopyRegisters(&registers)
			op(&opVals, r)
			if *r == expected {
				log.Printf("op %v works for %d", runtime.FuncForPC(reflect.ValueOf(op).Pointer()).Name(), opCode)
				if _, isOK := opIndexSuccess[opIndex][opCode]; !isOK {
					opIndexSuccess[opIndex][opCode] = true
				}
			} else {
				opIndexSuccess[opIndex][opCode] = false
			}
		}

		i += 4
	}

	result := map[int]Op{}
	opIndexDone := map[int]bool{}

	log.Printf("opIndexSuccess: %v", opIndexSuccess)

	for len(result) < 16 {
		for opIndex, x := range opIndexSuccess {
			if opIndexDone[opIndex] {
				continue
			}

			var okCodes []int
			for opCode, success := range x {
				if success {
					okCodes = append(okCodes, opCode)
				}
			}

			if len(okCodes) == 0 {
				panic("asdf")
			}
			if len(okCodes) == 1 {
				result[okCodes[0]] = ops[opIndex]
				opIndexDone[opIndex] = true

				for oi, x := range opIndexSuccess {
					if oi == opIndex {
						continue
					}
					x[okCodes[0]] = false
				}
			}
		}
	}

	return result
}

func ParseProgramSource(data string) []OpVals {
	rows := strings.Split(data, "\n")
	result := make([]OpVals, len(rows))

	for i, row := range rows {
		result[i] = ParseInts(row)
	}
	return result
}

func RunProgram(program []OpVals, opMap map[int]Op) *Registers {
	r := &Registers{}

	for _, opVals := range program {
		op := opMap[opVals[0]]
		op(&opVals, r)
	}

	return r
}


func main() {
	data, err := ioutil.ReadFile("./day16/inputa.txt")
	if err != nil {
		log.Fatalf("error reading inputa.txt: %s", err)
	}

	opMap := ComputeOpCodeMapping(string(data))

	log.Printf("opMap: %d\n", opMap)

	programSource, err := ioutil.ReadFile("./day16/inputb.txt")
	if err != nil {
		log.Fatalf("error reading inputb.txt")
	}

	program := ParseProgramSource(string(programSource))
	registers := RunProgram(program, opMap)
	log.Printf("final registers: %v", registers)
}