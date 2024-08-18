package algorithms

import (
	"math"

	"pathfinding_algorithms_test_runner/maze"
)

// DFSAlgorithm performs a depth-first search on the grid
func DFSAlgorithm(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
	visitedNodesInOrder := []maze.Node{}
	startNode.Distance = 0
	stack := []*maze.Node{startNode}

	for len(stack) > 0 {
		currentNode := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

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
			neighbor.PreviousID = currentNode.ID
			stack = append(stack, neighbor)
		}
	}

	return visitedNodesInOrder
}
