package main

import (
	"fmt"
	"strconv"
	"strings"
)

type coord2d struct {
	x int
	y int
}

func (c coord2d) String() string {
	return fmt.Sprintf("%d,%d", c.x, c.y)
}

func parseCoord2d(s string) (coord2d, error) {
	coords, err := parseInts(s, ",")
	if err != nil || len(coords) != 2 {
		return coord2d{}, fmt.Errorf("failed to parse coords %q", s)
	}
	return coord2d{x: coords[0], y: coords[1]}, nil
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sgn(x int) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

func clamp(x, min, max int) int {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func i2b(x int) bool {
	return x != 0
}

func b2i(x bool) int {
	if x {
		return 1
	}
	return 0
}

func dir4(i int) (int, int) {
	i &= 3
	dx := b2i(i == 0) - b2i(i == 2)
	dy := b2i(i == 1) - b2i(i == 3)
	return dx, dy
}

func dir8(i int) (int, int) {
	i &= 7
	dx := b2i(i == 7 || i == 0 || i == 2) - b2i(i == 4 || i == 5 || i == 6)
	dy := b2i(i == 1 || i == 2 || i == 3) - b2i(i == 5 || i == 6 || i == 7)
	return dx, dy
}
