package algorithms

import (
	"pathfinding_algorithms_test_runner/maze"
)

func DijkstraAlgorithm(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
	rows, cols := len(grid), len(grid[0])
	visitedNodes := make([]maze.Node, 0, rows*cols)

	startNode.Distance = 0
	unvisitedNodes := NewFibonacciHeap(rows * cols)
	unvisitedNodes.Enqueue(startNode)

	for !unvisitedNodes.IsEmpty() {
		closestNode := unvisitedNodes.Dequeue()

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

func updateUnvisitedNeighbors(node *maze.Node, grid [][]maze.Node, unvisitedNodes *FibonacciHeap) {
	for _, neighbor := range getUnvisitedNeighbors(node, grid) {
		if neighbor.IsWall || neighbor.IsVisited {
			continue
		}

		newDistance := node.Distance + 1
		if newDistance < neighbor.Distance {
			neighbor.Distance = newDistance
			neighbor.PreviousNode = node
			if !unvisitedNodes.Contains(neighbor) {
				unvisitedNodes.Enqueue(neighbor)
			} else {
				unvisitedNodes.Update(neighbor, newDistance)
			}
		}
	}
}
