package main

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

func readDay9() ([]string, error) {
	f, err := os.Open("day9.txt")
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

func i2(i, j int) [2]int {
	return [2]int{i, j}
}

func lowPoint(heights map[[2]int]int, i, j int) bool {
	b := heights[i2(i, j)]
	for d := 0; d < 4; d++ {
		di, dj := dir4(d)
		if n, ok := heights[i2(i+di, j+dj)]; ok && n <= b {
			return false
		}
	}
	return true
}

func day9() error {
	lines, err := readDay9()
	if err != nil {
		return err
	}
	M := len(lines)
	N := len(lines[0])
	heights := map[[2]int]int{}
	nlp := 0
	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			heights[i2(i, j)] = int(lines[i][j] - '0')
		}
	}
	risks := 0
	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			if lowPoint(heights, i, j) {
				nlp++
				risks += heights[i2(i, j)] + 1
			}
		}
	}
	partPrint(1, risks)

	uf := NewUnionFind(M * N)

	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			if heights[i2(i, j)] == 9 {
				continue
			}
			for dd := 0; dd < 4; dd++ {
				di, dj := dir4(dd)
				if i+di < 0 || i+di >= M || j+dj < 0 || j+dj >= N {
					continue
				}
				if heights[i2(i+di, j+dj)] <= heights[i2(i, j)] {
					uf.Union((i+di)*N+(j+dj), i*N+j)
				}
			}
		}
	}
	sizes := map[int]int{}
	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			if heights[i2(i, j)] == 9 {
				continue
			}
			sizes[uf.Find(i*N+j)]++
		}
	}
	var sizeSlice []int
	for _, s := range sizes {
		sizeSlice = append(sizeSlice, s)
	}
	sort.Ints(sizeSlice)
	n := len(sizeSlice) - 1
	partPrint(2, sizeSlice[n]*sizeSlice[n-1]*sizeSlice[n-2])

	return nil
}

func init() {
	RegisterDay(9, day9)
}
