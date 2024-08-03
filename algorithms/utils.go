package algorithms

import (
	"pathfinding_algorithms_test_runner/maze"
)

func getUnvisitedNeighbors(node *maze.Node, grid [][]maze.Node) []*maze.Node {
	neighbors := make([]*maze.Node, 0, 4)
	row, col := node.X, node.Y
	maxRow, maxCol := uint16(len(grid)-1), uint16(len(grid[0])-1)

	if row > 0 && !grid[row-1][col].IsVisited {
		neighbors = append(neighbors, &grid[row-1][col])
	}
	if row < maxRow && !grid[row+1][col].IsVisited {
		neighbors = append(neighbors, &grid[row+1][col])
	}
	if col > 0 && !grid[row][col-1].IsVisited {
		neighbors = append(neighbors, &grid[row][col-1])
	}
	if col < maxCol && !grid[row][col+1].IsVisited {
		neighbors = append(neighbors, &grid[row][col+1])
	}

	return neighbors
}
