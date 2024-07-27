package algorithms

import (
	"pathfinding_algorithms_test_runner/maze"
)

func DijkstraAlgorithm(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
    openList := NewPriorityQueue()
    closedList := make(map[*maze.Node]bool)
    visitedNodesInOrder := []maze.Node{}

    startNode.Distance = 0
    openList.Enqueue(startNode)

    for !openList.IsEmpty() {
        currentNode := openList.Dequeue()
        closedList[currentNode] = true
        visitedNodesInOrder = append(visitedNodesInOrder, *currentNode)

        if currentNode == endNode {
            return visitedNodesInOrder
        }

        neighbors := getUnvisitedNeighbors(currentNode, grid)
        for _, neighbor := range neighbors {
            if closedList[neighbor] || neighbor.IsWall {
                continue
            }

            gScore := currentNode.Distance + 1

            if !openList.Contains(neighbor) {
                neighbor.Distance = gScore
                neighbor.PreviousNode = currentNode
                openList.Enqueue(neighbor)
            } else if gScore < neighbor.Distance {
                neighbor.Distance = gScore
                neighbor.PreviousNode = currentNode
                openList.Update(neighbor, neighbor.Distance)
            }
        }
    }

    return visitedNodesInOrder
}