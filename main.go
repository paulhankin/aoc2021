package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	dayFlag = flag.Int("day", 1, "which day to run")
)

var days = initDays()

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
}

func main() {
	flag.Parse()
	day := *dayFlag - 1
	if day < 0 || day >= len(days) {
		log.Fatalf("day %d out of range", *dayFlag)
	}
	if err := days[day](); err != nil {
		log.Fatalf("error from day %d: %v", *dayFlag, err)
	}
}
