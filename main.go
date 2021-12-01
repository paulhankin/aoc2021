package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	dayFlag = flag.Int("day", 1, "which day to run")
)

var days []func() error

func RegisterDay(day int, f func() error) {
	day--
	for day >= len(days) {
		i := len(days)
		days = append(days, func() error { return fmt.Errorf("no code for day %d", i) })
	}
	days[day] = f
}

func main() {
	flag.Parse()
	day := *dayFlag - 1
	if day < 0 || day >= len(days) {
		log.Fatalf("day %d out of range", *dayFlag)
	}
	days[day]()
}
