package algorithms

import (
	"container/heap"

	"pathfinding_algorithms_test_runner/maze"
)

func DijkstraAlgorithm(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
	rows, cols := len(grid), len(grid[0])
	visitedNodes := make([]maze.Node, 0, rows*cols)

	startNode.Distance = 0
	unvisitedNodes := &PriorityQueue{useAstar: false}
	heap.Init(unvisitedNodes)
	heap.Push(unvisitedNodes, startNode)

	for unvisitedNodes.Len() > 0 {
		closestNode := heap.Pop(unvisitedNodes).(*maze.Node)

		if closestNode.IsWall || closestNode.IsVisited {
			continue
		}

		if closestNode == endNode {
			visitedNodes = append(visitedNodes, *closestNode)
			break
		}

		closestNode.IsVisited = true
		visitedNodes = append(visitedNodes, *closestNode)

		updateUnvisitedNeighbors(closestNode, grid, unvisitedNodes)
	}

	return visitedNodes
}

func updateUnvisitedNeighbors(node *maze.Node, grid [][]maze.Node, unvisitedNodes *PriorityQueue) {
	for _, neighbor := range getUnvisitedNeighbors(node, grid) {
		if neighbor.IsWall || neighbor.IsVisited {
			continue
		}

		newDistance := node.Distance + 1
		if newDistance < neighbor.Distance {
			neighbor.Distance = newDistance
			neighbor.PreviousNode = node
			if !contains(unvisitedNodes, neighbor) {
				heap.Push(unvisitedNodes, neighbor)
			} else {
				unvisitedNodes.update(neighbor, newDistance)
			}
		}
	}
}
