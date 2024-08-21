package main

import (
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"pathfinding_algorithms_test_runner/algorithms"
	"pathfinding_algorithms_test_runner/maze"
)

type Metrics struct {
	SinglePath        bool      `json:"singlePath"`
	Time              []float64 `json:"time"`
	VisitedNodes      []int     `json:"visitedNodes"`
	VisitedPercentage []float64 `json:"visitedPercentage"`
	PathLength        []int     `json:"pathLength"`
	MemoryUsed        []float64 `json:"memoryUsed"`
}

var (
	algorithmsMap = map[string]algorithms.Algorithm{
		"dijkstra":     algorithms.Dijkstra{},
		"astar":        algorithms.Astar{},
		"bfs":          algorithms.BFS{},
		"dfs":          algorithms.DFS{},
		"wallFollower": algorithms.WallFollower{},
	}

	grids        map[string][][]maze.Node
	startNodes   map[string]*maze.Node
	endNodes     map[string]*maze.Node
	metrics      map[string]*Metrics
	metricsMutex = &sync.Mutex{}
)

func main() {
	metrics = initializeMetrics()
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Replace with your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Routes
	router.GET("/api/maze", mazeHandler)
	router.GET("/api/solution", solutionHandler)

	router.Run("localhost:5000")
}

func mazeHandler(c *gin.Context) {
	mazeSizeStr := c.Query("mazeSize")
	singlePathStr := c.Query("singlePath")

	if mazeSizeStr == "" {
		mazeSizeStr = "50"
	}
	if singlePathStr == "" {
		singlePathStr = "true"
	}

	mazeSize, err := strconv.Atoi(mazeSizeStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid mazeSize"})
		return
	}

	singlePath, err := strconv.ParseBool(singlePathStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid singlePath"})
		return
	}

	grids, startNodes, endNodes = getInitialGrid(mazeSize, mazeSize, singlePath)

	c.JSON(200, gin.H{
		"grids": map[string][][]maze.Node{
			"dijkstra":     grids["dijkstra"],
			"astar":        grids["astar"],
			"bfs":          grids["bfs"],
			"dfs":          grids["dfs"],
			"wallFollower": grids["wallFollower"],
		},
		"startNodes": map[string]*maze.Node{
			"dijkstra":     startNodes["dijkstra"],
			"astar":        startNodes["astar"],
			"bfs":          startNodes["bfs"],
			"dfs":          startNodes["dfs"],
			"wallFollower": startNodes["wallFollower"],
		},
		"endNodes": map[string]*maze.Node{
			"dijkstra":     endNodes["dijkstra"],
			"astar":        endNodes["astar"],
			"bfs":          endNodes["bfs"],
			"dfs":          endNodes["dfs"],
			"wallFollower": endNodes["wallFollower"],
		},
	})
}

func getInitialGrid(
	numRows, numCols int,
	singlePath bool,
) (map[string][][]maze.Node, map[string]*maze.Node, map[string]*maze.Node) {
	grids := make(map[string][][]maze.Node)
	startNodes := make(map[string]*maze.Node)
	endNodes := make(map[string]*maze.Node)

	mazeData := maze.GenerateMaze(numRows, numCols, singlePath)

	grids["dijkstra"] = mazeData["gridDijkstra"].([][]maze.Node)
	grids["astar"] = mazeData["gridAstar"].([][]maze.Node)
	grids["bfs"] = mazeData["gridBFS"].([][]maze.Node)
	grids["dfs"] = mazeData["gridDFS"].([][]maze.Node)
	grids["wallFollower"] = mazeData["gridWallFollower"].([][]maze.Node)

	startNodes["dijkstra"] = mazeData["gridDijkstraStartNode"].(*maze.Node)
	startNodes["astar"] = mazeData["gridAstarStartNode"].(*maze.Node)
	startNodes["bfs"] = mazeData["gridBFSStartNode"].(*maze.Node)
	startNodes["dfs"] = mazeData["gridDFSStartNode"].(*maze.Node)
	startNodes["wallFollower"] = mazeData["gridWallFollowerStartNode"].(*maze.Node)

	endNodes["dijkstra"] = mazeData["gridDijkstraEndNode"].(*maze.Node)
	endNodes["astar"] = mazeData["gridAstarEndNode"].(*maze.Node)
	endNodes["bfs"] = mazeData["gridBFSEndNode"].(*maze.Node)
	endNodes["dfs"] = mazeData["gridDFSEndNode"].(*maze.Node)
	endNodes["wallFollower"] = mazeData["gridWallFollowerEndNode"].(*maze.Node)

	return grids, startNodes, endNodes
}

