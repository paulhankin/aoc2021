package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type signals struct {
	patterns []string // always 10
	out      []string // always 4
}

func readDay8() ([]signals, error) {
	f, err := os.Open("day8.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var sigs []signals
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " | ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("can't find two halves of %q", scanner.Text())
		}
		left := strings.Split(parts[0], " ")
		right := strings.Split(parts[1], " ")
		if len(left) != 10 || len(right) != 4 {
			return nil, fmt.Errorf("wrong number of i/o of %q", scanner.Text())
		}
		sigs = append(sigs, signals{
			patterns: left,
			out:      right,
		})
	}
	return sigs, scanner.Err()
}

func sortString(s string) string {
	d := []rune(s)
	sort.Slice(d, func(i, j int) bool {
		return d[i] < d[j]
	})
	return string(d)

}

var digitSegs = [10]string{
	0: "abcefg",
	1: "cf",
	2: "acdeg",
	3: "acdfg",
	4: "bcdf",
	5: "abdfg",
	6: "abdefg",
	7: "acf",
	8: "abcdefg",
	9: "abcdfg",
}

func applyWiring(perm string, pat string) string {
	var r []byte
	for _, p := range pat {
		r = append(r, perm[p-'a'])
	}
	return sortString(string(r))
}

var wiringPerms = getWiringPerms()

func perms(xs []rune) [][]rune {
	var r [][]rune
	if len(xs) == 0 {
		return append(r, []rune{})
	}

	for i, x := range xs {
		xs[i], xs[0] = xs[0], xs[i]
		for _, tail := range perms(xs[1:]) {
			pp := append([]rune{x}, tail...)
			r = append(r, pp)
		}
		xs[i], xs[0] = xs[0], xs[i]
	}
	return r
}

func getWiringPerms() []string {
	var r []string
	for _, p := range perms([]rune("abcdefg")) {
		r = append(r, string(p))
	}
	return r
}

func solveWiring(pats []string) (map[string]int, error) {
mainLoop:
	for _, perm := range wiringPerms {
		var digits [10]string
		for _, pat := range pats {
			got := applyWiring(perm, pat)
			found := false
			for d := 0; d < 10; d++ {
				if digitSegs[d] == got {
					if digits[d] != "" {
						continue mainLoop
					}
					digits[d] = pat
					found = true
					break
				}
			}
			if !found {
				continue mainLoop
			}
		}
		r := map[string]int{}
		for d, s := range digits {
			r[sortString(s)] = d
		}
		return r, nil
	}
	return nil, fmt.Errorf("wiring not found")
}

func day8() error {
	sigs, err := readDay8()
	if err != nil {
		return err
	}
	d1478 := 0
	digits := [4]int{1, 4, 7, 8}
	for _, sig := range sigs {
		for _, out := range sig.out {
			for _, d := range digits {
				if len(out) == len(digitSegs[d]) {
					d1478++
				}
			}
		}
	}
	partPrint(1, d1478)
	var sum int
	for sign, sig := range sigs {
		mapping, err := solveWiring(sig.patterns)
		if err != nil {
			return err
		}
		var ds [4]int
		for i, d := range sig.out {
			var ok bool
			ds[i], ok = mapping[sortString(d)]
			if !ok {
				return fmt.Errorf("failed to find mapping in case %d", sign)
			}
		}
		sum += ds[0]*1000 + ds[1]*100 + ds[2]*10 + ds[3]
	}
	partPrint(2, sum)
	return nil
}

func init() {
	RegisterDay(8, day8)
}
