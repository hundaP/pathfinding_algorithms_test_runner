export function generateMaze(numRows, numCols, singlePath) {
    let maze = new Maze(numRows, numCols);
  
    // Generate the maze
    while (maze.stack.length > 0 || !maze.currentCell.visited) {
      maze.generateMazeNotGlobal();
    }
  
    if(!singlePath){
      for (let i = 0; i < numRows * numCols * 0.1; i++) { // remove 10% of the walls
          let x = Math.floor(Math.random() * (numRows - 2)) + 1;
          let y = Math.floor(Math.random() * (numCols - 2)) + 1;
          maze.grid[x][y].isWall = false;
      }
  }
  
    // Convert the maze grid to the format expected by the rest of your application
    let grid1 = maze.grid.map(row => row.map(cell => createNode(cell.x, cell.y, cell.isWall, maze.start, maze.end, 1)));
    let grid2 = maze.grid.map(row => row.map(cell => createNode(cell.x, cell.y, cell.isWall, maze.start, maze.end, 2)));
    let grid3 = maze.grid.map(row => row.map(cell => createNode(cell.x, cell.y, cell.isWall, maze.start, maze.end, 3)));
    let grid4 = maze.grid.map(row => row.map(cell => createNode(cell.x, cell.y, cell.isWall, maze.start, maze.end, 4)));
    let grid5 = maze.grid.map(row => row.map(cell => createNode(cell.x, cell.y, cell.isWall, maze.start, maze.end, 5)));
  
    // Find start and end nodes
    let gridDijkstraStartNode = grid1[maze.start.y][maze.start.x];
    let gridDijkstraEndNode = grid1[maze.end.y][maze.end.x];
    let gridAstarStartNode = grid2[maze.start.y][maze.start.x];
    let gridAstarEndNode = grid2[maze.end.y][maze.end.x];
    let gridBFSStartNode = grid3[maze.start.y][maze.start.x];
    let gridBFSEndNode = grid3[maze.end.y][maze.end.x];
    let gridDFSStartNode = grid4[maze.start.y][maze.start.x];
    let gridDFSEndNode = grid4[maze.end.y][maze.end.x];
    let gridWallFollowerStartNode = grid5[maze.start.y][maze.start.x];
    let gridWallFollowerEndNode = grid5[maze.end.y][maze.end.x];
  
  
    return {
      gridDijsktra: grid1,
      gridAstar: grid2,
      gridBFS: grid3,
      gridDFS: grid4,
      gridWallFollower: grid5,
      gridDijkstraStartNode,
      gridDijkstraEndNode,
      gridAstarStartNode,
      gridAstarEndNode,
      gridBFSStartNode,
      gridBFSEndNode,
      gridDFSStartNode,
      gridDFSEndNode,
      gridWallFollowerStartNode,
      gridWallFollowerEndNode
    };
  }
  
  class Cell {
    constructor(x, y, isWall = false) {
      this.x = x;
      this.y = y;
      this.isWall = isWall;
      this.visited = false;
    }
  }
  
  class Maze {
    constructor(width, height) {
      this.width = width * 2 + 1; // Adjust for walls
      this.height = height * 2 + 1; // Adjust for walls
      this.grid = [];
      this.stack = [];
  
      for (let y = 0; y < this.height; y++) {
        let row = [];
        for (let x = 0; x < this.width; x++) {
          let isWall = x % 2 === 0 || y % 2 === 0; // All cells at even indices are walls
          row.push(new Cell(x, y, isWall));
        }
        this.grid.push(row);
      }
  
      this.currentCell = this.grid[1][1]; // Start at the top-left corner
      this.start = this.grid[1][1]; // Top-left corner
      this.end = this.grid[this.height - 2][this.width - 2]; // Bottom-right corner
    }
  
    getCell(x, y) {
      if (x < 0 || y < 0 || x >= this.width || y >= this.height) {
        return undefined;
      }
      return this.grid[y][x];
    }
  
    getNeighbors(cell) {
      let neighbors = [];
  
      // Get the cells above, below, to the left, and to the right of the given cell
      let top = this.getCell(cell.x, cell.y - 2);
      let right = this.getCell(cell.x + 2, cell.y);
      let bottom = this.getCell(cell.x, cell.y + 2);
      let left = this.getCell(cell.x - 2, cell.y);
  
      // If a cell is within the grid and has not been visited yet, add it to the list of neighbors
      if (top && !top.visited) neighbors.push(top);
      if (right && !right.visited) neighbors.push(right);
      if (bottom && !bottom.visited) neighbors.push(bottom);
      if (left && !left.visited) neighbors.push(left);
  
      if (neighbors.length > 0) {
        // 75% chance to return a random neighbor
        if (Math.random() < 0.75) {
          return neighbors[Math.floor(Math.random() * neighbors.length)];
        }
  
        // Otherwise, sort the neighbors based on their distance from the start point
        neighbors.sort((a, b) => {
          let aDistance = Math.hypot(a.x - this.start.x, a.y - this.start.y);
          let bDistance = Math.hypot(b.x - this.start.x, b.y - this.start.y);
          return bDistance - aDistance;
        });
  
        // Return the neighbor that is furthest from the start point
        return neighbors[0];
      } else {
        return undefined;
      }
    }
  
    generateMazeNotGlobal() {
      this.currentCell.visited = true;
      let nextCell = this.getNeighbors(this.currentCell);
  
      if (nextCell) {
        nextCell.visited = true;
  
        this.stack.push(this.currentCell);
  
        // Remove the wall between the current cell and the next cell
        let wallX = (this.currentCell.x + nextCell.x) / 2;
        let wallY = (this.currentCell.y + nextCell.y) / 2;
        this.grid[wallY][wallX].isWall = false;
  
        this.currentCell = nextCell;
      } else if (this.stack.length > 0) {
        // 40% chance to backtrack and carve a new path
        if (Math.random() < 0.4) {
          let backtrackCell = this.stack[Math.floor(Math.random() * this.stack.length)];
          this.currentCell = backtrackCell;
        } else {
          this.currentCell = this.stack.pop();
        }
      }
    }
  
  }
  
  function createNode(col, row, isWall, startNode, endNode, gridId) {
    return {
      col,
      row,
      isStart: startNode && col === startNode.x && row === startNode.y,
      isEnd: endNode && col === endNode.x && row === endNode.y,
      distance: Infinity,
      isVisited: false,
      isWall,
      previousNode: null,
      gridId
    };
  }