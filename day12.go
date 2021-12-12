package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type node12 struct {
	name  string
	small bool
	edges []int
}

type graph12 struct {
	nodes      []node12
	edges      [][2]int
	start      int
	end        int
	startLarge int
}

func readDay12() (*graph12, error) {
	f, err := os.Open("day12.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var r [][2]string
	for scanner.Scan() {
		parts := strings.Split(strings.TrimSpace(scanner.Text()), "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("missing from/to in %q", scanner.Text())
		}
		r = append(r, [2]string{parts[0], parts[1]})
	}
	nm := map[string]bool{}
	var nodes []node12
	var addNode = func(name string) {
		if nm[name] {
			return
		}
		nm[name] = true
		nodes = append(nodes, node12{
			name:  name,
			small: name[0] >= 'a' && name[0] <= 'z',
		})
	}
	for _, e := range r {
		addNode(e[0])
		addNode(e[1])
	}
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].small != nodes[j].small {
			return nodes[i].small
		}
		return nodes[i].name <= nodes[j].name
	})
	findNode := func(s string) int {
		for i, n := range nodes {
			if n.name == s {
				return i
			}
		}
		log.Fatalf("not found node %q", s)
		return -1
	}
	var edges [][2]int
	for _, e := range r {
		n0 := findNode(e[0])
		n1 := findNode(e[1])
		edges = append(edges, [2]int{n0, n1})
		nodes[n0].edges = append(nodes[n0].edges, n1)
		nodes[n1].edges = append(nodes[n1].edges, n0)
	}

	var startLarge int
	for i, n := range nodes {
		if !n.small {
			startLarge = i
			break
		}
	}

	return &graph12{
		nodes:      nodes,
		edges:      edges,
		start:      findNode("start"),
		end:        findNode("end"),
		startLarge: startLarge,
	}, nil

}

func countPathsSearch(g *graph12, seen []int, node int, visitSmall bool, path []string) int {
	if g.end == node {
		// fmt.Println(append(path, "end"))
		return 1
	}
	path = append(path, g.nodes[node].name)
	if g.nodes[node].small {
		seen[node]++
	}
	n := 0
	for _, out := range g.nodes[node].edges {
		allowedVisits := 1 + b2i(visitSmall)
		if seen[out] >= allowedVisits || out == g.start {
			continue
		}
		n += countPathsSearch(g, seen, out, visitSmall && seen[out] == 0, path)
	}
	if g.nodes[node].small {
		seen[node]--
	}
	return n
}

func countPaths(g *graph12, extraVisit bool) int {
	seen := make([]int, len(g.nodes))
	return countPathsSearch(g, seen, g.start, extraVisit, nil)
}

func day12() error {
	g, err := readDay12()
	if err != nil {
		return err
	}
	partPrint(1, countPaths(g, false))
	partPrint(2, countPaths(g, true))
	return nil
}

func init() {
	RegisterDay(12, day12)
}
