package main

import (
	"bufio"
	"os"
	"strings"
)

type cukeBoard struct {
	m, n int
	locs []byte
}

func newCukeBoard(m, n int) *cukeBoard {
	return &cukeBoard{m, n, make([]byte, m*n)}
}

func (b2 *cukeBoard) Reset() {
	for i := range b2.locs {
		b2.locs[i] = '.'
	}
}

func (cb *cukeBoard) set(i, j int, b byte) {
	i = (i + cb.m) % cb.m
	j = (j + cb.n) % cb.n
	cb.locs[i*cb.n+j] = b
}

func (cb *cukeBoard) get(i, j int) byte {
	i = (i + cb.m) % cb.m
	j = (j + cb.n) % cb.n
	return cb.locs[i*cb.n+j]
}

func readDay25(s string) (*cukeBoard, error) {
	f, err := os.Open(s)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	m := len(lines)
	n := len(lines[0])
	b := newCukeBoard(m, n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			b.set(i, j, lines[i][j])
		}
	}
	return b, nil
}

func (b1 *cukeBoard) step(b2 *cukeBoard) bool {
	stuck := true
	m, n := b1.m, b1.n
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if b1.get(i, j) != '>' {
				continue
			}
			if b1.get(i, j+1) == '.' {
				b2.set(i, j+1, '>')
				stuck = false
			} else {
				b2.set(i, j, '>')
			}
		}
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if b1.get(i, j) != 'v' {
				continue
			}
			if b1.get(i+1, j) != 'v' && b2.get(i+1, j) != '>' {
				b2.set(i+1, j, 'v')
				stuck = false
			} else {
				b2.set(i, j, 'v')
			}
		}
	}
	return stuck
}

func day25(s string) error {
	board, err := readDay25(s)
	if err != nil {
		return err
	}
	i := 0
	b1 := board
	b2 := newCukeBoard(board.m, board.n)

	for {
		b2.Reset()
		stuck := b1.step(b2)
		b1, b2 = b2, b1
		i++
		if stuck {
			break
		}
	}
	partPrint(1, i)
	return nil
}

func init() {
	RegisterDay(25, func() error {
		return day25("day25.txt")
	})
}
