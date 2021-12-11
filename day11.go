package main

import (
	"bufio"
	"fmt"
	"os"
)

func readDay11() (*[10][10]int, error) {
	f, err := os.Open("day11.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var b [10][10]int
	var j int
	for scanner.Scan() {
		var line [10]int
		for i, x := range []byte(scanner.Text()) {
			if x < '0' || x > '9' {
				return nil, fmt.Errorf("unknown digit %c", x)
			}
			line[i] = int(x - '0')
		}
		b[j] = line
		j++
	}
	return &b, scanner.Err()
}

func day11Step(xs *[10][10]int) int {
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			(*xs)[i][j]++
		}
	}
	var flashed [10][10]bool
	for anyFlash := true; anyFlash; {
		anyFlash = false
		for i := 0; i < 10; i++ {
			for j := 0; j < 10; j++ {
				if (*xs)[i][j] == 10 && !flashed[i][j] {
					flashed[i][j] = true
					anyFlash = true
					for d := 0; d < 8; d++ {
						di, dj := dir8(d)
						if i+di < 0 || i+di >= 10 || j+dj < 0 || j+dj >= 10 {
							continue
						}
						if (*xs)[i+di][j+dj] < 10 {
							(*xs)[i+di][j+dj]++
						}
					}
				}
			}
		}
	}
	var nflash int
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if (*xs)[i][j] == 10 {
				(*xs)[i][j] = 0
				nflash++
			}
		}
	}

	return nflash
}

func day11() error {
	xs, err := readDay11()
	if err != nil {
		return err
	}
	var totalFlashes int
	var sync int
	for steps := 0; steps < 100 || sync == 0; steps++ {
		if false {
			for _, row := range *xs {
				for _, b := range row {
					fmt.Printf("%c", rune(b)+'0')
				}
				fmt.Println()
			}
			fmt.Println()
		}
		nflash := day11Step(xs)
		if nflash == 100 {
			if sync == 0 {
				sync = steps + 1
			}
		}
		if steps < 100 {
			totalFlashes += nflash
		}
	}
	partPrint(1, totalFlashes)
	partPrint(2, sync)
	return nil
}

func init() {
	RegisterDay(11, day11)
}
