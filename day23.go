package main

import (
	"fmt"
	"log"
)

type state23 struct {
	// 7 on the top row, 4 on the next row, 4 on the next row
	b [15]byte
}

func (x state23) String() string {
	s := x.b
	top := fmt.Sprintf("%c%c.%c.%c.%c.%c%c", s[0], s[1], s[2], s[3], s[4], s[5], s[6])
	mid := fmt.Sprintf("  %c %c %c %c", s[7], s[8], s[9], s[10])
	bot := fmt.Sprintf("  %c %c %c %c", s[11], s[12], s[13], s[14])
	return top + "\n" + mid + "\n" + bot
}

func newState23(s string) state23 {
	var r state23
	if len(s) != 8 {
		log.Fatalf("bad input string %q", s)
	}
	for i := 0; i < 15; i++ {
		if i < 7 {
			r.b[i] = '.'
		} else {
			r.b[i] = s[i-7]
		}
	}
	return r
}

func (s *state23) getCorridor(i int) byte {
	if i == 2 || i == 4 || i == 6 || i == 8 {
		return '.'
	}
	i -= b2i(i >= 3) + b2i(i >= 5) + b2i(i >= 7) + b2i(i >= 9)
	return s.b[i]
}

func topToCorridor(i int) int {
	return i + b2i(i >= 2) + b2i(i >= 3) + b2i(i >= 4) + b2i(i >= 5)
}

func (s *state23) locSlugPath(i int, dest int) int {
	cost := 0
	loc := i
	for loc != dest {
		loc += sgn(dest - loc)
		if s.getCorridor(loc) != '.' {
			return 0
		}
		cost += 1
	}
	return cost
}

func slugCost(x byte) int {
	if x == 'A' {
		return 1
	} else if x == 'B' {
		return 10
	} else if x == 'C' {
		return 100
	} else if x == 'D' {
		return 1000
	}
	log.Fatalf("bad slug %c", x)
	return 0
}

// findSlugPath finds an x slug that's on the top row, that can move to
// the "home" location dest.
func (s *state23) findSlugPath(x byte, dest int) (int, int) {
	for i := 0; i < 7; i++ {
		if s.b[i] != x {
			continue
		}
		pathlen := s.locSlugPath(topToCorridor(i), (dest-7)%4*2+2)
		if pathlen == 0 {
			continue
		}
		pathlen += b2i(dest >= 7) + b2i(dest >= 11)
		return i, pathlen * slugCost(x)
	}
	return 0, 0
}

// normalize23 moves any slugs into their homes if possible as it's always wlog.
// It returns the new state, and the cost of the moves.
func normalize23(s state23) (state23, int) {
	tcost := 0
	s0 := s
	for changed := true; changed; {
		changed = false
		for i := 0; i < 4; i++ {
			slug := 'A' + byte(i)
			if s.b[11+i] == '.' {
				if s.b[7+i] != '.' {
					fmt.Println(s)
					log.Fatalf("found empty square on bottom home, but above square _not_ empty!")
				}
				if loc, cost := s.findSlugPath(slug, 11+i); cost > 0 {
					s.b[loc] = '.'
					s.b[11+i] = slug
					tcost += cost
					changed = true
				}
			}
			if s.b[11+i] == slug && s.b[7+i] == '.' {
				if loc, cost := s.findSlugPath(slug, 7+i); cost > 0 {
					s.b[loc] = '.'
					s.b[7+i] = slug
					tcost += cost
					changed = true
				}
			}
		}
	}
	if !s.ok() {
		log.Fatalf("normalized\n%s\nto\n%s", s0, s)
	}
	return s, tcost
}

// stateCost23 returns a state and corresponding cost.
type stateCost23 struct {
	s    state23
	cost int
}

// adjacent returns all states adjacent to ours.
// The starting state must be normalized, and the adjacent
// states are also normalized.
func (s state23) adjacent() []stateCost23 {
	var r []stateCost23
	for i := 0; i < 4; i++ {
		hs := 'A' + byte(i)
		if s.b[7+i] == '.' {
			continue
		}
		if s.b[7+i] == hs && s.b[11+i] == hs {
			continue
		}
		for j := 0; j < 7; j++ {
			if cost := s.locSlugPath(i*2+2, topToCorridor(j)); cost > 0 {
				ns := s
				slug := ns.b[7+i]
				ns.b[j] = slug
				ns.b[7+i] = '.'
				if !ns.ok() {
					log.Fatalf("moved from middle\n%s\nto\n%s", s, ns)
				}
				ns, nc := normalize23(ns)
				r = append(r, stateCost23{ns, (cost+1)*slugCost(slug) + nc})
			}
		}
	}

	for i := 0; i < 4; i++ {
		hs := 'A' + byte(i)
		if s.b[11+i] == '.' {
			continue
		}
		if s.b[11+i] == hs || s.b[7+i] != '.' {
			continue
		}
		for j := 0; j < 7; j++ {
			if cost := s.locSlugPath(i*2+2, topToCorridor(j)); cost > 0 {
				ns := s
				slug := ns.b[11+i]
				ns.b[j] = slug
				ns.b[11+i] = '.'
				if !ns.ok() {
					log.Fatalf("moved from bottom %d->%d\n%s\nto\n%s", i, j, s, ns)
				}
				ns, nc := normalize23(ns)
				r = append(r, stateCost23{ns, (cost+2)*slugCost(slug) + nc})
			}
		}
	}

	return r
}

func (s state23) ok() bool {
	var cs [4]int
	for _, c := range s.b {
		if c != '.' {
			cs[c-'A']++
		}
	}
	return cs == [4]int{2, 2, 2, 2}
}

func solve23(s state23, m map[state23]int) int {
	if cost, ok := m[s]; ok {
		return cost
	}
	const infinity = 99999999
	best := infinity
	for _, sc := range s.adjacent() {
		if !sc.s.ok() {
			log.Fatalf("from\n%s\nto\n%s\n\n", s, sc.s)
		}
		c := solve23(sc.s, m) + sc.cost
		if c < best {
			best = c
		}
	}
	return best
}

func day23s(s state23) {
	m := map[state23]int{}
	// insert the solved state
	m[newState23("ABCDABCD")] = 0
	partPrint(1, solve23(s, m))
}

func init() {
	s0example := newState23("BCBDADCA")
	s0 := newState23("BCADBCDA")
	_ = s0
	_ = s0example
	RegisterDay(23, func() error {
		day23s(s0)
		return nil
	})
}
