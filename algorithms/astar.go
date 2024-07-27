package algorithms

import (
	"container/heap"
	"math"
)

type Node struct {
	row, col, distance, h, f int
	isWall, isVisited        bool
	previousNode             *Node
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].f < pq[j].f
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*Node)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*pq = old[0 : n-1]
	return node
}

func (pq *PriorityQueue) update(node *Node, f int) {
	node.f = f
	heap.Fix(pq, node.index)
}

func heuristic(node, endNode *Node) int {
	return int(math.Abs(float64(node.row-endNode.row)) + math.Abs(float64(node.col-endNode.col)))
}

func getUnvisitedNeighbors(node *Node, grid [][]*Node) []*Node {
	var neighbors []*Node
	row, col := node.row, node.col
	if row > 0 {
		neighbors = append(neighbors, grid[row-1][col])
	}
	if row < len(grid)-1 {
		neighbors = append(neighbors, grid[row+1][col])
	}
	if col > 0 {
		neighbors = append(neighbors, grid[row][col-1])
	}
	if col < len(grid[0])-1 {
		neighbors = append(neighbors, grid[row][col+1])
	}
	var unvisited []*Node
	for _, neighbor := range neighbors {
		if !neighbor.isVisited {
			unvisited = append(unvisited, neighbor)
		}
	}
	return unvisited
}

func astar(grid [][]*Node, startNode, endNode *Node) []*Node {
	openList := make(PriorityQueue, 0)
	heap.Init(&openList)
	closedList := make(map[*Node]bool)
	var visitedNodesInOrder []*Node
	startNode.distance = 0
	startNode.h = heuristic(startNode, endNode)
	startNode.f = startNode.h
	heap.Push(&openList, startNode)
	for len(openList) != 0 {
		currentNode := heap.Pop(&openList).(*Node)
		closedList[currentNode] = true
		visitedNodesInOrder = append(visitedNodesInOrder, currentNode)
		if currentNode == endNode {
			return visitedNodesInOrder
		}
		neighbors := getUnvisitedNeighbors(currentNode, grid)
		for _, neighbor := range neighbors {
			if closedList[neighbor] || neighbor.isWall {
				continue
			}
			gScore := currentNode.distance + 1
			hScore := neighbor.h
			if hScore == 0 {
				hScore = heuristic(neighbor, endNode)
			}
			neighbor.h = hScore
			if !neighbor.isVisited {
				heap.Push(&openList, neighbor)
				neighbor.distance = gScore
				neighbor.f = gScore + hScore
				neighbor.previousNode = currentNode
			} else if gScore < neighbor.distance {
				neighbor.distance = gScore
				neighbor.f = gScore + hScore
				neighbor.previousNode = currentNode
				openList.update(neighbor, neighbor.f)
			}
		}
	}
	return visitedNodesInOrder
}
