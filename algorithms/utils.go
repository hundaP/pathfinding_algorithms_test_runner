package algorithms

import (
	"pathfinding_algorithms_test_runner/maze"
)

func getUnvisitedNeighbors(node *maze.Node, grid [][]maze.Node) []*maze.Node {
    neighbors := []*maze.Node{}
    row, col := node.X, node.Y

    if row > 0 {
        neighbors = append(neighbors, &grid[row-1][col])
    }
    if row < len(grid)-1 {
        neighbors = append(neighbors, &grid[row+1][col])
    }
    if col > 0 {
        neighbors = append(neighbors, &grid[row][col-1])
    }
    if col < len(grid[0])-1 {
        neighbors = append(neighbors, &grid[row][col+1])
    }

    unvisitedNeighbors := []*maze.Node{}
    for _, neighbor := range neighbors {
        if !neighbor.IsVisited {
            unvisitedNeighbors = append(unvisitedNeighbors, neighbor)
        }
    }

    return unvisitedNeighbors
}
