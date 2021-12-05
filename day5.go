package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type coord2d struct {
	x int
	y int
}

type vent struct {
	from, to coord2d
}

func (c coord2d) String() string {
	return fmt.Sprintf("%d,%d", c.x, c.y)
}

func (v vent) String() string {
	return fmt.Sprintf("[%s->%s]", v.from, v.to)
}

func parseCoord2d(s string) (coord2d, error) {
	coords, err := parseInts(s, ",")
	if err != nil || len(coords) != 2 {
		return coord2d{}, fmt.Errorf("failed to parse coords %q", s)
	}
	return coord2d{x: coords[0], y: coords[1]}, nil
}

func readDay5() ([]vent, error) {
	f, err := os.Open("day5.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var r []vent
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "->")
		if len(parts) != 2 {
			return nil, fmt.Errorf("can't find single -> in %q", scanner.Text())
		}
		from, ferr := parseCoord2d(parts[0])
		to, terr := parseCoord2d(parts[1])
		if ferr != nil || terr != nil {
			return nil, fmt.Errorf("failed to parse vent: %v or %v", ferr, terr)
		}
		r = append(r, vent{from: from, to: to})
	}
	return r, scanner.Err()
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

func day5() error {
	vents, err := readDay5()
	if err != nil {
		return err
	}
	for part := 1; part <= 2; part++ {
		board := map[coord2d]int{}
		for _, v := range vents {
			dx := v.to.x - v.from.x
			dy := v.to.y - v.from.y
			if part == 1 && dx != 0 && dy != 0 {
				continue
			}
			for i := 0; (i-1)*sgn(dx)+v.from.x != v.to.x || (i-1)*sgn(dy)+v.from.y != v.to.y; i++ {
				board[coord2d{v.from.x + i*sgn(dx), v.from.y + i*sgn(dy)}]++
			}
		}
		sum := 0
		for _, b := range board {
			if b > 1 {
				sum++
			}
		}
		fmt.Println("part", part, "=", sum)
	}
	return nil
}

func init() {
	RegisterDay(5, day5)
}
