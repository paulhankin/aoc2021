package main

import (
	"os"
	"log"
	"bufio"
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
	if n, ok := heights[i2(i-1, j)]; ok && n <= b {
		return false
	}
	if n, ok := heights[i2(i+1, j)]; ok && n <= b {
		return false
	}
	if n, ok := heights[i2(i, j-1)]; ok && n <= b {
		return false
	}
	if n, ok := heights[i2(i, j+1)]; ok && n <= b {
		return false
	}
	return true
}


type UnionFind2i struct {
	Nodes []int
	NodeIndex map[[2]int]int
	Names [][2]int
}

func newUnionFind2i() *UnionFind2i {
	return &UnionFind2i{
		NodeIndex: map[[2]int]int{},
	}
}

func (uf *UnionFind2i) Add(x [2]int) {
	if _, ok := uf.NodeIndex[x]; ok {
		return
	}
	uf.NodeIndex[x] = len(uf.Nodes)
	uf.Nodes = append(uf.Nodes, len(uf.Nodes))
	uf.Names = append(uf.Names, x)
}

func (uf *UnionFind2i) Root(x [2]int) [2]int {
	idx, ok := uf.NodeIndex[x]
	if !ok {
		log.Fatal("root of non-existent node")
	}
	r := idx
	for uf.Nodes[r] != r {
		r = uf.Nodes[r]
	}
	for idx != r {
		p := uf.Nodes[idx]
		uf.Nodes[idx] = r
		idx = p
	}
	return uf.Names[r]
}

func (uf *UnionFind2i) Union(x, y [2]int) {
	rx := uf.NodeIndex[uf.Root(x)]
	ry := uf.NodeIndex[uf.Root(y)]
	if rx != ry {
		uf.Nodes[rx] = ry
	}
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

	dirs := [][2]int {
		i2(1, 0),
		i2(-1, 0),
		i2(0, 1),
		i2(0, -1),
	}

	uf := newUnionFind2i()
	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			if heights[i2(i, j)] == 9 {
				continue
			}
			uf.Add(i2(i, j))
		}
	}

	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			if heights[i2(i, j)] == 9 {
				continue
			}
			for _, dij := range dirs {
				ii := i + dij[0]
				jj := j + dij[1]
				if ii < 0 || ii >= M || jj < 0 || jj >= N {
					continue
				}
				if heights[i2(ii, jj)] <= heights[i2(i, j)] {
					uf.Union(i2(ii, jj), i2(i, j))
				}
			}
		}
	}
	sizes := map[[2]int]int{}
	for i := 0; i < M; i++ {
		for j := 0; j < N; j++ {
			if heights[i2(i, j)] == 9 {
				continue
			}
			sizes[uf.Root(i2(i, j))]++
		}
	}
	var sizeSlice []int
	for _, s := range sizes {
		sizeSlice = append(sizeSlice, s)
	}
	sort.Ints(sizeSlice)
	n := len(sizeSlice) - 1
	partPrint(2, sizeSlice[n] * sizeSlice[n-1] * sizeSlice[n-2])

	return nil
}

func init() {
	RegisterDay(9, day9)
}
