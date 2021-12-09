package main

type UnionFind struct {
	Nodes []int
}

func NewUnionFind(size int) *UnionFind {
	nodes := make([]int, size)
	for i := 0; i < size; i++ {
		nodes[i] = i
	}
	return &UnionFind{
		Nodes: nodes,
	}
}

func (uf *UnionFind) Root(x int) int {
	r := x
	for uf.Nodes[r] != r {
		r = uf.Nodes[r]
	}
	for x != r {
		p := uf.Nodes[x]
		uf.Nodes[x] = r
		x = p
	}
	return r
}

func (uf *UnionFind) Union(x, y int) {
	rx := uf.Root(x)
	ry := uf.Root(y)
	if rx != ry {
		uf.Nodes[rx] = ry
	}
}

