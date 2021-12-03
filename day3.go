package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func readDay3() ([]uint32, error) {
	f, err := os.Open("day3.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var r []uint32
	for scanner.Scan() {
		x, err := strconv.ParseUint(scanner.Text(), 2, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q: %v", scanner.Text(), err)
		}
		r = append(r, uint32(x))
	}
	return r, nil
}

func filterBits(in []uint32, b bool) uint32 {
	for j := 11; j >= 0; j-- {
		if len(in) <= 1 {
			return in[0]
		}
		var got [2][]uint32
		for _, i := range in {
			got[(i>>j)&1] = append(got[(i>>j)&1], i)
		}
		if b {
			if len(got[1]) >= len(got[0]) {
				in = got[1]
			} else {
				in = got[0]
			}
		} else {
			if len(got[0]) > 0 && len(got[0]) <= len(got[1]) || len(got[1]) == 0 {
				in = got[0]
			} else {
				in = got[1]
			}
		}
	}
	if len(in) <= 1 {
		return in[0]
	}
	log.Fatal("failed to reduce to single number")
	return 0
}

func day3() error {
	in, err := readDay3()
	if err != nil {
		return err
	}
	bits := [12]int{}
	for _, i := range in {
		for j := 0; j < 12; j++ {
			if (i>>j)&1 == 1 {
				bits[j]++
			}
		}
	}
	var gamma, epsilon uint32
	for j := 0; j < 12; j++ {
		if bits[j] > len(in)-bits[j] {
			gamma |= uint32(1) << j
		} else {
			epsilon |= uint32(1) << j
		}
	}
	fmt.Println("part 1 =", gamma*epsilon)

	ogr := filterBits(in, true)
	co2 := filterBits(in, false)

	fmt.Println("part 2 =", ogr*co2)
	return nil
}

func init() {
	RegisterDay(3, day3)
}
