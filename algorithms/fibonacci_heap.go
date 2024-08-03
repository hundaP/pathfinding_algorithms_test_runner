package algorithms

import (
	"pathfinding_algorithms_test_runner/maze"
	"sync"
)

type FibonacciNode struct {
	Node     *maze.Node
	Degree   uint8
	Mark     bool
	Parent   *FibonacciNode
	Children []*FibonacciNode
}

type FibonacciHeap struct {
	MinNode   *FibonacciNode
	NumNodes  uint16
	NodeIndex map[*maze.Node]*FibonacciNode
	RootList  []*FibonacciNode
	nodePool  sync.Pool
}

func NewFibonacciHeap(capacity int) *FibonacciHeap {
	return &FibonacciHeap{
		MinNode:   nil,
		NumNodes:  0,
		NodeIndex: make(map[*maze.Node]*FibonacciNode, capacity),
		RootList:  make([]*FibonacciNode, 0, capacity),
		nodePool: sync.Pool{
			New: func() interface{} {
				return &FibonacciNode{
					Children: make([]*FibonacciNode, 0, 4),
				}
			},
		},
	}
}

func (h *FibonacciHeap) Enqueue(node *maze.Node) {
	fNode := h.nodePool.Get().(*FibonacciNode)
	fNode.Node = node
	fNode.Degree = 0
	fNode.Mark = false
	fNode.Parent = nil
	fNode.Children = fNode.Children[:0] // Reset slice

	h.NodeIndex[node] = fNode
	h.RootList = append(h.RootList, fNode)
	if h.MinNode == nil || node.Distance < h.MinNode.Node.Distance {
		h.MinNode = fNode
	}
	h.NumNodes++
}

func (h *FibonacciHeap) Dequeue() *maze.Node {
	if h.MinNode == nil {
		return nil
	}

	z := h.MinNode

	for _, child := range z.Children {
		h.RootList = append(h.RootList, child)
		child.Parent = nil
	}

	h.removeRoot(z)
	delete(h.NodeIndex, z.Node)

	if len(h.RootList) == 0 {
		h.MinNode = nil
	} else {
		h.MinNode = h.RootList[0]
		h.consolidate()
	}

	h.NumNodes--
	return z.Node
}

func (h *FibonacciHeap) consolidate() {
	if len(h.RootList) == 0 {
		return
	}

	maxDegree := uint8(0)
	for _, node := range h.RootList {
		if node.Degree > maxDegree {
			maxDegree = node.Degree
		}
	}

	A := make([]*FibonacciNode, maxDegree+1)

	for len(h.RootList) > 0 {
		x := h.removeRoot(h.RootList[0])
		d := x.Degree
		for d < uint8(len(A)) && A[d] != nil {
			y := A[d]
			if x.Node.Distance > y.Node.Distance {
				x, y = y, x
			}
			h.link(y, x)
			A[d] = nil
			d++
		}
		if d >= uint8(len(A)) {
			newA := make([]*FibonacciNode, d+1)
			copy(newA, A)
			A = newA
		}
		A[d] = x
	}

	h.MinNode = nil
	for _, node := range A {
		if node != nil {
			h.RootList = append(h.RootList, node)
			if h.MinNode == nil || node.Node.Distance < h.MinNode.Node.Distance {
				h.MinNode = node
			}
		}
	}
}

func (h *FibonacciHeap) link(y, x *FibonacciNode) {
	h.removeRoot(y)
	y.Parent = x
	x.Children = append(x.Children, y)
	x.Degree++
	y.Mark = false
}

func (h *FibonacciHeap) Update(node *maze.Node, newDistance uint32) {
	fNode := h.NodeIndex[node]
	if fNode == nil || newDistance >= node.Distance {
		return
	}
	fNode.Node.Distance = newDistance
	parent := fNode.Parent
	if parent != nil && fNode.Node.Distance < parent.Node.Distance {
		h.cut(fNode, parent)
		h.cascadingCut(parent)
	}
	if fNode.Node.Distance < h.MinNode.Node.Distance {
		h.MinNode = fNode
	}
}

func (h *FibonacciHeap) cut(x, y *FibonacciNode) {
	for i, child := range y.Children {
		if child == x {
			y.Children = append(y.Children[:i], y.Children[i+1:]...)
			break
		}
	}
	y.Degree--
	h.RootList = append(h.RootList, x)
	x.Parent = nil
	x.Mark = false
}

func (h *FibonacciHeap) cascadingCut(y *FibonacciNode) {
	z := y.Parent
	if z != nil {
		if !y.Mark {
			y.Mark = true
		} else {
			h.cut(y, z)
			h.cascadingCut(z)
		}
	}
}

func (h *FibonacciHeap) IsEmpty() bool {
	return h.NumNodes == 0
}

func (h *FibonacciHeap) Contains(node *maze.Node) bool {
	_, exists := h.NodeIndex[node]
	return exists
}

func (h *FibonacciHeap) removeRoot(node *FibonacciNode) *FibonacciNode {
	index := -1
	for i, root := range h.RootList {
		if root == node {
			index = i
			break
		}
	}
	if index != -1 {
		h.RootList = append(h.RootList[:index], h.RootList[index+1:]...)
	}
	return node
}
