package main

import (
	"log"
	"net/http"
	"pathfinding_algorithms_test_runner/algorithms"
	"pathfinding_algorithms_test_runner/maze"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type Metrics struct {
	Time              []float64
	VisitedNodes      []int
	VisitedPercentage []float64
	PathLength        []int
	MemoryUsed        []float64
}

var algorithmsMap = map[string]algorithms.Algorithm{
	"dijkstra":     algorithms.Dijkstra{},
	"astar":        algorithms.Astar{},
	"bfs":          algorithms.BFS{},
	"dfs":          algorithms.DFS{},
	"wallFollower": algorithms.WallFollower{},
}

var MazeSize int = 14

func initializeMetrics() map[string]*Metrics {
	return map[string]*Metrics{
		"dijkstra":     {},
		"astar":        {},
		"bfs":          {},
		"dfs":          {},
		"wallFollower": {},
	}
}

func runAlgorithm(
	algorithm string,
	grid [][]maze.Node,
	startNode, endNode *maze.Node,
	metrics map[string]*Metrics,
) {
	startTime := time.Now()
	var initialMemoryUsage runtime.MemStats
	runtime.ReadMemStats(&initialMemoryUsage)

	visitedNodes := algorithmsMap[algorithm].FindPath(grid, startNode, endNode)
	pathNodes := getShortestPath(endNode)

	var finalMemoryUsage runtime.MemStats
	runtime.ReadMemStats(&finalMemoryUsage)

	timeTaken := time.Since(startTime).Seconds() // Convert to seconds

	memoryUsed := float64(finalMemoryUsage.HeapAlloc-initialMemoryUsage.HeapAlloc) / (1024 * 1024) // Convert to MB

	totalNodes := len(grid) * len(grid[0])
	wallNodes := countWallNodes(grid)
	nonWallNodes := totalNodes - wallNodes
	visitedPercentage := (float64(len(visitedNodes)) / float64(nonWallNodes)) * 100

	metrics[algorithm].Time = append(metrics[algorithm].Time, timeTaken)
	metrics[algorithm].VisitedNodes = append(metrics[algorithm].VisitedNodes, len(visitedNodes))
	metrics[algorithm].VisitedPercentage = append(metrics[algorithm].VisitedPercentage, visitedPercentage)
	metrics[algorithm].PathLength = append(metrics[algorithm].PathLength, len(pathNodes))
	metrics[algorithm].MemoryUsed = append(metrics[algorithm].MemoryUsed, memoryUsed)

	// Debugging output
	log.Printf("Algorithm: %s", algorithm)
	log.Printf("Time Taken: %f seconds", timeTaken)
	log.Printf("Memory Used: %f MB", memoryUsed)
	log.Printf("Visited Nodes: %d", len(visitedNodes))
	log.Printf("Visited Percentage: %f%%", visitedPercentage)
	log.Printf("Path Length: %d", len(pathNodes))
}

func getShortestPath(endNode *maze.Node) []*maze.Node {
	var path []*maze.Node
	for currentNode := endNode; currentNode != nil; currentNode = currentNode.PreviousNode {
		path = append(path, currentNode)
	}
	// Reverse path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func countWallNodes(grid [][]maze.Node) int {
	count := 0
	for _, row := range grid {
		for _, node := range row {
			if node.IsWall {
				count++
			}
		}
	}
	return count
}

func calculateAverages(metrics map[string]*Metrics) map[string]map[string]float64 {
	averages := make(map[string]map[string]float64)
	for algorithm, metric := range metrics {
		averages[algorithm] = make(map[string]float64)
		numTests := float64(len(metric.Time))
		for _, time := range metric.Time {
			averages[algorithm]["time"] += time
		}
		for _, visitedNodes := range metric.VisitedNodes {
			averages[algorithm]["visitedNodes"] += float64(visitedNodes)
		}
		for _, visitedPercentage := range metric.VisitedPercentage {
			averages[algorithm]["visitedPercentage"] += visitedPercentage
		}
		for _, pathLength := range metric.PathLength {
			averages[algorithm]["pathLength"] += float64(pathLength)
		}
		for _, memoryUsed := range metric.MemoryUsed {
			averages[algorithm]["memoryUsed"] += memoryUsed
		}
		for key := range averages[algorithm] {
			averages[algorithm][key] /= numTests
		}
	}
	return averages
}

func getInitialGrid(numRows, numCols int, singlePath bool) (
	map[string][][]maze.Node,
	map[string]*maze.Node,
	map[string]*maze.Node,
) {
	mazeData := maze.GenerateMaze(numRows, numCols, singlePath)

	grids := map[string][][]maze.Node{
		"dijkstra":     mazeData["gridDijkstra"].([][]maze.Node),
		"astar":        mazeData["gridAstar"].([][]maze.Node),
		"bfs":          mazeData["gridBFS"].([][]maze.Node),
		"dfs":          mazeData["gridDFS"].([][]maze.Node),
		"wallFollower": mazeData["gridWallFollower"].([][]maze.Node),
	}
	startNodes := map[string]*maze.Node{
		"dijkstra":     mazeData["gridDijkstraStartNode"].(*maze.Node),
		"astar":        mazeData["gridAstarStartNode"].(*maze.Node),
		"bfs":          mazeData["gridBFSStartNode"].(*maze.Node),
		"dfs":          mazeData["gridDFSStartNode"].(*maze.Node),
		"wallFollower": mazeData["gridWallFollowerStartNode"].(*maze.Node),
	}
	endNodes := map[string]*maze.Node{
		"dijkstra":     mazeData["gridDijkstraEndNode"].(*maze.Node),
		"astar":        mazeData["gridAstarEndNode"].(*maze.Node),
		"bfs":          mazeData["gridBFSEndNode"].(*maze.Node),
		"dfs":          mazeData["gridDFSEndNode"].(*maze.Node),
		"wallFollower": mazeData["gridWallFollowerEndNode"].(*maze.Node),
	}

	return grids, startNodes, endNodes
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")

	r.GET("/", indexHandler)
	r.GET("/api/metrics", metricsHandler)
	r.GET("/api/initial-grid", gridHandler)

	r.Run(":3000")
}

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func metricsHandler(c *gin.Context) {
	metrics := initializeMetrics()
	grids, startNodes, endNodes := getInitialGrid(MazeSize, MazeSize, true)

	if grids == nil || startNodes == nil || endNodes == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Grid or node data is missing"})
		return
	}

	for algorithm := range algorithmsMap {
		grid := grids[algorithm]
		startNode := startNodes[algorithm]
		endNode := endNodes[algorithm]
		runAlgorithm(algorithm, grid, startNode, endNode, metrics)
	}

	type VisitedNode struct {
		X          uint16 `json:"col"`
		Y          uint16 `json:"row"`
		GridId     string `json:"gridId"`
		NoOfVisits uint8  `json:"noOfVisits"`
	}

	type PathNode struct {
		X      uint16 `json:"col"`
		Y      uint16 `json:"row"`
		GridId string `json:"gridId"`
	}

	type AlgorithmResult struct {
		Metrics             map[string]float64 `json:"metrics"`
		VisitedNodesInOrder []VisitedNode      `json:"visitedNodesInOrder"`
		ShortestPath        []PathNode         `json:"shortestPath"`
	}

	results := make(map[string]AlgorithmResult)

	for algorithm, algoImpl := range algorithmsMap {
		grid, ok := grids[algorithm]
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Grid data for algorithm not found"})
			return
		}

		startNode, ok := startNodes[algorithm]
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Start node for algorithm not found"})
			return
		}

		endNode, ok := endNodes[algorithm]
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "End node for algorithm not found"})
			return
		}

		visitedNodes := algoImpl.FindPath(grid, startNode, endNode)
		shortestPath := getShortestPath(endNode)

		visitedNodesInOrder := make([]VisitedNode, len(visitedNodes))
		for i, node := range visitedNodes {
			visitedNodesInOrder[i] = VisitedNode{
				X:          node.X,
				Y:          node.Y,
				GridId:     algorithm,
				NoOfVisits: node.NoOfVisits,
			}
		}

		shortestPathNodes := make([]PathNode, len(shortestPath))
		for i, node := range shortestPath {
			shortestPathNodes[i] = PathNode{
				X:      node.X,
				Y:      node.Y,
				GridId: algorithm,
			}
		}

		results[algorithm] = AlgorithmResult{
			Metrics:             calculateAverages(map[string]*Metrics{algorithm: metrics[algorithm]})[algorithm],
			VisitedNodesInOrder: visitedNodesInOrder,
			ShortestPath:        shortestPathNodes,
		}
	}

	c.JSON(http.StatusOK, results)
}

func gridHandler(c *gin.Context) {
	grids, _, _ := getInitialGrid(MazeSize, MazeSize, false)

	if grids == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Grid data is missing"})
		return
	}

	gridData := make(map[string][][]map[string]interface{})
	for algorithm, grid := range grids {
		gridData[algorithm] = make([][]map[string]interface{}, len(grid))
		for y, row := range grid {
			gridData[algorithm][y] = make([]map[string]interface{}, len(row))
			for x, node := range row {
				gridData[algorithm][y][x] = map[string]interface{}{
					"X":         node.X,
					"Y":         node.Y,
					"className": getClassName(node),
				}
			}
		}
	}

	c.JSON(http.StatusOK, gridData)
}

func getClassName(node maze.Node) string {
	switch {
	case node.IsStart:
		return "node-start"
	case node.IsEnd:
		return "node-end"
	case node.IsWall:
		return "node-wall"
	case node.IsVisited:
		return "node-visited"
	case node.PreviousNode != nil:
		return "node-shortest-path"
	default:
		return ""
	}
}
