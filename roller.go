package main

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	// args := os.Args[1:]

	// m := make(map[string]func(string, string) string)
	// m[""] = add
	// m["+"] = add
	// m["-"] = subtract
	// m["*"] = multiply
	// m["/"] = divide
	// m["^"] = power
	// m["d"] = strRoll

	// fmt.Println(args)
	// (10d6+2d8)^2+(1d20-6d4)*12/2d4

	// argString := strings.ToLower(strings.Join(args, ""))
	// fmt.Println(GetRoll(argString))
	// // argString = "(10d6+2d8)^2+(1d20-6d4)*12/2d4"
	// fmt.Println(argString)
	// termSlice := parse(argString)
	// fmt.Println(termSlice)

	// postfix := CreatePostfix(termSlice)
	// fmt.Println(evalPostfix(postfix, m))

	StartBot()
}

// GetRoll parses the passed string and returns the result of the roll
func GetRoll(argString string) string {
	m := make(map[string]func(string, string, chan string) string)
	m[""] = add
	m["+"] = add
	m["-"] = subtract
	m["*"] = multiply
	m["/"] = divide
	m["^"] = power
	m["d"] = strRoll

	// Format the input string, removing spaces and lowering case
	argString = strings.Replace(strings.ToLower(argString), " ", "", -1)

	// Use argString to create terms
	termSlice := parse(argString)

	// Create a channel with a buffer as long as the term slice
	outPut := make(chan string)

	// Prepare to evaluate the terms
	postfix := CreatePostfix(termSlice)
	resultString := argString

	// Keep an eye out for die rolls
	rollPattern, err := regexp.Compile("\\d*d\\d+")
	checkErr(err)

	// GoRoutine to put roll results into the output string
	go func() {
		for str := range outPut {
			resultString = strings.Replace(resultString, rollPattern.FindString(resultString), str, 1)
		}
	}()

	// Evaluate the terms
	resultSlice := evalPostfix(postfix, m, outPut)

	// Force the program to wait until every instance of "d[number]" is replaced
	for rollPattern.MatchString(resultString) {
	}

	// Return the results
	return resultString + " = " + resultSlice[0]
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
			if str == "d" {
				terms = append(terms, "1")
			}
			terms = append(terms, str)
		}
	}
	if num != "" {
		terms = append(terms, num)
	}

	terms = append(terms, ")")

	return terms
}

func evalPostfix(postfix *Postfix, m map[string]func(string, string, chan string) string, outPut chan string) []string {
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
			outStack, rightNum = Pop(outStack)
			outStack, leftNum = Pop(outStack)
			// fmt.Println(leftNum, rightNum)
			outStack = append(outStack, fn(leftNum, rightNum, outPut))
		}
		postfix.data[0] = ""
		if len(postfix.data) > 0 {
			postfix.data = postfix.data[1:]
		}
	}
	return outStack
}

func add(leftNum, rightNum string, outPut chan string) string {
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
	// c <- "+" + rightNum
	return result
}

func subtract(leftNum, rightNum string, outPut chan string) string {
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
	// c <- "-" + rightNum
	return result
}

func multiply(leftNum, rightNum string, outPut chan string) string {
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
	// c <- "*" + rightNum
	return result
}

func divide(leftNum, rightNum string, outPut chan string) string {
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
	// c <- "/" + rightNum
	return result
}

func power(leftNum, rightNum string, outPut chan string) string {
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
	// fmt.Println(left, right, result)
	// c <- "^" + rightNum
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

func strRoll(leftNum, rightNum string, outPut chan string) string {
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
	// fmt.Println(strRolls[:lastIndex])

	// Add each roll to the output channel
	outPut <- "(" + strRolls[:lastIndex] + ")"
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
