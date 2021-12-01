package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func readReport() ([]int, error) {
	f, err := os.Open("day1.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to read report: %v", err)
	}
	defer f.Close()
	var r []int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		x, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("failed to parse int from report: %v", err)
		}
		r = append(r, int(x))
	}
	return r, nil
}

func intSum(d []int) int {
	s := 0
	for _, v := range d {
		s += v
	}
	return s
}

func day1() error {
	data, err := readReport()
	if err != nil {
		return err
	}
	incs := 0
	inc3s := 0
	for i, v := range data {
		if i >= 1 && v > data[i-1] {
			incs++
		}
		if i >= 3 && intSum(data[i-2:i+1]) > intSum(data[i-3:i]) {
			inc3s++
		}
	}
	fmt.Println("part 1 =", incs)
	fmt.Println("part 3 =", inc3s)
	return nil
}

func init() {
	RegisterDay(1, day1)
}
