package main

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

func readDay9() ([][]int, error) {
	f, err := os.Open("day9.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var lines [][]int
	for scanner.Scan() {
		var line []int
		for _, b := range strings.TrimSpace(scanner.Text()) {
			line = append(line, int(b-'0'))
		}
		lines = append(lines, line)
	}
	return lines, scanner.Err()
}

func day9() error {
	h, err := readDay9()
	if err != nil {
		return err
	}
	M := len(h)
	N := len(h[0])

	uf := NewUnionFind(M * N)

	risks := 0
	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			if h[i][j] == 9 {
				continue
			}
			isLow := true
			for dd := 0; dd < 4; dd++ {
				di, dj := dir4(dd)
				if i+di < 0 || i+di >= M || j+dj < 0 || j+dj >= N {
					continue
				}
				if h[i+di][j+dj] <= h[i][j] {
					uf.Union((i+di)*N+(j+dj), i*N+j)
					isLow = false
				}
			}
			if isLow {
				risks += h[i][j] + 1
			}
		}
	}
	partPrint(1, risks)

	sizes := map[int]int{}
	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			if h[i][j] == 9 {
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
