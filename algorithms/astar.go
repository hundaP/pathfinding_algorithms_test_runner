package algorithms

import (
	"container/heap"
	"math"
	"pathfinding_algorithms_test_runner/maze"
)

func heuristic(node, endNode *maze.Node) float32 {
	// Manhattan distance
	return float32(math.Abs(float64(node.X-endNode.X)) + math.Abs(float64(node.Y-endNode.Y)))
}

func AstarAlgorithm(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
	openSet := &PriorityQueue{useAstar: true}
	heap.Init(openSet)
	heap.Push(openSet, startNode)

	rows, cols := uint16(len(grid)), uint16(len(grid[0]))
	totalNodes := uint32(rows) * uint32(cols) // Correct calculation of totalNodes
	cameFrom := make([]*maze.Node, totalNodes)
	gScore := make([]float32, totalNodes)
	for i := range gScore {
		gScore[i] = math.MaxFloat32
	}
	startIndex := uint32(startNode.Y)*uint32(cols) + uint32(startNode.X)
	gScore[startIndex] = 0

	startNode.H = heuristic(startNode, endNode)
	startNode.F = startNode.H

	visitedNodesInOrder := []maze.Node{}
	closedSet := make([]bool, totalNodes)

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*maze.Node)
		visitedNodesInOrder = append(visitedNodesInOrder, *current)

		if current == endNode {
			return visitedNodesInOrder
		}

		currentIndex := uint32(current.Y)*uint32(cols) + uint32(current.X)
		closedSet[currentIndex] = true

		for _, neighbor := range getUnvisitedNeighbors(current, grid) {
			// Check if neighbor is within grid boundaries
			if neighbor.Y >= rows || neighbor.X >= cols {
				continue
			}

			neighborIndex := uint32(neighbor.Y)*uint32(cols) + uint32(neighbor.X)
			if neighbor.IsWall || closedSet[neighborIndex] {
				continue
			}

			tentativeGScore := gScore[currentIndex] + 1

			if tentativeGScore < gScore[neighborIndex] {
				cameFrom[neighborIndex] = current
				gScore[neighborIndex] = tentativeGScore
				neighbor.H = heuristic(neighbor, endNode)
				neighbor.F = tentativeGScore + neighbor.H
				neighbor.PreviousNode = current // Set the previous node

				// Only push the neighbor to the open set if it's not already in it
				if !closedSet[neighborIndex] {
					heap.Push(openSet, neighbor)
				} else {
					// Update the position of the neighbor in the priority queue
					heap.Fix(openSet, findIndex(openSet, neighbor))
				}
			}
		}
	}

	return visitedNodesInOrder
}

func findIndex(pq *PriorityQueue, node *maze.Node) int {
	for i, n := range pq.nodes {
		if n == node {
			return i
		}
	}
	return -1
}
