package main

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

func readDay10() ([]string, error) {
	f, err := os.Open("day10.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	return lines, scanner.Err()
}

var brakPair = map[byte]byte{
	'{': '}',
	'[': ']',
	'(': ')',
	'<': '>',
}

func openBracket(b byte) bool {
	return brakPair[b] != 0
}
func closeBracket(b byte) bool {
	return !openBracket(b)
}

var brakScore = map[byte]int{
	')': 3,
	']': 57,
	'}': 1197,
	'>': 25137,
}

func scoreBraks(s string) ([]byte, int) {
	stack := make([]byte, 0, 1024)

	for i := 0; i < len(s); i++ {
		b := s[i]
		if openBracket(b) {
			stack = append(stack, b)
		} else {
			N := len(stack) - 1
			if brakPair[stack[N]] == b {
				stack = stack[:N]
			} else {
				return nil, brakScore[b]
			}

		}
	}
	return stack, 0
}

var scoreComplete = map[byte]int{
	'(': 1,
	'[': 2,
	'{': 3,
	'<': 4,
}

func day10() error {
	braks, err := readDay10()
	if err != nil {
		return err
	}
	var totalScore int
	var totalScore2 []int
	for _, b := range braks {
		complete, score := scoreBraks(b)
		if score != 0 {
			totalScore += score
		} else {
			ts2 := 0
			for n := len(complete) - 1; n >= 0; n-- {
				ts2 = 5*ts2 + scoreComplete[complete[n]]
			}
			totalScore2 = append(totalScore2, ts2)
		}
	}
	sort.Ints(totalScore2)
	partPrint(1, totalScore)
	partPrint(2, totalScore2[len(totalScore2)/2])
	return nil
}

func init() {
	RegisterDay(10, day10)
}
