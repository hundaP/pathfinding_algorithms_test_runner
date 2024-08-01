package algorithms

import (
	"math"
	"pathfinding_algorithms_test_runner/maze"
)

func heuristic(node, endNode *maze.Node) float32 {
	// Manhattan distance
	return float32(math.Abs(float64(node.X-endNode.X)) + math.Abs(float64(node.Y-endNode.Y)))
}

func AstarAlgorithm(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
	openList := NewPriorityQueue()
	closedList := make(map[*maze.Node]struct{})
	visitedNodesInOrder := make([]maze.Node, 0, len(grid)*len(grid[0])/2) // Preallocate slice with estimated capacity

	startNode.Distance = 0
	startNode.H = heuristic(startNode, endNode)
	startNode.F = startNode.H
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
				neighbor.H = heuristic(neighbor, endNode)
				neighbor.F = neighbor.Distance + neighbor.H

				if !openList.Contains(neighbor) {
					openList.Enqueue(neighbor)
				} else {
					openList.Update(neighbor, neighbor.F)
				}
			}
		}
	}

	return visitedNodesInOrder
}
