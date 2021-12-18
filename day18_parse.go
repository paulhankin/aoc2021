package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"text/scanner"
)

// snNumber is either a number n, or a pair of snNumber.
// n is valid if and only if p is nil.
type snNumber struct {
	n      int
	p      *[2]*snNumber
	parent *snNumber
}

func (s snNumber) Pair() bool {
	return s.p != nil
}

func (s *snNumber) Copy(parent *snNumber) *snNumber {
	r := &snNumber{parent: parent, n: s.n}
	if s.Pair() {
		r.p = &[2]*snNumber{}
		r.p[0] = s.p[0].Copy(r)
		r.p[1] = s.p[1].Copy(r)
	}
	return r
}

const redOn = "\x1b[1;31m"
const reset = "\x1b[0m"

func (s snNumber) StringDepth(n int, showDepth bool) (r string) {
	if s.Pair() && n > 4 || !s.Pair() && s.n >= 10 {
		defer func() {
			r = redOn + r + reset
		}()
	}

	if !s.Pair() {
		return fmt.Sprintf("%d", s.n)
	}
	if showDepth {
		return fmt.Sprintf("[%d:%s,%s]", n, s.p[0].StringDepth(n+1, showDepth), s.p[1].StringDepth(n+1, showDepth))
	} else {
		return fmt.Sprintf("[%s,%s]", s.p[0].StringDepth(n+1, showDepth), s.p[1].StringDepth(n+1, showDepth))
	}
}

func (s snNumber) String() string {
	return s.StringDepth(0, false)
}

type eofError struct{}

func readChar(r io.Reader) byte {
	b := []byte{0}
	n, err := r.Read(b)
	if n == 0 {
		panic(eofError{})
	}
	if err != nil {
		panic(err)
	}
	return b[0]
}

func parse18literal(s *scanner.Scanner) (*snNumber, error) {
	tok := s.Scan()
	if tok != scanner.Int {
		return nil, fmt.Errorf("expected literal int, but found %s", s.TokenText())
	}
	a, err := strconv.ParseInt(s.TokenText(), 10, 64)
	if err != nil {
		return nil, err
	}
	return &snNumber{n: int(a)}, nil
}

func parse18number(s *scanner.Scanner) (*snNumber, error) {
	tok := s.Peek()
	if tok == '[' {
		return parse18pair(s)
	} else if tok >= '0' && tok <= '9' {
		return parse18literal(s)
	}
	return nil, fmt.Errorf("expected pair or integer literal, found %c", tok)
}

func parse18rune(s *scanner.Scanner, r rune) error {
	tok := s.Scan()
	if tok != r {
		return fmt.Errorf("expected %c but found %q", r, s.TokenText())
	}
	return nil
}

func parse18pair(s *scanner.Scanner) (*snNumber, error) {
	if err := parse18rune(s, '['); err != nil {
		return nil, err
	}
	n1, err := parse18number(s)
	if err != nil {
		return nil, err
	}
	if err := parse18rune(s, ','); err != nil {
		return nil, err
	}
	n2, err := parse18number(s)
	if err != nil {
		return nil, err
	}
	if err := parse18rune(s, ']'); err != nil {
		return nil, err
	}
	r := &snNumber{p: &[2]*snNumber{n1, n2}}
	n1.parent = r
	n2.parent = r
	return r, nil
}

func parse18lines(s *scanner.Scanner) ([]*snNumber, error) {
	var r []*snNumber
	for {
		tok := s.Peek()
		if tok == scanner.EOF {
			return r, nil
		}
		if tok == '\n' || tok == '\r' {
			_ = s.Next()
			continue
		} else if tok == '[' || tok >= '0' && tok <= '9' {
			n, err := parse18number(s)
			if err != nil {
				return nil, err
			}
			r = append(r, n)
		} else {
			return nil, fmt.Errorf("unexpected token %c when parsing lines", tok)
		}
	}
}

func parse18FromReader(f io.Reader) ([]*snNumber, error) {
	var s scanner.Scanner
	s.Init(f)
	s.Mode = scanner.ScanInts
	r, err := parse18lines(&s)
	if err != nil {
		return nil, err
	}
	if s.ErrorCount > 0 {
		return nil, fmt.Errorf("parse errors")
	}
	return r, nil
}

func readDay18() ([]*snNumber, error) {
	f, err := os.Open("day18.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parse18FromReader(f)
}
