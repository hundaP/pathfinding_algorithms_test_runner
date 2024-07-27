package algorithms

import (
	"container/heap"
	"pathfinding_algorithms_test_runner/maze"
)

type PriorityQueue struct {
	items     []*maze.Node
	nodeIndex map[*maze.Node]int
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		items:     []*maze.Node{},
		nodeIndex: make(map[*maze.Node]int),
	}
}

func (pq *PriorityQueue) Len() int { return len(pq.items) }

func (pq *PriorityQueue) Less(i, j int) bool {
	return pq.items[i].F < pq.items[j].F
}

func (pq *PriorityQueue) Swap(i, j int) {
	pq.nodeIndex[pq.items[i]] = j
	pq.nodeIndex[pq.items[j]] = i
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	node := x.(*maze.Node)
	pq.nodeIndex[node] = len(pq.items)
	pq.items = append(pq.items, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := pq.items
	n := len(old)
	node := old[n-1]
	pq.nodeIndex[node] = -1
	pq.items = old[0 : n-1]
	return node
}

func (pq *PriorityQueue) Enqueue(node *maze.Node) {
	heap.Push(pq, node)
}

func (pq *PriorityQueue) Dequeue() *maze.Node {
	return heap.Pop(pq).(*maze.Node)
}

func (pq *PriorityQueue) Update(node *maze.Node, newF float64) {
	node.F = newF
	heap.Fix(pq, pq.nodeIndex[node])
}

func (pq *PriorityQueue) Contains(node *maze.Node) bool {
	_, exists := pq.nodeIndex[node]
	return exists
}

func (pq *PriorityQueue) IsEmpty() bool {
	return pq.Len() == 0
}
