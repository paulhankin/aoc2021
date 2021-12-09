package main

import (
	"bufio"
	"os"
	"sort"
)

func readDay7() ([]int, error) {
	f, err := os.Open("day7.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var r []int
	for scanner.Scan() {
		xs, err := parseInts(scanner.Text(), ",")
		if err != nil {
			return nil, err
		}
		r = append(r, xs...)
	}
	return r, scanner.Err()
}

// sumn returns the sum i=0..x (inclusive)
func sumn(x int) int {
	return x * (x + 1) / 2
}

func day7() error {
	crabs, err := readDay7()
	if err != nil {
		return err
	}
	sort.Ints(crabs)
	// The mode m of a sequence of ints minimizes sum(|x[i] - m|)
	// and it doesn't matter if the length of the array is odd
	// and the true mode lies between two elements -- either of
	// them gives the same answer.
	mode := crabs[len(crabs)/2]
	// The mean m of a sequence of ints minimizes sum(|x[i] - m|^2).
	// I guess that the optimal position of m that minimizes
	// sum(|x[i]-m|^2 + |x[i]-m|)/2 is the mean rounded up or down.
	mean1 := intSum(crabs) / len(crabs)
	mean2 := intSum(crabs)/len(crabs) + 1
	var fuel1, fuel21, fuel22 int
	for _, c := range crabs {
		fuel1 += abs(c - mode)
		fuel21 += sumn(abs(c - mean1))
		fuel22 += sumn(abs(c - mean2))
	}
	partPrint(1, fuel1)
	if fuel21 < fuel22 {
		partPrint(2, fuel21)
	} else {
		partPrint(2, fuel22)
	}
	return nil

}

func init() {
	RegisterDay(7, day7)
}
