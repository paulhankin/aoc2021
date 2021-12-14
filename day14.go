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

func step14(x []byte, r *reactions) []byte {
	y := make([]byte, 0, len(x)*2)
	for i := 0; i < len(x)-1; i++ {
		y = append(y, x[i])
		if got := r.rsmap[[2]byte{x[i], x[i+1]}]; got != 0 {
			y = append(y, got)
		}
	}
	y = append(y, x[len(x)-1])
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
	x := append([]byte{}, r.start...)
	for i := 0; i < 10; i++ {
		x = step14(x, r)
	}
	counts := map[byte]int{}
	for _, b := range x {
		counts[b]++
	}
	m, M := 100000000, 0
	for _, c := range counts {
		m = minint(m, c)
		M = maxint(M, c)
	}
	partPrint(1, M-m)
	return nil
}

func init() {
	RegisterDay(14, day14)
}
