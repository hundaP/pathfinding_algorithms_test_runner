package algorithms

import "pathfinding_algorithms_test_runner/maze"

// Algorithm is the interface that all pathfinding algorithms must implement.
type Algorithm interface {
	FindPath(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node
}

// Dijkstra implements the Algorithm interface.
type Dijkstra struct{}

func (d Dijkstra) FindPath(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
	return DijkstraAlgorithm(grid, startNode, endNode)
}

// Astar implements the Algorithm interface.
type Astar struct{}

func (a Astar) FindPath(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
	return AstarAlgorithm(grid, startNode, endNode)
}

// BFS implements the Algorithm interface.
type BFS struct{}

func (b BFS) FindPath(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
	return BFSAlgorithm(grid, startNode, endNode)
}

// DFS implements the Algorithm interface.
type DFS struct{}

func (d DFS) FindPath(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
	return DFSAlgorithm(grid, startNode, endNode)
}

// WallFollower implements the Algorithm interface.
type WallFollower struct{}

func (w WallFollower) FindPath(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
	return WallFollowerAlgorithm(grid, startNode, endNode)
}
