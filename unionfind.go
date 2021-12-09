package main

type rankParent struct {
	rank, parent int
}

// UnionFind is a disjoint set datstructure.
type UnionFind struct {
	nodes []rankParent
}

// NewUnionFind returns a disjoint set datastructure
// with elements 0 to size-1.
func NewUnionFind(size int) *UnionFind {
	nodes := make([]rankParent, size)
	for i := 0; i < size; i++ {
		nodes[i] = rankParent{0, i}
	}
	return &UnionFind{
		nodes: nodes,
	}
}

// Find returns the identified element of the set that
// contains x. All elements in the same set as x return
// the same Find. The root may change after a Union operation.
func (uf *UnionFind) Find(x int) int {
	// "Path splitting" -- setting each node's parent to their current
	// grandfather, is a cache-friendly way of path compression.
	for uf.nodes[x].parent != x {
		x, uf.nodes[x].parent = uf.nodes[x].parent, uf.nodes[uf.nodes[x].parent].parent
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
	rankx := uf.nodes[rx].rank
	ranky := uf.nodes[ry].rank
	if rankx < ranky {
		rx, ry = ry, rx
	}
	// rx has the larger rank and becomes the root node of the
	// merged set.
	uf.nodes[ry].parent = rx
	if rankx == ranky {
		uf.nodes[rx].rank++
	}
}
