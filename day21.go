package main

func day21x1(p1, p2 int) {
	s := [2]int{0, 0}
	rolls := 0
	p := [2]int{p1 - 1, p2 - 1}
	for s[0] < 1000 && s[1] < 1000 {
		d1 := (rolls*3)%100 + 1
		d2 := (rolls*3+1)%100 + 1
		d3 := (rolls*3+2)%100 + 1
		n := rolls % 2
		p[n] = (p[n] + d1 + d2 + d3) % 10
		s[n] += p[n] + 1
		rolls++
	}
	if s[0] >= 1000 {
		partPrint(1, s[1]*rolls*3)
	} else {
		partPrint(1, s[0]*rolls*3)
	}
}

type key struct {
	p1, p2 uint8
	s1, s2 uint8
}

func (k key) idx() int {
	return (int(k.p1)*10+int(k.p2))*22*22 + (22*int(k.s1) + int(k.s2))
}

var nrolls = []uint64{
	0, 0, 0, 1, 3, 6, 7, 6, 3, 1,
}

func universes21(k key, m *[10 * 10 * 22 * 22][2]uint64) [2]uint64 {
	v := (*m)[k.idx()]
	if v[0]+v[1] != 0 {
		return v
	}
	if k.s1 == 21 {
		v = [2]uint64{1, 0}
	} else if k.s2 == 21 {
		v = [2]uint64{0, 1}
	} else {
		for d := 3; d <= 9; d++ {
			s2 := k.s1 + (k.p1+uint8(d))%10 + 1
			if s2 > 21 {
				s2 = 21
			}
			k2 := key{
				p1: k.p2,
				p2: (k.p1 + uint8(d)) % 10,
				s1: k.s2,
				s2: s2,
			}
			v2 := universes21(k2, m)
			v = [2]uint64{v[0] + nrolls[d]*v2[1], v[1] + nrolls[d]*v2[0]}
		}
	}
	(*m)[k.idx()] = v
	return v
}

func day21x2(p1i, p2i int) {
	var m [10 * 10 * 22 * 22][2]uint64
	v := universes21(key{p1: uint8(p1i - 1), p2: uint8(p2i - 1)}, &m)
	if v[0] > v[1] {
		partPrint(2, v[0])
	} else {
		partPrint(2, v[1])
	}
}

func init() {
	RegisterDay(21, func() error {
		p1, p2 := 6, 9
		day21x1(p1, p2)
		day21x2(p1, p2)
		return nil
	})
}
