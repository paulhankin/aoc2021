package main

import "container/heap"

type NodeCost struct {
	Node int
	Cost int
}

type openSet struct {
	nodes map[int]int
	queue []NodeCost
}

type osAsHeap struct {
	o *openSet
}

func (os osAsHeap) Push(xi interface{}) {
	o := os.o
	x := xi.(NodeCost)
	o.nodes[x.Node] = len(o.queue)
	o.queue = append(o.queue, x)
}

func (os osAsHeap) Pop() interface{} {
	o := os.o
	x := o.queue[len(o.queue)-1]
	o.queue = o.queue[:len(o.queue)-1]
	delete(o.nodes, x.Node)
	return x
}

func (os osAsHeap) Len() int {
	o := os.o
	return len(o.queue)
}

func (os osAsHeap) Less(i, j int) bool {
	o := os.o
	return o.queue[i].Cost < o.queue[j].Cost
}

func (os osAsHeap) Swap(i, j int) {
	o := os.o
	o.queue[i], o.queue[j] = o.queue[j], o.queue[i]
	o.nodes[o.queue[i].Node] = i
	o.nodes[o.queue[j].Node] = j
}

func (o *openSet) update(state, cost int) {
	if idx, ok := o.nodes[state]; ok {
		o.queue[idx].Cost = cost
		heap.Fix(osAsHeap{o}, idx)
		return
	}
	heap.Push(osAsHeap{o}, NodeCost{state, cost})
}

func newOpenSet() *openSet {
	return &openSet{
		nodes: map[int]int{},
		queue: nil,
	}
}

func (o *openSet) len() int {
	return len(o.queue)
}

func (o *openSet) Pop() int {
	nc := heap.Pop(osAsHeap{o}).(NodeCost)
	return nc.Node
}

// MinPath finds a minimum path from start to target,
// returning the cost, or -1 if no such path exists.
// adjacent(nc, i) returns the edges from a given node along
// with the cost of traversing from i to the new node.
// It passes in a slice to append nodecosts to (to avoid garbage).
// heuristic(i) is an estimate of the cost of travelling
// from i to target, never over-estimating (ie: is admissible).
func MinPath(start, target int, adjacent func(int, []NodeCost) []NodeCost, heuristic func(int) int, cacheHeuristic bool) int {
	openSet := newOpenSet()
	openSet.update(start, heuristic(start))
	var hc map[int]int
	if cacheHeuristic {
		hc = map[int]int{}
	}
	gs := map[int]int{}
	gs[start] = 0

	nc := make([]NodeCost, 0, 120)

	for openSet.len() > 0 {
		current := openSet.Pop()
		if current == target {
			return gs[current]
		}
		for _, ec := range adjacent(current, nc[:0]) {
			tgs := gs[current] + ec.Cost
			if cgs, ok := gs[ec.Node]; !ok || tgs < cgs {
				gs[ec.Node] = tgs
				h := 0
				if cacheHeuristic {
					var ok bool
					h, ok = hc[ec.Node]
					if !ok {
						h = heuristic(ec.Node)
					}
					hc[ec.Node] = h
				} else {
					h = heuristic(ec.Node)
				}
				openSet.update(ec.Node, tgs+h)
			}
		}
	}
	return -1
}
