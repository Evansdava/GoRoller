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

	// fmt.Println(args)
	// (10d6+2d8)^2+(1d20-6d4)*12/2d4

	// // argString = "(10d6+2d8)^2+(1d20-6d4)*12/2d4"
	// argString := strings.ToLower(strings.Join(args, ""))
	// fmt.Println(GetRoll(argString))

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
	// m["k"] = keepHigh
	// m["kh"] = keepHigh
	// m["kl"] = keepLow
	// m["dh"] = dropHigh
	// m["dl"] = dropLow
	// m["d"] = dropLow

	// Format the input string, removing spaces and lowering case
	argString = strings.Replace(strings.ToLower(argString), " ", "", -1)

	// Use argString to create terms
	termSlice := parse(argString)

	// Create a channel with a buffer as long as the term slice
	outPut := make(chan string)
	dieRolls := make(chan string, len(termSlice))

	// Prepare to evaluate the terms
	postfix := CreatePostfix(termSlice)
	resultString := strings.Join(termSlice[1:len(termSlice)-1], "")

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
	resultSlice := evalPostfix(postfix, m, outPut, dieRolls)

	// Force the program to wait until every instance of "d[number]" is replaced
	for rollPattern.MatchString(resultString) {
	}

	resultString = strings.ReplaceAll(resultString, "*", "\\*")

	// Return the results
	if len(resultSlice) > 0 {
		return resultString + " = " + resultSlice[0]
	}
	return resultString
}

func parse(argString string) []string {
	terms := make([]string, 0)
	terms = append(terms, "(")

	num := ""
	for i := 0; i < len(argString); i++ {
		str := string(argString[i])
		// fmt.Println(str, i, len(argString))
		if _, err := strconv.Atoi(str); err == nil {
			num += str
		} else if num != "" && strings.Contains("+-*/^d()", str) {
			terms = append(terms, num)
			terms = append(terms, str)
			num = ""
		} else if strings.Contains("+-*/^d()", str) {
			if str == "d" && !strings.Contains(")", terms[len(terms)-1]) {
				// fmt.Println("Inserting 1 before operator")
				terms = append(terms, "1")
			} else if !strings.Contains(")", terms[len(terms)-1]) && !strings.Contains("()", str) {
				// fmt.Println("Inserting 0 before operator")
				terms = append(terms, "0")
			}
			terms = append(terms, str)
			if i == len(argString)-1 {
				// fmt.Println("Inserting 0 after operator")
				terms = append(terms, "0")
			}
		} // else if strings.Contains("kd", str) {
		// 	if i < len(argString)-1 {
		// 		if strings.Contains("hl", string(argString[i+1])) {
		// 			terms = append(terms, str+string(argString[i+1]))
		// 		} else {
		// 			terms = append(terms, str)
		// 		}
		// 	} else {
		// 		terms = append(terms, str)
		// 	}
		// }
	}
	if num != "" {
		terms = append(terms, num)
	}
	if strings.Contains("+-*/d^", terms[len(terms)-1]) {
		terms = append(terms, "0")
	}

	terms = append(terms, ")")
	fmt.Println(terms)

	return terms
}

func evalPostfix(postfix *Postfix, m map[string]func(string, string, chan string) string, outPut, dieRolls chan string) []string {
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

func strRoll(leftNum, rightNum string, dieRolls chan string) string {
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
	if len(strRolls) == 1 {
		dieRolls <- strRolls
		return strRolls
	}
	strRolls = strRolls[:lastIndex] + "=" + strRolls[lastIndex+1:]
	// result := strings.Join(strRolls[0:len(rolls)-1], "+")
	// fmt.Println(strRolls[:lastIndex])

	// Add each roll to the output channel
	dieRolls <- "(" + strRolls[:lastIndex] + ")"
	return strRolls[lastIndex+1:]
}

// func keepHigh(leftNum, rightNum string, dieRolls chan string) string {

// }

// func keepLow(leftNum, rightNum string, dieRolls chan string) string {

// }

// func dropLow(leftNum, rightNum string, dieRolls chan string) string {

// }

// func dropHigh(leftNum, rightNum string, dieRolls chan string) string {

// }

func dieRoll(numDice, dieSize int) []int {
	if numDice == 0 {
		return make([]int, 1)
	}
	if dieSize == 0 {
		return make([]int, numDice+1)
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
