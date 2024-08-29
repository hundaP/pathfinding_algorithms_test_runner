package algorithms

import "pathfinding_algorithms_test_runner/maze"

// Directions
const (
	Left = iota
	Up
	Right
	Down
)

// WallFollowerAlgorithm performs the wall follower algorithm on the grid
func WallFollowerAlgorithm(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
	visitedNodesInOrder := []maze.Node{}
	startNode.Distance = 0
	currentNode := startNode
	var previousNode *maze.Node
	currentDirection := Right                     // Assume starting direction is right
	maxIterations := len(grid) * len(grid[0]) * 2 // Set a maximum number of iterations

	for currentNode != endNode && maxIterations > 0 {
		currentNode.IsVisited = true
		currentNode.NoOfVisits++
		visitedNodesInOrder = append(visitedNodesInOrder, *currentNode)

		neighbors := getPrioritizedNeighbors(currentNode, grid, currentDirection)
		var nextNode *maze.Node

		for _, neighbor := range neighbors {
			if !neighbor.IsWall && (neighbor.NoOfVisits == 0 || neighbor == endNode) {
				nextNode = neighbor
				break
			}
		}

		if nextNode != nil {
			nextNode.Distance = currentNode.Distance + 1
			nextNode.PreviousNode = currentNode
			previousNode = currentNode
			currentDirection = getDirection(currentNode, nextNode)
			currentNode = nextNode
		} else {
			// Backtrack
			if previousNode != nil {
				currentNode = previousNode
				previousNode = currentNode.PreviousNode
				currentDirection = (currentDirection + 2) % 4 // Reverse direction
			} else {
				break // No more backtracking possible
			}
		}

		maxIterations--
	}

	return visitedNodesInOrder
}

// getPrioritizedNeighbors returns the neighbors of the node in the order of left, up, right, down relative to the current direction
func getPrioritizedNeighbors(node *maze.Node, grid [][]maze.Node, direction int) []*maze.Node {
	var neighbors []*maze.Node
	row, col := node.Y, node.X
	maxRow, maxCol := uint16(len(grid)-1), uint16(len(grid[0])-1)

	// Define the order of directions to check based on the current direction
	var directions [4][2]int
	switch direction {
	case Left:
		directions = [4][2]int{{0, -1}, {1, 0}, {-1, 0}, {0, 1}} // Left, Down, Up, Right
	case Up:
		directions = [4][2]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}} // Up, Left, Right, Down
	case Right:
		directions = [4][2]int{{0, 1}, {-1, 0}, {1, 0}, {0, -1}} // Right, Up, Down, Left
	case Down:
		directions = [4][2]int{{1, 0}, {0, 1}, {0, -1}, {-1, 0}} // Down, Right, Left, Up
	}

	// Check neighbors in the prioritized order
	for _, dir := range directions {
		newRow, newCol := row+uint16(dir[0]), col+uint16(dir[1])
		if newRow <= maxRow && newCol <= maxCol {
			neighbors = append(neighbors, &grid[newRow][newCol])
		}
	}

	return neighbors
}

// getDirection determines the direction from currentNode to nextNode
func getDirection(currentNode, nextNode *maze.Node) int {
	if nextNode.Y == currentNode.Y {
		if nextNode.X < currentNode.X {
			return Left
		}
		return Right
	}
	if nextNode.Y < currentNode.Y {
		return Up
	}
	return Down
}
