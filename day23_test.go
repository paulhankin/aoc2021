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

func BenchmarkDay23(b *testing.B) {
	cases := map[string][2]interface{}{
		"BCBDADCA": {12521, 44169},
		"BCADBCDA": {11120, 49232},
	}
	for n := 0; n < b.N; n++ {
		for s, want := range cases {
			got, err := collectResults(func() error {
				day23s(s)
				return nil
			})
			if err != nil {
				b.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				b.Errorf("got results %v want %v", got, want)
			}
		}
	}

}
