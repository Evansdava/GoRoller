package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]

	m := make(map[string]func(string, string) string)
	m["+"] = add
	m["-"] = subtract
	m["*"] = multiply
	m["/"] = divide
	m["^"] = power
	m["d"] = strRoll

	fmt.Println(args)
	// (10d6+2d8)^2+(1d20-6d4)*12/2d4

	argString := strings.ToLower(strings.Join(args, ""))
	fmt.Println(argString)
	termSlice := parse(argString)
	fmt.Println(termSlice)

	tree := Create(termSlice)
	fmt.Println(eval(tree.root, m))
	fmt.Println(tree.String())
}

func parse(argString string) []string {
	terms := make([]string, 0)
	// terms = append(terms, "(")

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

	// terms = append(terms, ")")

	return terms
}

func eval(curNode *node, m map[string]func(string, string) string) string {
	fmt.Println("Current", curNode)
	if curNode.left == nil && curNode.right == nil {
		fmt.Println("Neither")
		return curNode.data
	} else {
		fn := m[curNode.data]
		if curNode.left == nil {
			fmt.Println("No left")
			fmt.Println("Right", curNode.right)
			return fn("", eval(curNode.right, m))
		} else if curNode.right == nil {
			fmt.Println("No right")
			fmt.Println("Left", curNode.left)
			return fn(eval(curNode.left, m), "")
		} else {
			fmt.Println("Both")
			fmt.Println("Left", curNode.left)
			fmt.Println("Right", curNode.right)
			return fn(eval(curNode.left, m), eval(curNode.right, m))
		}
	}
}

func add(leftNum string, rightNum string) string {
	if leftNum == "" {
		leftNum = "0"
	}
	if rightNum == "" {
		rightNum = "0"
	}
	left, _ := strconv.ParseFloat(leftNum, 64)
	right, _ := strconv.ParseFloat(rightNum, 64)
	return string(strconv.FormatFloat(left+right, 'g', -1, 64))
}

func subtract(leftNum string, rightNum string) string {
	if leftNum == "" {
		leftNum = "0"
	}
	if rightNum == "" {
		rightNum = "0"
	}
	left, _ := strconv.ParseFloat(leftNum, 64)
	right, _ := strconv.ParseFloat(rightNum, 64)
	return string(strconv.FormatFloat(left-right, 'g', -1, 64))
}

func multiply(leftNum string, rightNum string) string {
	left, _ := strconv.ParseFloat(leftNum, 64)
	right, _ := strconv.ParseFloat(rightNum, 64)
	return string(strconv.FormatFloat(left*right, 'g', -1, 64))
}

func divide(leftNum string, rightNum string) string {
	left, _ := strconv.ParseFloat(leftNum, 64)
	right, _ := strconv.ParseFloat(rightNum, 64)
	return string(strconv.FormatFloat(left/right, 'g', -1, 64))
}

func power(leftNum string, rightNum string) string {
	left, _ := strconv.ParseFloat(leftNum, 64)
	right, _ := strconv.ParseFloat(rightNum, 64)
	return string(strconv.FormatFloat(math.Pow(left, right), 'g', -1, 64))
}

func addDice(dieString string) string {
	var total int
	for _, die := range dieString {
		dieInt, _ := strconv.Atoi(string(die))
		total += dieInt
	}
	return strconv.Itoa(total)
}

func strRoll(leftNum string, rightNum string) string {
	if leftNum == "" {
		leftNum = "1"
	}
	if rightNum == "" {
		rightNum = "20"
	}
	left, _ := strconv.Atoi(leftNum)
	right, _ := strconv.Atoi(rightNum)
	rolls := dieRoll(left, right)

	strRolls := make([]string, len(rolls))

	return strings.Join(strRolls[0:len(strRolls)-1], "+")
}

func dieRoll(numDice, dieSize int) []int {
	if numDice < 0 {
		numDice = -numDice
	}
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
