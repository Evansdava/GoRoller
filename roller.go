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
	m[""] = add
	m["+"] = add
	m["-"] = subtract
	m["*"] = multiply
	m["/"] = divide
	m["^"] = power
	m["d"] = strRoll

	// fmt.Println(args)
	// (10d6+2d8)^2+(1d20-6d4)*12/2d4

	argString := strings.ToLower(strings.Join(args, ""))
	argString = "(10d6+2d8)^2+(1d20-6d4)*12/2d4"
	fmt.Println(argString)
	termSlice := parse(argString)
	fmt.Println(termSlice)

	postfix := CreatePostfix(termSlice)
	fmt.Println(evalPostfix(postfix, m))
	// tree := Create(termSlice)
	// fmt.Println(tree.String())
	// fmt.Println(evalTree(tree.root, m))
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

func evalPostfix(postfix *Postfix, m map[string]func(string, string) string) []string {
	outStack := []string{}
	var leftNum, rightNum string
	for len(postfix.data) > 0 {
		// fmt.Println(postfix.data)
		// fmt.Println(outStack)
		token := postfix.data[0]
		if strings.ContainsAny(token, "1234567890") {
			outStack = append(outStack, token)
		} else if strings.Contains("+-*/d^", token) {
			fn := m[token]
			if postfix.assoc[token] == "l" {
				outStack, leftNum = Pop(outStack)
				outStack, rightNum = Pop(outStack)
			} else if postfix.assoc[token] == "r" {
				outStack, rightNum = Pop(outStack)
				outStack, leftNum = Pop(outStack)
			}
			outStack = append(outStack, fn(leftNum, rightNum))
		}
		postfix.data[0] = ""
		if len(postfix.data) > 0 {
			postfix.data = postfix.data[1:]
		}
	}
	return outStack
}

func evalTree(curNode *node, m map[string]func(string, string) string) string {
	// fmt.Println("Current", curNode)
	// fmt.Println(curNode.left, curNode.right)
	if curNode.left == nil && curNode.right == nil {
		// fmt.Println("Neither")
		// fmt.Println(curNode.data)
		return curNode.data
	}
	fn := m[curNode.data]
	if curNode.left == nil {
		// fmt.Println("No left")
		// fmt.Println("Right", curNode.right)
		result := fn("", evalTree(curNode.right, m))
		// fmt.Println(result)
		return result
	} else if curNode.right == nil {
		// fmt.Println("No right")
		// fmt.Println("Left", curNode.left)
		result := fn(evalTree(curNode.left, m), "")
		// fmt.Println(result)
		return result
	} else {
		// fmt.Println("Both")
		// fmt.Println("Left", curNode.left)
		// fmt.Println("Right", curNode.right)
		result := fn(evalTree(curNode.left, m), evalTree(curNode.right, m))
		// fmt.Println(result)
		return result
	}
}

func add(leftNum string, rightNum string) string {
	// fmt.Print("add: ")
	if leftNum == "" {
		leftNum = "0"
	}
	if rightNum == "" {
		rightNum = "0"
	}
	left, _ := strconv.ParseFloat(leftNum, 64)
	right, _ := strconv.ParseFloat(rightNum, 64)
	result := string(strconv.FormatFloat(left+right, 'g', -1, 64))
	// fmt.Println(result)
	return result
}

func subtract(leftNum string, rightNum string) string {
	// fmt.Print("subtract: ")
	if leftNum == "" {
		leftNum = "0"
	}
	if rightNum == "" {
		rightNum = "0"
	}
	left, _ := strconv.ParseFloat(leftNum, 64)
	right, _ := strconv.ParseFloat(rightNum, 64)
	result := string(strconv.FormatFloat(left-right, 'g', -1, 64))
	// fmt.Println(result)
	return result
}

func multiply(leftNum string, rightNum string) string {
	// fmt.Print("multiply: ")
	if leftNum == "" {
		leftNum = "1"
	}
	if rightNum == "" {
		rightNum = "1"
	}
	left, _ := strconv.ParseFloat(leftNum, 64)
	right, _ := strconv.ParseFloat(rightNum, 64)
	result := string(strconv.FormatFloat(left*right, 'g', -1, 64))
	// fmt.Println(result)
	return result
}

func divide(leftNum string, rightNum string) string {
	// fmt.Print("divide: ")
	if leftNum == "" {
		leftNum = "1"
	}
	if rightNum == "" {
		rightNum = "1"
	}
	left, _ := strconv.ParseFloat(leftNum, 64)
	right, _ := strconv.ParseFloat(rightNum, 64)
	result := string(strconv.FormatFloat(left/right, 'g', -1, 64))
	// fmt.Println(result)
	return result
}

func power(leftNum string, rightNum string) string {
	// fmt.Print("power: ")
	if leftNum == "" {
		leftNum = "1"
	}
	if rightNum == "" {
		rightNum = "1"
	}
	left, _ := strconv.ParseFloat(leftNum, 64)
	right, _ := strconv.ParseFloat(rightNum, 64)
	result := string(strconv.FormatFloat(math.Pow(left, right), 'g', -1, 64))
	fmt.Println(left, right, result)
	return result
}

func addDice(dieString string) string {
	// fmt.Print("addDice: ")
	var total int
	for _, die := range dieString {
		dieInt, _ := strconv.Atoi(string(die))
		total += dieInt
	}
	result := strconv.Itoa(total)
	// fmt.Println(result)
	return result
}

func strRoll(leftNum string, rightNum string) string {
	// fmt.Print("strRoll: ")
	// fmt.Println(leftNum, rightNum)
	if leftNum == "" {
		leftNum = "1"
	}
	if rightNum == "" {
		rightNum = "20"
	}
	left, _ := strconv.Atoi(leftNum)
	right, _ := strconv.Atoi(rightNum)
	rolls := dieRoll(left, right)

	strRolls := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(rolls)), "+"), "[]")
	lastIndex := strings.LastIndex(strRolls, "+")
	// fmt.Println(strRolls)
	strRolls = strRolls[:lastIndex] + "=" + strRolls[lastIndex+1:]

	// result := strings.Join(strRolls[0:len(rolls)-1], "+")
	// fmt.Println(result)
	return strRolls[lastIndex+1:]
}

func dieRoll(numDice, dieSize int) []int {
	if numDice < 0 {
		numDice = -numDice
	}
	// fmt.Println("Dice:", numDice, dieSize)
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
