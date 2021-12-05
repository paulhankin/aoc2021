package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type bingoData struct {
	draw   []int
	boards [][]int
}

func parseInts(s, sep string) ([]int, error) {
	s = strings.TrimSpace(s)
	parts := strings.FieldsFunc(s, func(r rune) bool { return strings.ContainsRune(sep, r) })
	var r []int
	for _, part := range parts {
		x, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q: %v", part, err)
		}
		r = append(r, x)
	}
	return r, nil
}

func readDraw(scanner *bufio.Scanner) ([]int, error) {
	scanner.Scan()
	return parseInts(scanner.Text(), ",")
}

func readBoard(scanner *bufio.Scanner) ([]int, error) {
	var board []int
	for scanner.Scan() {
		xs, err := parseInts(scanner.Text(), " ")
		if err != nil {
			return nil, err
		}
		if len(xs) == 0 {
			continue
		}
		if len(xs) != 5 {
			return nil, fmt.Errorf("expected 5 bingo numbers but found %q", scanner.Text())
		}
		board = append(board, xs...)
		if len(board) == 25 {
			return board, scanner.Err()
		}
	}
	if len(board) == 0 {
		return nil, nil
	}
	return nil, fmt.Errorf("scanner ended with only %d numbers in board", len(board))
}

func readDay4() (bingoData, error) {
	f, err := os.Open("day4.txt")
	if err != nil {
		return bingoData{}, err
	}
	defer f.Close()
	var bingo bingoData

	scanner := bufio.NewScanner(f)

	bingo.draw, err = readDraw(scanner)
	if err != nil {
		return bingo, err
	}
	for {
		board, err := readBoard(scanner)
		if err != nil {
			return bingo, err
		}
		if board == nil {
			return bingo, scanner.Err()
		}
		bingo.boards = append(bingo.boards, board)
	}
}

func scoreBingo(b []int, lastCall int, calls map[int]bool) int {
	sum := 0
	for _, i := range b {
		if !calls[i] {
			sum += i
		}
	}
	return sum * lastCall
}

func isBingoWon(b []int, calls map[int]bool) bool {
	for i := 0; i < 5; i++ {
		row := 0
		col := 0
		for j := 0; j < 5; j++ {
			if calls[b[i*5+j]] {
				row++
			}
			if calls[b[j*5+i]] {
				col++
			}
		}
		if row == 5 || col == 5 {
			return true
		}
	}
	return false
}

func playBingo(bingo bingoData) (int, error) {
	calls := map[int]bool{}
	for _, call := range bingo.draw {
		calls[call] = true
		for _, b := range bingo.boards {
			if isBingoWon(b, calls) {
				return scoreBingo(b, call, calls), nil
			}
		}
	}
	return 0, fmt.Errorf("no board won?")
}

func playBingoToLose(bingo bingoData) (int, error) {
	won := map[int]bool{}
	calls := map[int]bool{}
	for _, call := range bingo.draw {
		calls[call] = true
		for bi, b := range bingo.boards {
			if won[bi] {
				continue
			}
			if isBingoWon(b, calls) {
				won[bi] = true
			}
			if len(won) == len(bingo.boards) {
				return scoreBingo(b, call, calls), nil
			}
		}
	}
	return 0, fmt.Errorf("not all boards won?")
}

func day4() error {
	bingo, err := readDay4()
	if err != nil {
		return err
	}
	score, err := playBingo(bingo)
	if err != nil {
		return err
	}
	fmt.Println("part 1 =", score)
	score2, err := playBingoToLose(bingo)
	if err != nil {
		return err
	}
	fmt.Println("part 2 =", score2)
	return nil
}

func init() {
	RegisterDay(4, day4)
}
