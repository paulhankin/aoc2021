package main

import (
	"bufio"
	"os"
	"strings"
)

func readDay15() ([][]int, error) {
	f, err := os.Open("day15.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var r [][]int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var line []int
		for _, b := range strings.TrimSpace(scanner.Text()) {
			line = append(line, int(b-'0'))
		}
		r = append(r, line)
	}
	return r, scanner.Err()
}

func minPath(lines [][]int) int {
	m := len(lines)
	n := len(lines[0])
	c2 := func(i, j int) int {
		return i*n + j
	}
	c2i := func(i int) (int, int) {
		return i / n, i % n
	}
	adj := func(x int, r []NodeCost) []NodeCost {
		i, j := c2i(x)
		for d := 0; d < 4; d++ {
			di, dj := dir4(d)
			if i+di < 0 || i+di >= m || j+dj < 0 || j+dj >= n {
				continue
			}
			v := c2(i+di, j+dj)
			r = append(r, NodeCost{v, int(lines[i+di][j+dj])})
		}
		return r
	}
	heuristic := func(x int) int {
		i, j := c2i(x)
		return (m - 1 - i) + (n - 1 - j)
	}
	return MinPath(c2(0, 0), c2(m-1, n-1), adj, heuristic, false)
}

func expand15(x [][]int) [][]int {
	m := len(x)
	n := len(x[0])
	var result [][]int
	for i := 0; i < m*5; i++ {
		var line []int
		for j := 0; j < n*5; j++ {
			r := x[i%m][j%n]
			r += i/m + j/n
			r = (r-1)%9 + 1
			line = append(line, r)
		}
		result = append(result, line)
	}
	return result
}

func day15() error {
	lines, err := readDay15()
	if err != nil {
		return err
	}
	lines2 := expand15(lines)
	partPrint(1, minPath(lines))
	partPrint(2, minPath(lines2))
	return nil
}

func init() {
	RegisterDay(15, day15)
}
