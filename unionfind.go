package main

type rankParent struct {
	rank, parent int
}

type UnionFind struct {
	Nodes []rankParent
}

func NewUnionFind(size int) *UnionFind {
	nodes := make([]rankParent, size)
	for i := 0; i < size; i++ {
		nodes[i] = rankParent{0, i}
	}
	return &UnionFind{
		Nodes: nodes,
	}
}

func (uf *UnionFind) Root(x int) int {
	// "Path splitting" -- setting each node's parent to their current
	// grandfather, is a cache-friendly way of path compression.
	for uf.Nodes[x].parent != x {
		x, uf.Nodes[x].parent = uf.Nodes[x].parent, uf.Nodes[uf.Nodes[x].parent].parent
	}
	return x
}

func (uf *UnionFind) Union(x, y int) {
	rx := uf.Root(x)
	ry := uf.Root(y)
	if rx == ry {
		return
	}
	rankx := uf.Nodes[rx].rank
	ranky := uf.Nodes[ry].rank
	if rankx < ranky {
		rx, ry = ry, rx
	}
	// rx has the larger rank and becomes the root node of the
	// merged set.
	uf.Nodes[ry].parent = rx
	if rankx == ranky {
		uf.Nodes[rx].rank++
	}
}
