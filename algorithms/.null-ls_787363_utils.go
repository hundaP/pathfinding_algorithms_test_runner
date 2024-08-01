package algorithms

import (
	"pathfinding_algorithms_test_runner/maze"
)

func getUnvisitedNeighbors(node *maze.Node, grid [][]maze.Node) []*maze.Node {
	neighbors := []*maze.Node{}
	row, col := node.X, node.Y

	if row > 0 && !grid[row-1][col].IsVisited {
		neighbors = append(neighbors, &grid[row-1][col])
	}
	if row < len(grid)-1 && !grid[row+1][col].IsVisited {
		neighbors = append(neighbors, &grid[row+1][col])
	}
	if col > 0 && !grid[row][col-1].IsVisited {
		neighbors = append(neighbors, &grid[row][col-1])
	}
	if col < len(grid[0])-1 && !grid[row][col+1].IsVisited {
		neighbors = append(neighbors, &grid[row][col+1])
	}

	return neighbors
}
