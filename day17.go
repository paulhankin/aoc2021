package main

type day17range struct {
	xMin, xMax, yMin, yMax int
}

var (
	day17input = day17range{207, 263, -115, -63}
)

// The possible number of steps that can result
// in hitting the x-range of the target to how many
// initial dx produce that value.
func solve17x(r day17range) map[int][]int {
	steps := map[int][]int{}
	for dxx := 1; dxx <= r.xMax; dxx++ {
		x := 0
		dx := dxx
		nStep := 0
		stopSteps := 1000
		for stopSteps > 0 && x <= r.xMax {
			if x >= r.xMin {
				steps[nStep] = append(steps[nStep], dxx)
			}
			nStep++
			x += dx
			dx -= sgn(dx)
			if dx == 0 {
				stopSteps--
			}
		}
	}
	return steps
}

type y17 struct {
	maxY int
	ys   []int
}

// Map from number of steps to the maximum y coord achieved and number of possible initial dy
func solve17y(r day17range) map[int]y17 {
	steps := map[int]y17{}
	for dyy := -1000; dyy < 5000; dyy++ {
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
	viableX := solve17x(r)
	viableY := solve17y(r)
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
