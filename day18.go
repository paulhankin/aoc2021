package main

import (
	"fmt"
	"regexp"
)

func (s *snNumber) magnitude() int {
	if s.Pair() {
		return s.p[0].magnitude()*3 + 2*s.p[1].magnitude()
	}
	return s.n
}

func findExtremeRegular(s *snNumber, first int) (*snNumber, bool) {
	if !s.Pair() {
		return s, true
	}
	if n, ok := findExtremeRegular(s.p[first], first); ok {
		return n, ok
	}
	return findExtremeRegular(s.p[1-first], first)
}

func findRegularLR(s *snNumber, first int) (*snNumber, bool) {
	if s.parent == nil {
		return nil, false
	}
	if !s.parent.Pair() {
		panic("pair child of non-pair?")
	}
	if s.parent.p[first] == s {
		// look on the left branch of the parent, if we're the right branch (first=1)
		// look on the right branch of the parent, if we're the left branch (first=0)
		r, ok := findExtremeRegular(s.parent.p[1-first], first)
		if ok {
			return r, ok
		}
	}
	return findRegularLR(s.parent, first)
}

func explode18(s *snNumber) {
	if !s.Pair() {
		panic("exploding pairs should be pairs!")
	}
	if s.p[0].Pair() || s.p[1].Pair() {
		panic("exploding pairs should always have regular numbers!")
	}
	n0 := s.p[0].n
	n1 := s.p[1].n
	left, leftok := findRegularLR(s, 1)
	right, rightok := findRegularLR(s, 0)
	if leftok {
		left.n += n0
	}
	if rightok {
		right.n += n1
	}
	s.p = nil
	s.n = 0
}

func reduce18explode(s *snNumber, nest int) bool {
	if !s.Pair() {
		return false
	}
	if nest >= 4 {
		explode18(s)
		return true
	}
	if reduce18explode(s.p[0], nest+1) {
		return true
	}
	if reduce18explode(s.p[1], nest+1) {
		return true
	}
	return false
}

func reduce18split(s *snNumber) bool {
	if !s.Pair() {
		if s.n >= 10 {
			s.p = &[2]*snNumber{
				{n: s.n / 2, parent: s}, {n: (s.n + 1) / 2, parent: s},
			}
			return true
		}
		return false
	}
	return reduce18split(s.p[0]) || reduce18split(s.p[1])
}

func reduce18(s *snNumber) bool {
	return reduce18explode(s, 0) || reduce18split(s)
}

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var re = regexp.MustCompile(ansi)

func stripAnsi(s string) string {
	return re.ReplaceAllString(s, "")
}

func printDiff(was, now string) {
	sim := 0
	wasPure := stripAnsi(was)
	nowPure := stripAnsi(now)
	for i := 0; i < len(nowPure); i++ {
		if i >= len(wasPure) || nowPure[i] != wasPure[i] {
			sim = i
			break
		}
	}
	siM := 0
	for i := 0; i < len(nowPure); i++ {
		if len(wasPure)-i-1 < 0 || nowPure[len(nowPure)-i-1] != wasPure[len(wasPure)-i-1] {
			siM = i
			break
		}
	}
	fmt.Println(now)
	if now != was {
		for i := 0; i < sim; i++ {
			fmt.Print(" ")
		}
		for i := 0; i < len(nowPure)-sim-siM; i++ {
			fmt.Print("^")
		}
		fmt.Println()
	}
}

var verbose = false

func add18(left, right *snNumber) *snNumber {
	if reduce18(left) {
		panic("left not reduced")
	}
	if reduce18(right) {
		panic("right not reduced")
	}
	root := &snNumber{p: &[2]*snNumber{left, right}}
	left.parent = root
	right.parent = root
	var was string
	if verbose {
		was = root.StringDepth(0, true)
		fmt.Println(was)
	}
	for reduce18(root) {
		now := root.StringDepth(0, true)
		if verbose {
			printDiff(was, now)
		}
		was = now
	}
	return root
}

func day18() error {
	lines, err := readDay18()
	if err != nil {
		return err
	}
	lc := []*snNumber{}
	for _, l := range lines {
		lc = append(lc, l.Copy(nil))
	}
	got := lc[0]
	for i := 1; i < len(lc); i++ {
		got = add18(got, lc[i])
	}
	partPrint(1, got.magnitude())
	var bm int
	for i := 0; i < len(lines); i++ {
		for j := 0; j < len(lines); j++ {
			if i == j {
				continue
			}
			r := add18(lines[i].Copy(nil), lines[j].Copy(nil))
			bm = maxint(bm, r.magnitude())
		}
	}
	partPrint(2, bm)
	return nil
}

func init() {
	RegisterDay(18, day18)
}
