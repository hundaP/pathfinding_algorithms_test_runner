package algorithms

import (
	"pathfinding_algorithms_test_runner/maze"
)

func getUnvisitedNeighbors(node *maze.Node, grid [][]maze.Node) []*maze.Node {
	neighbors := make([]*maze.Node, 0, 4) // Preallocate slice with capacity 4
	row, col := node.X, node.Y
	maxRow, maxCol := uint16(len(grid)-1), uint16(len(grid[0])-1)

	if row > 0 {
		neighbor := &grid[row-1][col]
		if !neighbor.IsVisited {
			neighbors = append(neighbors, neighbor)
		}
	}
	if row < maxRow {
		neighbor := &grid[row+1][col]
		if !neighbor.IsVisited {
			neighbors = append(neighbors, neighbor)
		}
	}
	if col > 0 {
		neighbor := &grid[row][col-1]
		if !neighbor.IsVisited {
			neighbors = append(neighbors, neighbor)
		}
	}
	if col < maxCol {
		neighbor := &grid[row][col+1]
		if !neighbor.IsVisited {
			neighbors = append(neighbors, neighbor)
		}
	}

	return neighbors
}
