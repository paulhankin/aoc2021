package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type fold struct {
	v    int
	axis int // 0=x, 1=y
}

type fold13 struct {
	points []coord2d
	folds  []fold
}

func readDay13() (*fold13, error) {
	f, err := os.Open("day13.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var f13 fold13
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == "" {
			continue
		}
		if strings.HasPrefix(scanner.Text(), "fold along") {
			var axisc rune
			var v int
			_, err := fmt.Sscanf(scanner.Text(), "fold along %c=%d", &axisc, &v)
			if err != nil {
				return nil, err
			}
			if axisc != 'x' && axisc != 'y' {
				return nil, fmt.Errorf("axis %c in %q", axisc, scanner.Text())
			}
			var axis = 0
			if axisc == 'y' {
				axis = 1
			}
			f13.folds = append(f13.folds, fold{v: v, axis: axis})
		} else {
			var x, y int
			_, err := fmt.Sscanf(scanner.Text(), "%d,%d", &x, &y)
			if err != nil {
				return nil, err
			}
			f13.points = append(f13.points, coord2d{x, y})
		}
	}
	return &f13, scanner.Err()
}

func fold1d(x, v int) int {
	if x > v {
		return 2*v - x
	}
	return x
}

func applyFold(p coord2d, f fold) coord2d {
	if f.axis == 0 {
		return coord2d{fold1d(p.x, f.v), p.y}
	} else {
		return coord2d{p.x, fold1d(p.y, f.v)}
	}
}

func renderPoints(ps map[coord2d]bool) string {
	Mx, My := -1000000, -1000000
	for p := range ps {
		Mx = maxint(Mx, p.x)
		My = maxint(My, p.y)
	}
	rows := make([][]byte, My+1)
	for i := range rows {
		rows[i] = make([]byte, Mx+1)
	}
	for _, r := range rows {
		for j := range r {
			r[j] = ' '
		}
	}
	for p := range ps {
		rows[p.y][p.x] = '#'
	}
	var sb strings.Builder
	for _, r := range rows {
		sb.WriteString("\n")
		sb.Write(r)
	}
	return sb.String()
}

func day13() error {
	inst13, err := readDay13()
	if err != nil {
		return err
	}
	points := map[coord2d]bool{}
	for _, p := range inst13.points {
		points[p] = true
	}
	for part := 1; part <= 2; part++ {
		for _, f := range inst13.folds {
			np := map[coord2d]bool{}
			for p := range points {
				np[applyFold(p, f)] = true
			}
			points = np
			if part == 1 {
				break
			}
		}
		if part == 1 {
			partPrint(part, len(points))
		} else {
			partPrint(part, renderPoints(points))
		}
	}
	return nil
}

func init() {
	RegisterDay(13, day13)
}
