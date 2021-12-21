package main

import (
	"fmt"
	"reflect"
	"testing"
)

var expectedResults = map[int][2]interface{}{
	1:  {1390, 1457},
	2:  {2322630, 2105273490},
	3:  {uint32(2595824), uint32(2135254)},
	4:  {35711, 5586},
	5:  {7438, 21406},
	6:  {uint64(345793), uint64(1572643095893)},
	7:  {352997, 101571302},
	8:  {519, 1027483},
	9:  {444, 1168440},
	10: {311895, 2904180541},
	11: {1647, 348},
	12: {5576, 152837},
	13: {847, `
###   ##  #### ###   ##  ####  ##  ### 
#  # #  #    # #  # #  # #    #  # #  #
###  #      #  #  # #    ###  #  # ### 
#  # #     #   ###  #    #    #### #  #
#  # #  # #    # #  #  # #    #  # #  #
###   ##  #### #  #  ##  #### #  # ### `},
	14: {2899, 3528317079545},
	15: {702, 2955},
	16: {895, int64(1148595959144)},
	17: {6555, 4973},
	18: {3654, 4578},
	19: {459, 19130},
	20: {5571, 17965},
	21: {925605, uint64(486638407378784)},
}

func collectResults(f func() error) ([2]interface{}, error) {
	oldPartPrint := partPrint
	defer func() {
		partPrint = oldPartPrint
	}()
	var got [2]interface{}
	partPrint = func(i int, x interface{}) {
		got[i-1] = x
	}
	err := f()
	return got, err
}

func TestDays(t *testing.T) {
	for i := 0; i < 25; i++ {
		if (validDays>>i)&1 != 1 {
			if expectedResults[i+1][0] != nil {
				t.Errorf("expected results %v for day %d, but it's not registered", expectedResults[i+1], i+1)
			}
			continue
		}
		t.Run(fmt.Sprintf("day %d", i+1), func(t *testing.T) {
			got, err := collectResults(days[i])
			if err != nil {
				t.Fatalf("error %v", err)
			}
			want := expectedResults[i+1]
			for i := 0; i < 2; i++ {
				if reflect.TypeOf(got[i]) != reflect.TypeOf(want[i]) {
					t.Errorf("mistyped got[%d] (%T) vs want[%d] (%T)", i, got[i], i, want[i])
				}
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("day%d() = %v, want %v", i+1, got, want)
			}
		})
	}
}
