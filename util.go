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
