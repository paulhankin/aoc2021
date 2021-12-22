package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

func readDay19(filename string) ([][]coord3i, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var r [][]coord3i
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		} else if strings.Contains(line, "--- scanner") {
			wantScan := fmt.Sprintf("--- scanner %d ---", len(r))
			if line != wantScan {
				return nil, fmt.Errorf("expected %q, found %q", wantScan, line)
			}
			r = append(r, nil)
		} else {
			c, err := parseCoord3i(line)
			if err != nil {
				return nil, err
			}
			r[len(r)-1] = append(r[len(r)-1], c)
		}
	}
	return r, nil
}

// axrot3d is a number from 0 to 23 that represents
// a axis-aligned rotation by multiples of 90 degrees in 3d.
type axrot3d int

// axform3d is a rotation by r, then a translation by t.
type axform3d struct {
	t coord3i
	r axrot3d
}

func (a axrot3d) rotate(c coord3i) coord3i {
	ax := perm3[a%6]
	p := a % 6 % 2
	sign := a / 6 // 0, 1, 2, 3
	return coord3i{c.get(ax.x) * sign3[p][sign].x, c.get(ax.y) * sign3[p][sign].y, c.get(ax.z) * sign3[p][sign].z}
}

var axrot3dinvs [24]axrot3d
var axrot3dcomps [24][24]axrot3d

func (a axrot3d) inverse() axrot3d {
	return axrot3dinvs[a]
}

func (x axrot3d) String() string {
	return fmt.Sprintf("rot(%s)", x.rotate(coord3i{1, 2, 3}))
}

func (x axrot3d) compose(y axrot3d) axrot3d {
	return axrot3dcomps[x][y]
}

func init() {
	for i := axrot3d(0); i < 24; i++ {
		for j := axrot3d(0); j < 24; j++ {
			got := j.rotate(i.rotate(coord3i{1, 2, 3}))
			if got.x == 1 && got.y == 2 && got.z == 3 {
				axrot3dinvs[i] = j
			}
			for k := axrot3d(0); k < 24; k++ {
				if got == k.rotate(coord3i{1, 2, 3}) {
					axrot3dcomps[i][j] = k
				}
			}
		}
	}
}

func (x axform3d) inverse() axform3d {
	ri := x.r.inverse()
	return axform3d{
		t: ri.rotate(x.t).mul(-1),
		r: ri,
	}
}

// compose is the result of doing xform y after x.
func (x axform3d) compose(y axform3d) axform3d {
	return axform3d{
		r: x.r.compose(y.r),
		t: y.r.rotate(x.t).add(y.t),
	}
}

func (x axform3d) String() string {
	return fmt.Sprintf("{r:%s, t:%s}", x.r, x.t)
}

func (x axform3d) apply(c coord3i) coord3i {
	return x.r.rotate(c).add(x.t)
}

// perm3 is all permutations of {0, 1, 2} such that the even permutations
// are at even indexes.
var perm3 = []coord3i{
	{0, 1, 2},
	{0, 2, 1},
	{1, 2, 0},
	{1, 0, 2},
	{2, 0, 1},
	{2, 1, 0},
}

var sign3 = [2][]coord3i{
	// Even numbers of change of sign
	{
		{1, 1, 1},
		{1, -1, -1},
		{-1, -1, 1},
		{-1, 1, -1},
	},
	// Odd numbers of change of sign
	{
		{-1, 1, 1},
		{1, -1, 1},
		{1, 1, -1},
		{-1, -1, -1},
	},
}

// align19 finds an offset between s1 and s2 such that at least
// 12 beacons in common.
// Positions in s2 are rotated by r!
func align19(s1, s2 []coord3i, r axrot3d) (coord3i, bool) {
	diffs := map[coord3i]int{}
	maxdiff := 0
	for i2, c2 := range s2 {
		if maxdiff+len(s2)-i2 < 12 {
			return coord3i{}, false
		}
		c2r := r.rotate(c2)
		for _, c1 := range s1 {
			d := c1.sub(c2r)
			diffs[d] += 1
			maxdiff = maxint(maxdiff, diffs[d])
			if maxdiff >= 12 {
				return d, true
			}
		}
	}
	return coord3i{}, false
}

func alignrot19(s1, s2 []coord3i) (axform3d, bool) {
	for r := axrot3d(0); r < 24; r++ {
		if c, ok := align19(s1, s2, r); ok {
			return axform3d{r: r, t: c}, true
		}
	}
	return axform3d{}, false
}

type sensalign struct {
	sensor int
	xf     axform3d
}

func day19withFilename(s string) error {
	ss, err := readDay19(s)
	if err != nil {
		return err
	}
	done := map[int]bool{0: true}
	todo := []int{0}
	aligns := make([]axform3d, len(ss))
	aligns[0] = axform3d{}
	for len(todo) > 0 {
		i := todo[len(todo)-1]
		todo = todo[:len(todo)-1]
		var wg sync.WaitGroup
		var mut sync.Mutex
		for j := 0; j < len(ss); j++ {
			if done[j] {
				continue
			}
			wg.Add(1)
			go func(j int) {
				xf, ok := alignrot19(ss[i], ss[j])
				if ok {
					mut.Lock()
					defer mut.Unlock()
					done[j] = true
					aligns[j] = xf.compose(aligns[i])
					todo = append(todo, j)
				}
				wg.Done()
			}(j)
		}
		wg.Wait()
	}
	beacons := map[coord3i]bool{}
	for i, s := range ss {
		xf := aligns[i]
		for _, c := range s {
			beacons[xf.apply(c)] = true
		}
	}
	partPrint(1, len(beacons))
	maxMD := 0
	for i := 0; i < len(ss); i++ {
		for j := i + 1; j < len(ss); j++ {
			maxMD = maxint(maxMD, aligns[j].t.sub(aligns[i].t).l1norm())
		}
	}
	partPrint(2, maxMD)
	return nil
}

func init() {
	RegisterDay(19, func() error {
		return day19withFilename("day19.txt")
	})
}
