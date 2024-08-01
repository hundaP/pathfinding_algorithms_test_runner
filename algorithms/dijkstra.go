package algorithms

import (
	"pathfinding_algorithms_test_runner/maze"
)

func DijkstraAlgorithm(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
	openList := NewPriorityQueue()
	closedList := make(map[*maze.Node]struct{})
	visitedNodesInOrder := make([]maze.Node, 0, len(grid)*len(grid[0])/2) // Preallocate slice with estimated capacity

	startNode.Distance = 0
	openList.Enqueue(startNode)

	for !openList.IsEmpty() {
		currentNode := openList.Dequeue()
		closedList[currentNode] = struct{}{}
		visitedNodesInOrder = append(visitedNodesInOrder, *currentNode)

		if currentNode == endNode {
			return visitedNodesInOrder
		}

		for _, neighbor := range getUnvisitedNeighbors(currentNode, grid) {
			if _, found := closedList[neighbor]; found || neighbor.IsWall {
				continue
			}

			tentativeGScore := currentNode.Distance + 1

			if tentativeGScore < neighbor.Distance {
				neighbor.PreviousNode = currentNode
				neighbor.Distance = tentativeGScore

				if !openList.Contains(neighbor) {
					openList.Enqueue(neighbor)
				} else {
					openList.Update(neighbor, neighbor.Distance)
				}
			}
		}
	}

	return visitedNodesInOrder
}
