document.addEventListener('DOMContentLoaded', () => {
    const generateButton = document.getElementById('generate-button');
    const solveButton = document.getElementById('solve-button');
    const singlePathCheckbox = document.getElementById('singlePath');
    const mazeSizeSlider = document.getElementById('mazeSize');
    const mazeSizeValue = document.getElementById('mazeSizeValue');
    const gridContainer = document.getElementById('grid-container');
    const metricsTableBody = document.getElementById('metrics-table').querySelector('tbody');

    mazeSizeSlider.addEventListener('input', () => {
        mazeSizeValue.textContent = mazeSizeSlider.value;
    });

    generateButton.addEventListener('click', async () => {
        console.log('Generate button clicked');
        const mazeSize = mazeSizeSlider.value;
        const singlePath = singlePathCheckbox.checked;
        const response = await fetch(`/api/maze?mazeSize=${mazeSize}&singlePath=${singlePath}`);
        const data = await response.json();
        console.log('Maze data received:', data);
        renderGrids(data.grids);
        renderMetrics({}); // Reset metrics table
    });

    solveButton.addEventListener('click', async () => {
        console.log('Solve button clicked');
        const response = await fetch('/api/solution');
        const data = await response.json();
        console.log('Solution data received:', data);
        renderMetrics(data);
        for (const algorithm in data) {
            await visualizeAlgorithm(algorithm, data[algorithm]);
        }
    });

    function addNodeIds(grid, algorithmName) {
        return grid.map((row, y) =>
            row.map((node, x) => ({
                ...node,
                id: `grid${algorithmName.charAt(0).toUpperCase() + algorithmName.slice(1)}-node-${y}-${x}`,
                x: x,
                y: y,
            }))
        );
    }

    function renderGrids(grids) {
        console.log('Rendering grids');
        gridContainer.innerHTML = '';
        for (const [algorithm, grid] of Object.entries(grids)) {
            const gridElement = document.createElement('div');
            gridElement.classList.add('grid');
            gridElement.dataset.algorithm = algorithm;
            gridElement.style.gridTemplateColumns = `repeat(${grid[0].length}, 30px)`;

            const gridWithIds = addNodeIds(grid, algorithm);

            gridWithIds.forEach((row) => {
                row.forEach((node) => {
                    const nodeElement = document.createElement('div');
                    nodeElement.classList.add('node');
                    nodeElement.id = node.id;
                    nodeElement.dataset.x = node.x;
                    nodeElement.dataset.y = node.y;
                    if (node.isStart) nodeElement.classList.add('node-start');
                    if (node.isEnd) nodeElement.classList.add('node-end');
                    if (node.isWall) nodeElement.classList.add('node-wall');
                    gridElement.appendChild(nodeElement);
                });
            });
            gridContainer.appendChild(gridElement);
        }
    }

    function renderMetrics(data) {
        console.log('Rendering metrics:', data);
        metricsTableBody.innerHTML = '';
        for (const [algorithm, result] of Object.entries(data)) {
            if (result && result.metrics) {
                const metric = result.metrics;
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${algorithm}</td>
                    <td>${metric.time[metric.time.length - 1].toFixed(2)}</td>
                    <td>${metric.visitedNodes[metric.visitedNodes.length - 1]}</td>
                    <td>${metric.visitedPercentage[metric.visitedPercentage.length - 1].toFixed(2)}%</td>
                    <td>${metric.pathLength[metric.pathLength.length - 1]}</td>
                    <td>${metric.memoryUsed[metric.memoryUsed.length - 1].toFixed(2)}</td>
                `;
                metricsTableBody.appendChild(row);
            }
        }
    }

    function animateShortestPath(nodesInShortestPathOrder) {
        console.log('Animating shortest path');
        return new Promise(resolve => {
            let i = 0;
            const animate = () => {
                if (i < nodesInShortestPathOrder.length) {
                    const node = nodesInShortestPathOrder[i];
                    const nodeElement = document.getElementById(node.id);
                    if (nodeElement) {
                        nodeElement.className = 'node node-shortest-path';
                    } else {
                        console.warn(`Node element not found: ${node.id}`);
                    }
                    i++;
                    requestAnimationFrame(animate);
                } else {
                    resolve();
                }
            };
            animate();
        });
    }

    function animateAlgorithm(visitedNodesInOrder, nodesInShortestPathOrder) {
        console.log('Animating algorithm', visitedNodesInOrder.length, nodesInShortestPathOrder.length);
        return new Promise(resolve => {
            let i = 0;
            const animate = () => {
                if (i < visitedNodesInOrder.length) {
                    const node = visitedNodesInOrder[i];
                    const nodeElement = document.getElementById(node.id);
                    if (nodeElement) {
                        nodeElement.classList.add('node', 'node-visited');
                        if (node.noOfVisits > 1 && visitedNodesInOrder[i + 1] === node.previousNode) {
                            const hue = 174 + (node.noOfVisits - 1) * 10;
                            const lightness = 30 - (node.noOfVisits - 1) * 5;
                            nodeElement.style.backgroundColor = `hsl(${hue}, 50%, ${lightness}%)`;
                        }
                    } else {
                        console.warn(`Node element not found: ${node.id}`);
                    }
                    i++;
                    requestAnimationFrame(animate);
                } else {
                    animateShortestPath(nodesInShortestPathOrder).then(resolve);
                }
            };
            animate();
        });
    }

    async function visualizeAlgorithm(algorithmName, data) {
        console.log(`Visualizing algorithm: ${algorithmName}`, data);
        const visitedNodesInOrder = data.visitedNodesInOrder;
        const nodesInShortestPathOrder = data.nodesInShortestPathOrder;

        // Update node IDs to match the current grid
        visitedNodesInOrder.forEach(node => {
            node.id = `grid${algorithmName.charAt(0).toUpperCase() + algorithmName.slice(1)}-node-${node.row}-${node.col}`;
        });
        nodesInShortestPathOrder.forEach(node => {
            node.id = `grid${algorithmName.charAt(0).toUpperCase() + algorithmName.slice(1)}-node-${node.row}-${node.col}`;
        });

        await animateAlgorithm(visitedNodesInOrder, nodesInShortestPathOrder);
    }
}););););
