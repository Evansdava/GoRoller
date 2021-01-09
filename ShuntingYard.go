package main

import (
	"fmt"
	"strings"
)

// Postfix struct, with precedence and association maps, and data queue
type Postfix struct {
	prec  map[string]int
	assoc map[string]string
	data  []string
}

// CreatePostfix - Make a new Postfix from a string of data
func CreatePostfix(data []string) *Postfix {
	// Set precedence for each operator
	prec := make(map[string]int)
	prec["+"] = 1
	prec["-"] = 1
	prec["*"] = 2
	prec["/"] = 2
	prec["^"] = 3
	prec["d"] = 4

	// Set association for each operator
	assoc := make(map[string]string)
	assoc["+"] = "l"
	assoc["-"] = "l"
	assoc["*"] = "l"
	assoc["/"] = "l"
	assoc["^"] = "r"
	assoc["d"] = "l"

	// Create stack for operators
	opStack := []string{}
	// Queue for output
	outQueue := []string{}

	var op string
	fmt.Println(op)

	for _, char := range data {
		fmt.Println(char)
		fmt.Println(opStack)
		if strings.ContainsAny(char, "1234567890") {
			outQueue = append(outQueue, char)
		} else if strings.Contains("+-*/d^", char) {
			if len(opStack) > 0 {
				for strings.Contains("+-*/d^", opStack[len(opStack)-1]) &&
					(prec[opStack[len(opStack)-1]] > prec[char] ||
						(prec[opStack[len(opStack)-1]] == prec[char] &&
							assoc[char] == "l")) &&
					opStack[len(opStack)-1] != "(" {
					// fmt.Println(outQueue)
					opStack, op = Pop(opStack)
					outQueue = append(outQueue, op)
				}
			}
			opStack = append(opStack, char)
		} else if char == "(" {
			opStack = append(opStack, char)
		} else if char == ")" {
			for opStack[len(opStack)-1] != "(" {
				// fmt.Println(opStack, opStack[len(opStack)-1])
				opStack, op = Pop(opStack)
				outQueue = append(outQueue, op)
				// fmt.Println(opStack, opStack[len(opStack)-1])
				// fmt.Println(outQueue)
			}
			if opStack[len(opStack)-1] == "(" {
				opStack, op = Pop(opStack)
			}
		}
	}
	for len(opStack) > 0 {
		opStack, op = Pop(opStack)
		outQueue = append(outQueue, op)
	}

	return &Postfix{prec: prec, assoc: assoc, data: outQueue}
}

// (10d6+2d8)^2+(1d20-6d4)*12/2d4
// Output queue: 10 6 d 2 8 d + 2 ^ 1 20 d 6 4 d - 12 * 2 4 d / +
// Operator stack:

// Pop - Remove the last element of a string slice, return both that slice and the element
func Pop(stack []string) ([]string, string) {
	// fmt.Println("Stack", stack)
	output := string(stack[len(stack)-1])
	// fmt.Println(output)
	stack[len(stack)-1] = ""
	stack = stack[:len(stack)-1]
	// fmt.Println("Stack", stack)

	return stack, output
}
