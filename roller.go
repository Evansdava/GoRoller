package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]

	fmt.Println(args)
	// (10d6+2d8)^2+(1d20-6d4)*12/2d4

	argString := strings.ToLower(strings.Join(args, ""))
	fmt.Println(argString)
	termSlice := parse(argString)
	fmt.Println(termSlice)

	Create(termSlice)
}

func parse(argString string) []string {
	terms := make([]string, 0)
	terms = append(terms, "(")

	num := ""
	for i := 0; i < len(argString); i++ {
		str := string(argString[i])
		if _, err := strconv.Atoi(str); err == nil {
			num += str
		} else if num != "" {
			terms = append(terms, num)
			terms = append(terms, str)
			num = ""
		} else {
			terms = append(terms, str)
		}
	}
	if num != "" {
		terms = append(terms, num)
	}

	terms = append(terms, ")")

	return terms
}

func dieRoll(numDice, dieSize int) []int {
	fmt.Println("Dice:", numDice, dieSize)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	output := make([]int, numDice+1)
	var total int

	for i := 0; i < numDice; i++ {
		newRoll := r1.Intn(dieSize) + 1
		total += newRoll
		output[i] = newRoll
	}

	output[len(output)-1] = total

	return output
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
