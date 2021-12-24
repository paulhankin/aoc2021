package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type opcode uint8

const (
	opnop = opcode(0)
	opinp = opcode(1)
	opadd = opcode(2)
	opmul = opcode(3)
	opdiv = opcode(4)
	opmod = opcode(5)
	opeql = opcode(6)
)

type opval int8

const (
	regw = opval(-125)
	regx = opval(-126)
	regy = opval(-127)
	regz = opval(-128)
)

func (v opval) decode() (int, bool) {
	if v <= -125 {
		return int(-125 - v), true
	}
	return int(v), false
}

type op struct {
	c  opcode
	v1 opval
	v2 opval
}

func (c opcode) String() string {
	switch c {
	case opinp:
		return "inp"
	case opadd:
		return "add"
	case opmul:
		return "mul"
	case opdiv:
		return "div"
	case opmod:
		return "mod"
	case opeql:
		return "eql"
	default:
		return fmt.Sprintf("<%x>", int(c))
	}
}

func (c opval) String() string {
	i, reg := c.decode()
	if reg {
		return fmt.Sprintf("%c", 'w'+byte(i))
	}
	return fmt.Sprintf("%d", i)
}

func (c opcode) vals() int {
	switch c {
	case opinp:
		return 1
	case opadd, opmul, opdiv, opmod, opeql:
		return 2
	default:
		return 0
	}
}

func (op op) String() string {
	switch op.c.vals() {
	case 2:
		return fmt.Sprintf("%s %s %s", op.c, op.v1, op.v2)
	case 1:
		return fmt.Sprintf("%s %s", op.c, op.v1)
	default:
		return fmt.Sprintf("%s", op.c)
	}
}

func parseReg(s string) (opval, error) {
	if len(s) != 1 || s[0] < 'w' || s[0] > 'z' {
		return 0, fmt.Errorf("bad register %q", s)
	}
	return regw - opval(s[0]-'w'), nil
}

func parseVal(s string) (opval, error) {
	if s != "" && s[0] >= 'a' && s[0] <= 'z' {
		return parseReg(s)
	}
	x, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return 0, err
	}
	if x > 127 || x <= int64(regw) {
		return 0, fmt.Errorf("int argument out of range!")
	}
	return opval(x), nil
}

func parseOp(c opcode, p []string, argp ...func(s string) (opval, error)) (op, error) {
	if len(p) != len(argp) {
		return op{}, fmt.Errorf("got %d args, want %d", len(p), len(argp))
	}
	o := op{c: c}
	var err error
	o.v1, err = argp[0](p[0])
	if err != nil {
		return op{}, err
	}
	if len(p) > 1 {
		o.v2, err = argp[1](p[1])
		if err != nil {
			return op{}, err
		}
	}
	return o, nil
}

func readDay24(name string) ([]op, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var ops []op
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		var o op
		var err error
		switch parts[0] {
		case "inp":
			o, err = parseOp(opinp, parts[1:], parseReg)
		case "add":
			o, err = parseOp(opadd, parts[1:], parseReg, parseVal)
		case "mul":
			o, err = parseOp(opmul, parts[1:], parseReg, parseVal)
		case "div":
			o, err = parseOp(opdiv, parts[1:], parseReg, parseVal)
		case "mod":
			o, err = parseOp(opmod, parts[1:], parseReg, parseVal)
		case "eql":
			o, err = parseOp(opeql, parts[1:], parseReg, parseVal)
		default:
			return nil, fmt.Errorf("unknown opcode %q", parts[0])
		}
		if err != nil {
			return nil, fmt.Errorf("failure when parsing opcode %q: %v", parts[0], err)
		}
		ops = append(ops, o)
	}
	return ops, scanner.Err()
}

type opblock []op

func (opb opblock) String() string {
	var parts []string
	for _, op := range opb {
		parts = append(parts, op.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(parts, "; "))
}

type con24 struct {
	c, d, n int
}

func extractConstants(b opblock) con24 {
	d, dr := b[4].v2.decode()
	c, cr := b[5].v2.decode()
	n, nr := b[15].v2.decode()
	if dr || cr || nr {
		panic("constant extract failed: some ops were registers")
	}
	return con24{c: int(c), d: int(d), n: int(n)}
}

// eval24 evaluates the block that the constants have been extracted
// from, with the given input value of z.
// z0 is the new z if the w inputed is the given w.
// z1 is the new z if the w inputed is not given w.
func eval24(c con24, z int) []int {
	xcond := z%26 + c.c
	zs := make([]int, 9)
	for w := 1; w <= 9; w++ {
		if xcond == w {
			zs[w-1] = z / c.d
		} else {
			zs[w-1] = (z/c.d)*26 + w + c.n
		}
	}
	return zs
}

// solve24x returns the largest string of w digits
// that, starting with the given z results in z=0 by the end.
func solve24x(cs []con24, z int, cache map[[2]int]int, smallest bool) int {
	key := [2]int{len(cs), z}
	if r, ok := cache[key]; ok {
		return r
	}
	if len(cs) == 0 {
		if z == 0 {
			return 0
		} else {
			return -1
		}
	}
	zs := eval24(cs[0], z)
	for ww := 0; ww < 9; ww++ {
		w := 9 - ww
		if smallest {
			w = 1 + ww
		}
		b := solve24x(cs[1:], zs[w-1], cache, smallest)
		if b != -1 {
			r := b*10 + w
			cache[key] = r
			return r
		}
	}
	cache[key] = -1
	return -1
}

func revint(x int, n int) int {
	var r int
	for i := 0; i < n; i++ {
		r = r*10 + (x % 10)
		x /= 10
	}
	return r
}

func solve24(cs []con24, smallest bool) int {
	return revint(solve24x(cs, 0, map[[2]int]int{}, smallest), 14)
}

func day24s(name string) error {
	ops, err := readDay24(name)
	if err != nil {
		return err
	}
	var blocks []opblock
	for _, op := range ops {
		if op.c == opinp {
			blocks = append(blocks, nil)
		}
		blocks[len(blocks)-1] = append(blocks[len(blocks)-1], op)
	}
	var cs []con24
	for _, b := range blocks {
		con := extractConstants(b)
		cs = append(cs, con)
	}
	for part := 1; part <= 2; part++ {
		partPrint(part, solve24(cs, part == 2))
	}

	return nil
}

func init() {
	RegisterDay(24, func() error {
		return day24s("day24.txt")
	})
}
