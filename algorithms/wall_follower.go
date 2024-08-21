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
	currentDirection := Right // Assume starting direction is right

	for currentNode != endNode {
		currentNode.IsVisited = true
		currentNode.NoOfVisits++
		visitedNodesInOrder = append(visitedNodesInOrder, *currentNode)

		neighbors := getPrioritizedNeighbors(currentNode, grid, currentDirection)
		var nextNode *maze.Node

		for _, neighbor := range neighbors {
			if neighbor != previousNode && !neighbor.IsWall && !neighbor.IsVisited {
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
			if previousNode != nil {
				currentNode = previousNode
				previousNode = currentNode.PreviousNode
			} else {
				break
			}
		}
	}

	return visitedNodesInOrder
}

// getPrioritizedNeighbors returns the neighbors of the node in the order of left, up, right, down relative to the current direction
func getPrioritizedNeighbors(node *maze.Node, grid [][]maze.Node, direction int) []*maze.Node {
	var neighbors []*maze.Node
	row, col := node.Y, node.X

	switch direction {
	case Left:
		// Prioritize Down, Left, Up, Right
		if row < uint16(len(grid)-1) {
			neighbors = append(neighbors, &grid[row+1][col]) // Down
		}
		if col > 0 {
			neighbors = append(neighbors, &grid[row][col-1]) // Left
		}
		if row > 0 {
			neighbors = append(neighbors, &grid[row-1][col]) // Up
		}
		if col < uint16(len(grid[0])-1) {
			neighbors = append(neighbors, &grid[row][col+1]) // Right
		}
	case Up:
		// Prioritize Left, Up, Right, Down
		if col > 0 {
			neighbors = append(neighbors, &grid[row][col-1]) // Left
		}
		if row > 0 {
			neighbors = append(neighbors, &grid[row-1][col]) // Up
		}
		if col < uint16(len(grid[0])-1) {
			neighbors = append(neighbors, &grid[row][col+1]) // Right
		}
		if row < uint16(len(grid)-1) {
			neighbors = append(neighbors, &grid[row+1][col]) // Down
		}
	case Right:
		// Prioritize Up, Right, Down, Left
		if row > 0 {
			neighbors = append(neighbors, &grid[row-1][col]) // Up
		}
		if col < uint16(len(grid[0])-1) {
			neighbors = append(neighbors, &grid[row][col+1]) // Right
		}
		if row < uint16(len(grid)-1) {
			neighbors = append(neighbors, &grid[row+1][col]) // Down
		}
		if col > 0 {
			neighbors = append(neighbors, &grid[row][col-1]) // Left
		}
	case Down:
		// Prioritize Right, Down, Left, Up
		if col < uint16(len(grid[0])-1) {
			neighbors = append(neighbors, &grid[row][col+1]) // Right
		}
		if row < uint16(len(grid)-1) {
			neighbors = append(neighbors, &grid[row+1][col]) // Down
		}
		if col > 0 {
			neighbors = append(neighbors, &grid[row][col-1]) // Left
		}
		if row > 0 {
			neighbors = append(neighbors, &grid[row-1][col]) // Up
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
