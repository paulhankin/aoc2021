package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	dayFlag = flag.Int("day", 0, "which day to run")
)

var (
	days      = initDays()
	validDays uint32
)

func initDays() []func() error {
	var r []func() error
	for i := 0; i < 25; i++ {
		day := i
		r = append(r, func() error { return fmt.Errorf("no code for day %d", day+1) })
	}
	return r
}

func RegisterDay(day int, f func() error) {
	days[day-1] = f
	validDays |= 1 << uint32(day-1)
}

func main() {
	flag.Parse()
	if *dayFlag == 0 {
		exit := 0
		for i := 0; i < 25; i++ {
			if (validDays>>i)&1 == 1 {
				fmt.Println("day", i+1)
				if err := days[i](); err != nil {
					fmt.Printf("***error: %v\n", err)
					exit = 1
				}
				fmt.Println()
			}
		}
		os.Exit(exit)
	}
	day := *dayFlag - 1
	if day < 0 || day >= len(days) {
		log.Fatalf("day %d out of range", *dayFlag)
	}
	if err := days[day](); err != nil {
		log.Fatalf("error from day %d: %v", *dayFlag, err)
	}
}
