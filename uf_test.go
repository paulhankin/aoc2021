package main

import (
	"testing"
)

func TestUnionFind1(t *testing.T) {
	const N = 999983
	uf := NewUnionFind(N)
	for i := 0; i < N; i++ {
		i0 := ((i + 1234) * 457) % N
		i1 := ((i+1234)*457 + 1) % N
		uf.Union(i0, i1)
	}
	S := map[int]int{}
	for i := 0; i < N; i++ {
		S[uf.Find(i)]++
	}
	if len(S) > 2 {
		t.Errorf("too many sets (%d)", len(S))
	}
	sum := 0
	for _, i := range S {
		sum += i
	}
	if sum != N {
		t.Fatalf("sum of sets is %d, want %d", sum, N)
	}
}

func TestUnionFind2(t *testing.T) {
	const N = 1000000
	uf := NewUnionFind(N)
	for i := 0; i < N; i++ {
		uf.Union(i, (i+123)%N)
	}
	x0 := uf.Find(0)
	for i := 0; i < N; i++ {
		if uf.Find(i) != x0 {
			t.Fatalf("expected only one set, found %d disjoint from 0", i)
		}
	}
}
