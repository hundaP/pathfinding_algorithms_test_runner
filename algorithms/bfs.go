package algorithms

import (
	"math"

	"pathfinding_algorithms_test_runner/maze"
)

// BFS performs a breadth-first search on the grid
func BFSAlgorithm(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
	visitedNodesInOrder := []maze.Node{}
	queue := []*maze.Node{}
	startNode.Distance = 0
	queue = append(queue, startNode)

	for len(queue) != 0 {
		currentNode := queue[0]
		queue = queue[1:]

		if currentNode.IsWall {
			continue
		}
		if currentNode.Distance == math.MaxInt32 { // Equivalent to Infinity
			return visitedNodesInOrder
		}
		currentNode.IsVisited = true
		visitedNodesInOrder = append(visitedNodesInOrder, *currentNode)

		if currentNode == endNode {
			return visitedNodesInOrder
		}

		unvisitedNeighbors := getUnvisitedNeighbors(currentNode, grid)
		for _, neighbor := range unvisitedNeighbors {
			neighbor.Distance = currentNode.Distance + 1
			neighbor.PreviousNode = currentNode
			neighbor.IsVisited = true // Mark as visited here
			queue = append(queue, neighbor)
		}
	}

	return visitedNodesInOrder
}
