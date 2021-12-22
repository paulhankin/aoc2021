package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type cuboid struct {
	set      bool
	min, max coord3i
}

var cubRE = regexp.MustCompile("^(on|off) x=(-?[0-9]+)..(-?[0-9]+),y=(-?[0-9]+)..(-?[0-9]+),z=(-?[0-9]+)..(-?[0-9]+)$")

func parseInt(s string) int {
	r, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return int(r)
}

func parseCuboid(s string) (cuboid, error) {
	parts := cubRE.FindStringSubmatch(s)
	if parts == nil {
		return cuboid{}, fmt.Errorf("failed to match %q", s)
	}
	xmin := parseInt(parts[2])
	xmax := parseInt(parts[3])
	ymin := parseInt(parts[4])
	ymax := parseInt(parts[5])
	zmin := parseInt(parts[6])
	zmax := parseInt(parts[7])
	if xmin > xmax || ymin > ymax || zmin > zmax {
		return cuboid{}, fmt.Errorf("min values not smaller than max values in %q", s)
	}
	return cuboid{
		set: parts[1] == "on",
		min: coord3i{xmin, ymin, zmin},
		max: coord3i{xmax, ymax, zmax},
	}, nil
}

func readDay22(name string) ([]cuboid, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var cs []cuboid
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		cub, err := parseCuboid(line)
		if err != nil {
			return nil, err
		}
		cs = append(cs, cub)
	}
	return cs, scanner.Err()

}

func cuboidIntersect(a, b cuboid) (r cuboid, ok bool, all bool) {
	if b.min.x > a.max.x || b.min.y > a.max.y || b.min.z > a.max.z {
		return cuboid{}, false, false
	}
	if b.max.x < a.min.x || b.max.y < a.min.y || b.max.z < a.min.z {
		return cuboid{}, false, false
	}
	xmin := clamp(b.min.x, a.min.x, a.max.x)
	xmax := clamp(b.max.x, a.min.x, a.max.x)
	ymin := clamp(b.min.y, a.min.y, a.max.y)
	ymax := clamp(b.max.y, a.min.y, a.max.y)
	zmin := clamp(b.min.z, a.min.z, a.max.z)
	zmax := clamp(b.max.z, a.min.z, a.max.z)
	r = cuboid{
		set: b.set,
		min: coord3i{xmin, ymin, zmin},
		max: coord3i{xmax, ymax, zmax},
	}
	return r, true, r.min == a.min && r.max == a.max
}

func (c cuboid) dx() int {
	return c.max.x - c.min.x + 1
}
func (c cuboid) dy() int {
	return c.max.y - c.min.y + 1
}
func (c cuboid) dz() int {
	return c.max.z - c.min.z + 1
}

func (c cuboid) size() int {
	return c.dx() * c.dy() * c.dz()
}

func filterCuboids(R cuboid, cs []cuboid) []cuboid {
	var rcs []cuboid
	for _, c := range cs {
		in, ok, all := cuboidIntersect(R, c)
		if !ok {
			continue
		}
		if all {
			rcs = nil
		}
		rcs = append(rcs, in)
	}
	return rcs
}

func (c cuboid) divide(axis, at int) (cuboid, cuboid) {
	left, right := c, c
	if axis == 0 {
		if c.dx() < 2 {
			panic("dx")
		}
		left.max.x = at
		right.min.x = at + 1
	} else if axis == 1 {
		if c.dy() < 2 {
			panic("dy")
		}
		left.max.y = at
		right.min.y = at + 1
	} else {
		if c.dz() < 2 {
			panic("dz")
		}
		left.max.z = at
		right.min.z = at + 1
	}
	return left, right
}

func countCuboids(R cuboid, cs []cuboid) int {
	cs = filterCuboids(R, cs)
	if len(cs) == 0 {
		return 0
	}
	if len(cs) == 1 {
		c := cs[0]
		if !c.set {
			return 0
		}
		return c.size()
	}
	if len(cs) == 2 {
		ci, ok, _ := cuboidIntersect(cs[0], cs[1])
		is := 0
		if ok {
			is = ci.size()
		}
		r := 0
		if cs[0].set {
			r += cs[0].size() - is
		}
		if cs[1].set {
			r += cs[1].size()
		}
		return r
	}
	// Find a place to cut the region such that it's divided
	// into two regions, and that the cut is on the edge of
	// at least one cuboid. That guarantees that we can't do
	// too many useless cuts which don't reduce the number of
	// cuboids we're considering: either the region shrinks to
	// the bounds of the cuboids in one axis, or at least one
	// cuboid is entirely on one side of the cut.
	axis := 0
	at := 0
	for _, c := range cs {
		if c.min.x > R.min.x {
			axis = 0
			at = c.min.x - 1
			break
		} else if c.max.x < R.max.x {
			axis = 0
			at = c.max.x
			break
		}
		if c.min.y > R.min.y {
			axis = 1
			at = c.min.y - 1
			break
		} else if c.max.y < R.max.y {
			axis = 1
			at = c.max.y
			break
		}
		if c.min.z > R.min.z {
			axis = 2
			at = c.min.z - 1
			break
		} else if c.max.z < R.max.z {
			axis = 2
			at = c.max.z
			break
		}
	}
	left, right := R.divide(axis, at)
	return countCuboids(left, cs) + countCuboids(right, cs)
}

func day22filename(name string) error {
	cs, err := readDay22(name)
	if err != nil {
		return err
	}
	const M = 130000
	Rs := []cuboid{
		{min: coord3i{-50, -50, -50}, max: coord3i{50, 50, 50}},
		{min: coord3i{-M, -M, -M}, max: coord3i{M, M, M}},
	}
	for pm1, rs := range Rs {
		partPrint(pm1+1, countCuboids(rs, cs))
	}
	return nil
}

func init() {
	RegisterDay(22, func() error {
		return day22filename("day22.txt")
	})
}
