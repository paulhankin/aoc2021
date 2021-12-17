package main

type day17range struct {
	xMin, xMax, yMin, yMax int
}

var (
	day17input = day17range{207, 263, -115, -63}
)

// The possible number of steps (up to maxSteps) that can
// result in hitting the x-range of the target to which
// initial dx produce that value.
func solve17x(r day17range, maxSteps int) map[int][]int {
	steps := map[int][]int{}
	for dxx := 1; dxx <= r.xMax; dxx++ {
		x := 0
		dx := dxx
		nStep := 0
		for nStep <= maxSteps && x <= r.xMax {
			if x >= r.xMin {
				steps[nStep] = append(steps[nStep], dxx)
			}
			nStep++
			x += dx
			dx -= sgn(dx)
		}
	}
	return steps
}

type y17 struct {
	maxY int
	ys   []int
}

// Map from number of steps to the maximum y coord achieved and the possible initial dy
func solve17y(r day17range) map[int]y17 {
	steps := map[int]y17{}
	// If we choose an initial positive velocity V, then at some later point we'll
	// return to y=0 with velocity -V, and on the next step be below V.
	// So we only need to consider initial velocities from yMin to -yMin!
	for dyy := r.yMin; dyy <= -r.yMin; dyy++ {
		y := 0
		maxY := 0
		dy := dyy
		nStep := 0
		for y >= r.yMin {
			if y <= r.yMax {
				steps[nStep] = y17{maxint(steps[nStep].maxY, maxY), append(steps[nStep].ys, dyy)}
			}
			nStep++
			y += dy
			maxY = maxint(maxY, y)
			dy -= 1
		}
	}
	return steps
}

func day17() error {
	r := day17input
	viableY := solve17y(r)
	maxSteps := 0
	for s := range viableY {
		maxSteps = maxint(maxSteps, s)
	}
	viableX := solve17x(r, maxSteps)
	maxY := 0
	got := map[[2]int]bool{}
	for s, xs := range viableX {
		maxY = maxint(viableY[s].maxY, maxY)
		for _, x := range xs {
			for _, y := range viableY[s].ys {
				got[[2]int{x, y}] = true
			}
		}
	}
	partPrint(1, maxY)
	partPrint(2, len(got))
	return nil
}

func init() {
	RegisterDay(17, day17)
}
