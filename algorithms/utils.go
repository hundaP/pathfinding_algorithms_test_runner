package algorithms

import (
	"pathfinding_algorithms_test_runner/maze"
)

func getUnvisitedNeighbors(node *maze.Node, grid [][]maze.Node) []*maze.Node {
	neighbors := make([]*maze.Node, 0, 4) // Preallocate slice with capacity 4
	row, col := node.Y, node.X
	maxRow, maxCol := uint16(len(grid)-1), uint16(len(grid[0])-1)

	// Check the top neighbor
	if row > 0 {
		neighbor := &grid[row-1][col]
		if !neighbor.IsVisited {
			neighbors = append(neighbors, neighbor)
		}
	}

	// Check the bottom neighbor
	if row < maxRow {
		neighbor := &grid[row+1][col]
		if !neighbor.IsVisited {
			neighbors = append(neighbors, neighbor)
		}
	}

	// Check the left neighbor
	if col > 0 {
		neighbor := &grid[row][col-1]
		if !neighbor.IsVisited {
			neighbors = append(neighbors, neighbor)
		}
	}

	// Check the right neighbor
	if col < maxCol {
		neighbor := &grid[row][col+1]
		if !neighbor.IsVisited {
			neighbors = append(neighbors, neighbor)
		}
	}

	return neighbors
}
