package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type vent struct {
	from, to coord2d
}

func (v vent) String() string {
	return fmt.Sprintf("[%s->%s]", v.from, v.to)
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
			if dx != 0 && dy != 0 && dx != dy && dx != -dy {
				log.Fatal("found non-horizontal/vertical/diagonal line", v)
			}
			for i := 0; (i-1)*sgn(dx)+v.from.x != v.to.x || (i-1)*sgn(dy)+v.from.y != v.to.y; i++ {
				board[coord2d{v.from.x + i*sgn(dx), v.from.y + i*sgn(dy)}]++
			}
		}
		sum := 0
		for _, b := range board {
			sum += clamp(b-1, 0, 1)
		}
		fmt.Println("part", part, "=", sum)
	}
	return nil
}

func init() {
	RegisterDay(5, day5)
}
