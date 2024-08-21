package algorithms

import (
	"pathfinding_algorithms_test_runner/maze"
)

func getUnvisitedNeighbors(node *maze.Node, grid [][]maze.Node) []*maze.Node {
	var neighbors []*maze.Node
	row, col := node.Y, node.X
	gridRows, gridCols := uint16(len(grid)), uint16(len(grid[0]))

	// Use a bitfield to track which neighbors have been checked
	checked := 0

	// Check the top neighbor
	if row > 0 {
		neighbor := &grid[row-1][col]
		if !neighbor.IsVisited {
			neighbors = append(neighbors, neighbor)
		}
		checked |= 1
	}

	// Check the bottom neighbor
	if row < gridRows-1 {
		neighbor := &grid[row+1][col]
		if !neighbor.IsVisited {
			neighbors = append(neighbors, neighbor)
		}
		checked |= 2
	}

	// Check the left neighbor
	if col > 0 {
		neighbor := &grid[row][col-1]
		if !neighbor.IsVisited {
			neighbors = append(neighbors, neighbor)
		}
		checked |= 4
	}

	// Check the right neighbor
	if col < gridCols-1 {
		neighbor := &grid[row][col+1]
		if !neighbor.IsVisited {
			neighbors = append(neighbors, neighbor)
		}
		checked |= 8
	}

	return neighbors
}
