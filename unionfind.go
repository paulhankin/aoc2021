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
	r := x
	for uf.Nodes[r].parent != r {
		r = uf.Nodes[r].parent
	}
	for x != r {
		p := uf.Nodes[x].parent
		uf.Nodes[x].parent = r
		x = p
	}
	return r
}

func (uf *UnionFind) Union(x, y int) {
	rx := uf.Root(x)
	ry := uf.Root(y)
	if rx == ry {
		return
	}
	rankx := uf.Nodes[rx].rank
	ranky := uf.Nodes[ry].rank
	if rankx > ranky {
		rx, ry = ry, rx
		rankx, ranky = ranky, rankx
	}
	uf.Nodes[ry].parent = rx
	if rankx == ranky {
		uf.Nodes[ry].rank++
	}
}
