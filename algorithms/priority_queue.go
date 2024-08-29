package algorithms

import (
	"container/heap"
	"pathfinding_algorithms_test_runner/maze"
)

type PriorityQueue struct {
	nodes    []*maze.Node
	useAstar bool
}

func (pq PriorityQueue) Len() int { return len(pq.nodes) }

func (pq PriorityQueue) Less(i, j int) bool {
	if pq.useAstar {
		if pq.nodes[i].F == pq.nodes[j].F {
			return pq.nodes[i].G < pq.nodes[j].G
		}
		return pq.nodes[i].F < pq.nodes[j].F
	}
	return pq.nodes[i].Distance < pq.nodes[j].Distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq.nodes[i], pq.nodes[j] = pq.nodes[j], pq.nodes[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	node := x.(*maze.Node)
	pq.nodes = append(pq.nodes, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := pq.nodes
	n := len(old)
	node := old[n-1]
	pq.nodes = old[0 : n-1]
	return node
}

func (pq *PriorityQueue) update(node *maze.Node, distance uint32) {
	for i, n := range pq.nodes {
		if n == node {
			pq.nodes[i].Distance = distance
			heap.Fix(pq, i)
			break
		}
	}
}

func (pq *PriorityQueue) IndexOf(node *maze.Node) int {
	for i, n := range pq.nodes {
		if n == node {
			return i
		}
	}
	return -1
}
