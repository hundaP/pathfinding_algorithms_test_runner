package algorithms

import (
    "container/heap"
    "math"
)

// Node represents a cell in the grid.
type Node struct {
    X, Y       int
    IsWall     bool
    Distance   int
    IsVisited  bool
    Previous   *Node
}

// PriorityQueue implements heap.Interface for Nodes.
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Distance < pq[j].Distance }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
    *pq = append(*pq, x.(*Node))
}

func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    x := old[n-1]
    *pq = old[0 : n-1]
    return x
}

// Dijkstra performs the Dijkstra's algorithm on a grid.
func Dijkstra(grid [][]*Node, startNode, endNode *Node) []*Node {
    pq := &PriorityQueue{}
    heap.Init(pq)

    for _, row := range grid {
        for _, node := range row {
            node.Distance = math.MaxInt32
            node.Previous = nil
        }
    }
    startNode.Distance = 0
    heap.Push(pq, startNode)

    var visitedNodes []*Node

    for pq.Len() > 0 {
        currentNode := heap.Pop(pq).(*Node)
        if currentNode.IsWall {
            continue
        }
        if currentNode.Distance == math.MaxInt32 {
            return visitedNodes
        }
        currentNode.IsVisited = true
        visitedNodes = append(visitedNodes, currentNode)
        if currentNode == endNode {
            return visitedNodes
        }
        updateUnvisitedNeighbors(currentNode, grid, pq)
    }

    return visitedNodes
}

// updateUnvisitedNeighbors updates the distances of the neighbors of the current node.
func updateUnvisitedNeighbors(node *Node, grid [][]*Node, pq *PriorityQueue) {
    neighbors := getUnvisitedNeighbors(node, grid)
    for _, neighbor := range neighbors {
        if neighbor.IsVisited {
            continue
        }
        newDistance := node.Distance + 1
        if newDistance < neighbor.Distance {
            neighbor.Distance = newDistance
            neighbor.Previous = node
            if !contains(pq, neighbor) {
                heap.Push(pq, neighbor)
            } else {
                updatePriorityQueue(pq, neighbor, newDistance)
            }
        }
    }
}

// getUnvisitedNeighbors returns the unvisited neighbors of a node.
func getUnvisitedNeighbors(node *Node, grid [][]*Node) []*Node {
    neighbors := []*Node{}
    if node.Y > 0 {
        neighbors = append(neighbors, grid[node.Y-1][node.X])
    }
    if node.Y < len(grid)-1 {
        neighbors = append(neighbors, grid[node.Y+1][node.X])
    }
    if node.X > 0 {
        neighbors = append(neighbors, grid[node.Y][node.X-1])
    }
    if node.X < len(grid[0])-1 {
        neighbors = append(neighbors, grid[node.Y][node.X+1])
    }
    return neighbors
}

// contains checks if a node is in the priority queue.
func contains(pq *PriorityQueue, node *Node) bool {
    for _, n := range *pq {
        if n == node {
            return true
        }
    }
    return false
}

// updatePriorityQueue updates the priority queue with a new distance for a node.
func updatePriorityQueue(pq *PriorityQueue, node *Node, newDistance int) {
    for _, n := range *pq {
        if n == node {
            n.Distance = newDistance
            heap.Fix(pq, indexOf(pq, node))
            return
        }
    }
}

// indexOf returns the index of a node in the priority queue.
func indexOf(pq *PriorityQueue, node *Node) int {
    for i, n := range *pq {
        if n == node {
            return i
        }
    }
    return -1
}

