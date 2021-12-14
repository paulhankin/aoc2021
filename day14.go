package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type reactions struct {
	start []byte
	rs    [][3]byte
	rsmap map[[2]byte]byte
}

func step14(pairs map[[2]byte]int, r *reactions) map[[2]byte]int {
	y := map[[2]byte]int{}
	for p, n := range pairs {
		if got := r.rsmap[p]; got != 0 {
			y[[2]byte{p[0], got}] += n
			y[[2]byte{got, p[1]}] += n
		} else {
			y[p] += n
		}
	}
	return y
}

func readDay14() (*reactions, error) {
	f, err := os.Open("day14.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := reactions{
		rsmap: map[[2]byte]byte{},
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(scanner.Text(), " -> ")
		if len(parts) == 1 {
			if len(r.start) != 0 {
				return nil, fmt.Errorf("two start lines: second %q", line)
			}
			r.start = []byte(line)
		} else if len(parts) != 2 || len(parts[0]) != 2 || len(parts[1]) != 1 {
			return nil, fmt.Errorf("can't parse reaction line %q", line)
		} else {
			r.rs = append(r.rs, [3]byte{parts[0][0], parts[0][1], parts[1][0]})
			r.rsmap[[2]byte{parts[0][0], parts[0][1]}] = parts[1][0]
		}
	}
	return &r, scanner.Err()
}

func day14() error {
	r, err := readDay14()
	if err != nil {
		return err
	}
	for part := 1; part <= 2; part++ {
		pairs := map[[2]byte]int{}
		getx := func(i int) byte {
			if i < 0 || i >= len(r.start) {
				return '*'
			}
			return r.start[i]
		}
		for i := -1; i < len(r.start); i++ {
			pairs[[2]byte{getx(i), getx(i + 1)}]++
		}

		steps := 10 + 30*b2i(part == 2)
		for i := 0; i < steps; i++ {
			pairs = step14(pairs, r)
		}
		counts := map[byte]int{}
		for p, n := range pairs {
			if p[0] == '*' {
				continue
			}
			counts[p[0]] += n
		}
		m, M := 1<<62, 0
		for _, c := range counts {
			m = minint(m, c)
			M = maxint(M, c)
		}
		partPrint(part, M-m)
	}
	return nil
}

func init() {
	RegisterDay(14, day14)
}