func solutionHandler(c *gin.Context) {
	var wg sync.WaitGroup
	results := make(map[string]interface{})

	// Reset metrics
	metricsMutex.Lock()
	metrics = initializeMetrics()
	metricsMutex.Unlock()

	for algorithm, alg := range algorithmsMap {
		wg.Add(1)
		go func(algorithm string, alg algorithms.Algorithm) {
			defer wg.Done()
			grid := grids[algorithm]
			startNode := startNodes[algorithm]
			endNode := endNodes[algorithm]

			nodesInShortestPathOrder, visitedNodesInOrder := runAlgorithm(
				algorithm,
				grid,
				startNode,
				endNode,
			)

			results[algorithm] = gin.H{
				"visitedNodesInOrder":      visitedNodesInOrder,
				"nodesInShortestPathOrder": nodesInShortestPathOrder,
				"metrics":                  metrics[algorithm],
			}
		}(algorithm, alg)
	}

	wg.Wait()
	c.JSON(200, results)
}

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
	startNode *maze.Node,
	endNode *maze.Node,
) ([]*maze.Node, []maze.Node) {
	startTime := time.Now()
	var initialMemoryUsage runtime.MemStats
	runtime.ReadMemStats(&initialMemoryUsage)

	visitedNodesInOrder := algorithmsMap[algorithm].FindPath(grid, startNode, endNode)

	var midMemoryUsage runtime.MemStats
	runtime.ReadMemStats(&midMemoryUsage)

	nodesInShortestPathOrder := getNodesInShortestPathOrder(endNode, grid)

	var finalMemoryUsage runtime.MemStats
	runtime.ReadMemStats(&finalMemoryUsage)

	endTime := time.Now()
	timeTaken := endTime.Sub(startTime).Nanoseconds() // Convert to milliseconds

	var memoryUsed float64
	if finalMemoryUsage.HeapAlloc >= initialMemoryUsage.HeapAlloc {
		memoryUsed = float64(
			finalMemoryUsage.HeapAlloc-initialMemoryUsage.HeapAlloc,
		) / (1024 * 1024) // Convert to MB
	} else {
		memoryUsed = 0
	}

	totalNodes := len(grid) * len(grid[0])
	wallNodes := countWallNodes(grid)
	nonWallNodes := totalNodes - wallNodes
	visitedPercentage := (float64(len(visitedNodesInOrder)) / float64(nonWallNodes)) * 100

	metricsMutex.Lock()
	metrics[algorithm].Time = append(metrics[algorithm].Time, float64(timeTaken))
	metrics[algorithm].VisitedNodes = append(
		metrics[algorithm].VisitedNodes,
		len(visitedNodesInOrder),
	)
	metrics[algorithm].VisitedPercentage = append(
		metrics[algorithm].VisitedPercentage,
		visitedPercentage,
	)
	metrics[algorithm].PathLength = append(
		metrics[algorithm].PathLength,
		len(nodesInShortestPathOrder),
	)
	metrics[algorithm].MemoryUsed = append(metrics[algorithm].MemoryUsed, memoryUsed)
	metricsMutex.Unlock()

	return nodesInShortestPathOrder, visitedNodesInOrder
}

func getNodesInShortestPathOrder(endNode *maze.Node, grid [][]maze.Node) []*maze.Node {
	nodeMap := make(map[uint32]*maze.Node)
	for _, row := range grid {
		for _, node := range row {
			nodeMap[node.ID] = &node
		}
	}

	var nodesInShortestPathOrder []*maze.Node
	currentNode := endNode
	for currentNode != nil {
		nodesInShortestPathOrder = append(nodesInShortestPathOrder, currentNode)
		currentNode = nodeMap[currentNode.PreviousID]
	}

	// Reverse the slice
	for i, j := 0, len(nodesInShortestPathOrder)-1; i < j; i, j = i+1, j-1 {
		nodesInShortestPathOrder[i], nodesInShortestPathOrder[j] = nodesInShortestPathOrder[j], nodesInShortestPathOrder[i]
	}

	return nodesInShortestPathOrder
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
