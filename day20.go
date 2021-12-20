package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func readDay20(name string) ([]bool, pic20, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, pic20{}, err
	}
	defer f.Close()
	b := bufio.NewReader(f)
	comp := make([]bool, 0, 512)
	for {
		got, err := b.ReadByte()
		if err != nil {
			return nil, pic20{}, err
		}
		if got == '#' {
			comp = append(comp, true)
		} else if got == '.' {
			comp = append(comp, false)
		}
		if len(comp) == 512 {
			break
		}
	}
	var pic [][]bool
	var line []bool
	for {
		got, err := b.ReadByte()
		if got == '\n' || err == io.EOF {
			if len(line) > 0 {
				pic = append(pic, line)
				line = nil
			}
		}
		if err != nil {
			if err == io.EOF {
				return comp, pic20{b: pic}, nil
			}
			return nil, pic20{}, err
		}
		if got == '#' {
			line = append(line, true)
		} else if got == '.' {
			line = append(line, false)
		}
	}
}

type pic20 struct {
	b   [][]bool
	out bool
}

func (p pic20) String() string {
	m := p.Height()
	n := p.Width()
	var s strings.Builder
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if p.Get(i, j) {
				s.WriteByte('#')
			} else {
				s.WriteByte('.')
			}
		}
		s.WriteByte('\n')
	}
	return s.String()
}

func (p pic20) Height() int {
	return len(p.b)
}
func (p pic20) Width() int {
	return len(p.b[0])
}
func (p pic20) Get(i, j int) bool {
	if i < 0 || j < 0 || i >= p.Height() || j >= p.Width() {
		return p.out
	}
	return p.b[i][j]
}

func newPic20(m, n int) pic20 {
	r := make([][]bool, m)
	for i := range r {
		r[i] = make([]bool, n)
	}
	return pic20{b: r}
}

func (p pic20) decompress(c []bool, i, j int) bool {
	b := uint32(0)
	for k := 0; k < 9; k++ {
		b |= uint32(b2i(p.Get(i-1+(k/3), j-1+(k%3)))) << (8 - k)
	}
	return c[b]
}

func (p pic20) Step(c []bool) pic20 {
	q := newPic20(p.Height()+2, p.Width()+2)
	for i := 0; i < q.Height(); i++ {
		for j := 0; j < q.Width(); j++ {
			q.b[i][j] = p.decompress(c, i-1, j-1)
		}
	}
	if !p.out {
		q.out = c[0]
	} else {
		q.out = c[511]
	}
	return q
}

func (p pic20) Count() int {
	var i int
	for _, r := range p.b {
		for _, b := range r {
			if b {
				i++
			}
		}
	}
	return i
}

func day20Filename(name string) error {
	comp, pic0, err := readDay20(name)
	if err != nil {
		return err
	}
	for part := 1; part <= 2; part++ {
		pic := pic0
		steps := []int{2, 50}[part-1]
		for s := 0; s < steps; s++ {
			pic = pic.Step(comp)
		}
		partPrint(part, pic.Count())
	}
	return nil
}

func init() {
	RegisterDay(20, func() error { return day20Filename("day20.txt") })
}
