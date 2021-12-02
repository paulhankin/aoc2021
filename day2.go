package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type dxy struct {
	dx int
	dy int
}

func getDir(s string) (dxy, error) {
	if s == "forward" {
		return dxy{1, 0}, nil
	}
	if s == "down" {
		return dxy{0, 1}, nil
	}
	if s == "up" {
		return dxy{0, -1}, nil
	}
	return dxy{}, fmt.Errorf("bad dir %q", s)
}

func readDay2() ([]dxy, error) {
	f, err := os.Open("day2.txt")
	if err != nil {
		return nil, err
	}
	var r []dxy
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) != 2 {
			return nil, fmt.Errorf("bad line %q (wanted 2 parts, got %d)", scanner.Text(), len(parts))
		}
		dir, err := getDir(parts[0])
		if err != nil {
			return nil, err
		}
		amt, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("bad number %q: %v", parts[1], err)
		}
		r = append(r, dxy{dir.dx * amt, dir.dy * amt})
	}
	return r, scanner.Err()
}

func day2() error {
	dirs, err := readDay2()
	if err != nil {
		return err
	}
	hor := 0
	ver := 0
	for _, d := range dirs {
		hor += d.dx
		ver += d.dy
	}
	fmt.Println("part 1 =", hor*ver)

	hor = 0
	ver = 0
	aim := 0
	for _, d := range dirs {
		aim += d.dy
		hor += d.dx
		ver += d.dx * aim
	}
	fmt.Println("part 2 =", hor*ver)

	return nil
}

func init() {
	RegisterDay(2, day2)
}
