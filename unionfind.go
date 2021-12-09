package main

import "log"

// UnionFind is a disjoint set datstructure.
type UnionFind struct {
	parent []int
	// rank is at most log_2(n) so uint8 is sufficient
	// for 2^256 elements.
	rank []uint8
}

// NewUnionFind returns a disjoint set datastructure
// with elements 0 to size-1.
func NewUnionFind(size int) *UnionFind {
	parent := make([]int, size)
	for i := 0; i < size; i++ {
		parent[i] = i
	}
	return &UnionFind{
		parent: parent,
		rank:   make([]uint8, size),
	}
}

// Find returns the identified element of the set that
// contains x. All elements in the same set as x return
// the same Find. The root may change after a Union operation.
func (uf *UnionFind) Find(x int) int {
	// "Path splitting" -- setting each node's parent to their current
	// grandfather, is a cache-friendly way of path compression.
	for uf.parent[x] != x {
		x, uf.parent[x] = uf.parent[x], uf.parent[uf.parent[x]]
	}
	return x
}

// Union merges the two sets which contain x and y.
func (uf *UnionFind) Union(x, y int) {
	rx := uf.Find(x)
	ry := uf.Find(y)
	if rx == ry {
		return
	}
	rankx := uf.rank[rx]
	ranky := uf.rank[ry]
	if rankx < ranky {
		rx, ry = ry, rx
	}
	// rx has the larger rank and becomes the root node of the
	// merged set.
	uf.parent[ry] = rx
	if rankx == ranky {
		if uf.rank[rx] == 255 {
			log.Fatal("rank 255 should not be possible without a huge number of elements")
		}
		uf.rank[rx]++
	}
}
