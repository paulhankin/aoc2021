package main

import (
	"bufio"
	"os"
)

func readDay6() ([]int, error) {
	f, err := os.Open("day6.txt")
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	var r []int
	for scanner.Scan() {
		xs, err := parseInts(scanner.Text(), ",")
		if err != nil {
			return nil, err
		}
		r = append(r, xs...)
	}
	return r, scanner.Err()
}

func day6() error {
	fish, err := readDay6()
	if err != nil {
		return err
	}
	for part := 1; part <= 2; part++ {
		timers := make([]uint64, 9)
		for _, f := range fish {
			timers[f]++
		}
		maxDay := []int{0, 80, 256}[part]
		for days := 0; days < maxDay; days++ {
			var nt []uint64
			nt = append(nt, timers[1:]...)
			nt = append(nt, 0)
			nt[6] += timers[0]
			nt[8] += timers[0]
			timers = nt
		}
		var sum uint64
		for _, t := range timers {
			sum += t
		}
		partPrint(part, sum)
	}

	return nil
}

func init() {
	RegisterDay(6, day6)
}
