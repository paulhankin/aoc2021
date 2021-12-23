package main

import (
	"reflect"
	"testing"
)

func TestExample20(t *testing.T) {
	got, err := collectResults(func() error {
		return day20Filename("day20_example.txt")
	})
	if err != nil {
		t.Fatal(err)
	}
	want := [2]interface{}{35, 3351}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got results %v for day20 example, want %v", got, want)
	}
}
