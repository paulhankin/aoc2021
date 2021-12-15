package main

import (
	"bufio"
	"container/heap"
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

type pointDist struct {
	i    int
	dist int
}

type distHeap struct {
	h   []*pointDist
	idx []int // map from i to index
}

func (d *distHeap) Len() int {
	return len(d.h)
}
func (d *distHeap) Less(i, j int) bool {
	return d.h[i].dist < d.h[j].dist
}
func (d *distHeap) Swap(i, j int) {
	d.h[i], d.h[j] = d.h[j], d.h[i]
	d.idx[d.h[i].i] = i
	d.idx[d.h[j].i] = j
}
func (d *distHeap) Push(x interface{}) {
	xpd := x.(*pointDist)
	d.idx[xpd.i] = len(d.h)
	d.h = append(d.h, xpd)
}

func (d *distHeap) Pop() interface{} {
	old := d.h
	n := len(old)
	x := old[n-1]
	d.h = old[0 : n-1]
	d.idx[x.i] = -1
	return x
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
	prev := make([]int, m*n)
	Q := &distHeap{
		h:   make([]*pointDist, m*n),
		idx: make([]int, m*n),
	}
	infinity := 1 << 62
	dist := make([]pointDist, m*n)
	for i := 0; i < m*n; i++ {
		prev[i] = -1
		dist[i] = pointDist{i, infinity}
		Q.h[i] = &dist[i]
		Q.idx[i] = i
	}
	Q.h[c2(0, 0)].dist = 0
	heap.Init(Q)
	for Q.Len() > 0 {
		u := heap.Pop(Q).(*pointDist)
		i, j := c2i(u.i)
		for d := 0; d < 4; d++ {
			di, dj := dir4(d)
			if i+di < 0 || i+di >= m || j+dj < 0 || j+dj >= n {
				continue
			}
			v := c2(i+di, j+dj)
			if Q.idx[v] == -1 {
				continue
			}
			alt := u.dist + lines[i+di][j+dj]
			if alt < dist[v].dist {
				vi := Q.idx[v]
				heap.Remove(Q, vi)
				dist[v].dist = alt
				prev[v] = u.i
				heap.Push(Q, &dist[v])
			}
		}
	}
	return dist[m*n-1].dist
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
