package main

import (
	"reflect"
	"testing"
)

func TestExample23(t *testing.T) {
	got, err := collectResults(func() error {
		s0example := "BCBDADCA"
		day23s(s0example)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	want := [2]interface{}{12521, 44169}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got results %v want %v", got, want)
	}
}
