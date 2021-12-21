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
	n      int
	p1, p2 int
	s1, s2 int
}

var nrolls = []uint64{
	0, 0, 0, 1, 3, 6, 7, 6, 3, 1,
}

func (k *key) posAdd(n, v int) {
	if n == 0 {
		k.p1 = (k.p1 + v) % 10
	} else {
		k.p2 = (k.p2 + v) % 10
	}
}
func (k *key) posGet(n int) int {
	if n == 0 {
		return k.p1
	} else {
		return k.p2
	}
}

func (k *key) scoreAdd(n, v int) {
	if n == 0 {
		k.s1 = minint(21, k.s1+v)
	} else {
		k.s2 = minint(21, k.s2+v)
	}
}

func universes21(k key, m map[key][2]uint64) [2]uint64 {
	v, ok := m[k]
	if ok {
		return v
	}
	if k.s1 == 21 {
		v = [2]uint64{1, 0}
	} else if k.s2 == 21 {
		v = [2]uint64{0, 1}
	} else {
		for d := 3; d <= 9; d++ {
			k2 := k
			k2.n = 1 - k.n
			k2.posAdd(k.n, d)
			k2.scoreAdd(k.n, k2.posGet(k.n)+1)
			v2 := universes21(k2, m)
			v = [2]uint64{v[0] + nrolls[d]*v2[0], v[1] + nrolls[d]*v2[1]}
		}
	}
	m[k] = v
	return v
}

func day21x2(p1i, p2i int) {
	m := map[key][2]uint64{}
	v := universes21(key{n: 0, p1: p1i - 1, p2: p2i - 1}, m)
	if v[0] > v[1] {
		partPrint(2, v[0])
	} else {
		partPrint(2, v[1])
	}
}

func init() {
	RegisterDay(21, func() error {
		p1, p2 := 6, 9 // 4, 8
		day21x1(p1, p2)
		day21x2(p1, p2)
		return nil
	})
}
