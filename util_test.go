package main

import "testing"

func TestDir8(t *testing.T) {
	var got uint32
	var want uint32
	for di := -1; di <= 1; di++ {
		for dj := -1; dj <= 1; dj++ {
			if di == 0 && dj == 0 {
				continue
			}
			want |= 1 << ((di+1)*3 + dj + 1)
		}
	}
	for i := 0; i < 8; i++ {
		di, dj := dir8(i)
		if di != clamp(di, -1, 1) || dj != clamp(dj, -1, 1) {
			t.Errorf("dir8(%d)=%d,%d out of range", i, di, dj)
			continue
		}
		bit := uint32(1) << ((di+1)*3 + dj + 1)
		if got&bit != 0 {
			t.Errorf("dir8(%d)=%d,%d equal to previous result", i, di, dj)
		}
		got |= bit
	}
	if got != want {
		t.Errorf("bitset of directions from dir8 = %d, want %d", got, want)
	}
}
