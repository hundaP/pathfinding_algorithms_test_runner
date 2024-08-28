package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

	initializeHTMLFile("visualization.html")

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
	finalizeHTMLFile("visualization.html")
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

	nodesInShortestPathOrder := getNodesInShortestPathOrder(endNode)

	var finalMemoryUsage runtime.MemStats
	runtime.ReadMemStats(&finalMemoryUsage)

	endTime := time.Now()
	timeTaken := endTime.Sub(startTime).Nanoseconds()

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

	// Generate the HTML visualization for this algorithm
	htmlContent := generateHTMLGrid(algorithm, grid, visitedNodesInOrder, nodesInShortestPathOrder)
	appendHTMLToFile("visualization.html", htmlContent)

	return nodesInShortestPathOrder, visitedNodesInOrder
}

func getNodesInShortestPathOrder(endNode *maze.Node) []*maze.Node {
	var nodesInShortestPathOrder []*maze.Node
	currentNode := endNode
	for currentNode != nil {
		nodesInShortestPathOrder = append(nodesInShortestPathOrder, currentNode)
		currentNode = currentNode.PreviousNode
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

func initializeHTMLFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	initialContent := "<html><head><style>"
	initialContent += "table { border-collapse: collapse; }"
	initialContent += "td { width: 20px; height: 20px; text-align: center; }"
	initialContent += ".wall { background-color: black; }"
	initialContent += ".visited { background-color: blue; }"
	initialContent += ".path { background-color: yellow !important; }" // Yellow for the path, with higher priority
	initialContent += ".empty { background-color: white; }"
	initialContent += "</style></head><body>"

	if _, err := file.WriteString(initialContent); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func appendHTMLToFile(filename, content string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func finalizeHTMLFile(filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString("</body></html>"); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func generateHTMLGrid(
	algorithm string,
	grid [][]maze.Node,
	visitedNodes []maze.Node,
	nodesInShortestPathOrder []*maze.Node,
) string {
	// Create a set of visited nodes for fast lookup
	visitedNodesSet := make(map[maze.Node]struct{})
	for _, node := range visitedNodes {
		visitedNodesSet[node] = struct{}{}
	}

	// Create a set of nodes in the shortest path for fast lookup
	shortestPathNodesSet := make(map[*maze.Node]struct{})
	for _, node := range nodesInShortestPathOrder {
		shortestPathNodesSet[node] = struct{}{}
	}

	// Start building the HTML content
	html := fmt.Sprintf("<h2>Algorithm: %s</h2>", algorithm)
	html += "<table>"

	// Generate the HTML table for the grid
	for _, row := range grid {
		html += "<tr>"
		for _, node := range row {
			cellClass := "empty"
			if node.IsWall {
				cellClass = "wall"
			} else if _, found := visitedNodesSet[node]; found {
				cellClass = "visited" // Set visited initially
				if _, inPath := shortestPathNodesSet[&node]; inPath {
					cellClass = "path" // Override with path if in shortest path
				}
			} else if _, found := shortestPathNodesSet[&node]; found {
				cellClass = "path" // Shortest path nodes in yellow
			}
			html += fmt.Sprintf("<td class='%s'></td>", cellClass)
		}
		html += "</tr>"
	}

	html += "</table>"
	return html
}
