package main

import (
	"reflect"
	"testing"
)

func TestExample22s(t *testing.T) {
	cs := map[string][2]interface{}{
		"day22_example2.txt": {590784, 39769202357779},
		"day22_example3.txt": {474140, 2758514936282235},
	}
	for c, want := range cs {
		t.Run(c, func(t *testing.T) {
			got, err := collectResults(func() error {
				return day22filename(c)
			})
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("got results %v want %v", got, want)
			}

		})
	}
}
