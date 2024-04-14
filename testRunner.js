import path from 'path';
import fs from 'fs';
import { Worker, isMainThread, parentPort, workerData } from 'worker_threads';
import { generateMaze } from './Algorithms/mazeGenerator.js';
import { fileURLToPath } from 'url';
import { dirname } from 'path';
import { dijkstra } from './Algorithms/dijkstra.js';
import { astar } from './Algorithms/astar.js';
import { bfs } from './Algorithms/bfs.js';
import { dfs } from './Algorithms/dfs.js';
import { wallFollower } from './Algorithms/wall_follower.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const algorithms = {
    dijkstra,
    astar,
    bfs,
    dfs,
    wallFollower
};

if (isMainThread) {
    // This code will be executed in the main thread
    function runAlgorithmInWorker(workerData) {
        return new Promise((resolve, reject) => {
            const worker = new Worker(__filename, { workerData });

            worker.on('message', resolve);
            worker.on('error', reject);
            worker.on('exit', (code) => {
                if (code !== 0) reject(new Error(`Worker stopped with exit code ${code}`));
            });
        });
    }

    async function runTest(numMazes) {
        const numRows = 750;
        const numCols = 750;
        const metricsSPOn = {
            dijkstra: { time: [], visitedNodes: [], visitedPercentage: [], pathLength: [], memoryUsed: [] },
            astar: { time: [], visitedNodes: [], visitedPercentage: [], pathLength: [], memoryUsed: [] },
            bfs: { time: [], visitedNodes: [], visitedPercentage: [], pathLength: [], memoryUsed: [] },
            dfs: { time: [], visitedNodes: [], visitedPercentage: [], pathLength: [], memoryUsed: [] },
            wallFollower: { time: [], visitedNodes: [], visitedPercentage: [], pathLength: [], memoryUsed: [] }
        };
        const metricsSPOff = {
            dijkstra: { time: [], visitedNodes: [], visitedPercentage: [], pathLength: [], memoryUsed: [] },
            astar: { time: [], visitedNodes: [], visitedPercentage: [], pathLength: [], memoryUsed: [] },
            bfs: { time: [], visitedNodes: [], visitedPercentage: [], pathLength: [], memoryUsed: [] },
            dfs: { time: [], visitedNodes: [], visitedPercentage: [], pathLength: [], memoryUsed: [] },
            wallFollower: { time: [], visitedNodes: [], visitedPercentage: [], pathLength: [], memoryUsed: [] }
        };

        // Test mazes with a single path
        for (let i = 0; i < numMazes; i++) {
            const { griddijkstra, gridastar, gridbfs, griddfs, gridwallFollower, griddijkstraStartNode, griddijkstraEndNode, gridastarStartNode, gridastarEndNode, gridbfsStartNode, gridbfsEndNode, griddfsStartNode, griddfsEndNode, gridwallFollowerStartNode, gridwallFollowerEndNode } = getInitialGrid(numRows, numCols, true);

            const algorithms = ['astar', 'bfs', 'dfs', 'wallFollower'];
            const promises = [];

            const grids = {
                //dijkstra: griddijkstra,
                astar: gridastar,
                bfs: gridbfs,
                dfs: griddfs,
                wallFollower: gridwallFollower
            };

            const startNodes = {
                //dijkstra: griddijkstraStartNode,
                astar: gridastarStartNode,
                bfs: gridbfsStartNode,
                dfs: griddfsStartNode,
                wallFollower: gridwallFollowerStartNode
            };

            const endNodes = {
                //dijkstra: griddijkstraEndNode,
                astar: gridastarEndNode,
                bfs: gridbfsEndNode,
                dfs: griddfsEndNode,
                wallFollower: gridwallFollowerEndNode
            };

            algorithms.forEach(algorithm => {
                const promise = runAlgorithmInWorker({
                    algorithm,
                    grid: grids[algorithm],
                    startNode: startNodes[algorithm],
                    endNode: endNodes[algorithm],
                }).then(({ algorithm, time, visitedNodesInOrder, nodesInShortestPathOrder, memoryUsed, visitedPercentage }) => {
                    const visitedNodes = visitedNodesInOrder.length;
                    const pathLength = nodesInShortestPathOrder.length;
                    metricsSPOn[algorithm].time.push(time);
                    metricsSPOn[algorithm].visitedNodes.push(visitedNodes);
                    metricsSPOn[algorithm].visitedPercentage.push(visitedPercentage);
                    metricsSPOn[algorithm].pathLength.push(pathLength);
                    metricsSPOn[algorithm].memoryUsed.push(memoryUsed);
                });
                promises.push(promise);
            });

            await Promise.all(promises);
            console.log(`Completed test ${i + 1} of ${numMazes} for mazes with a single path`);
        }

        // Test mazes with multiple paths
        for (let i = 0; i < numMazes; i++) {
            const { griddijkstra, gridastar, gridbfs, griddfs, gridwallFollower, griddijkstraStartNode, griddijkstraEndNode, gridastarStartNode, gridastarEndNode, gridbfsStartNode, gridbfsEndNode, griddfsStartNode, griddfsEndNode, gridwallFollowerStartNode, gridwallFollowerEndNode } = getInitialGrid(numRows, numCols, false);

            const algorithms = ['astar', 'bfs', 'dfs', 'wallFollower'];
            const promises = [];

            const grids = {
                //dijkstra: griddijkstra,
                astar: gridastar,
                bfs: gridbfs,
                dfs: griddfs,
                wallFollower: gridwallFollower
            };

            const startNodes = {
                //dijkstra: griddijkstraStartNode,
                astar: gridastarStartNode,
                bfs: gridbfsStartNode,
                dfs: griddfsStartNode,
                wallFollower: gridwallFollowerStartNode
            };

            const endNodes = {
                //dijkstra: griddijkstraEndNode,
                astar: gridastarEndNode,
                bfs: gridbfsEndNode,
                dfs: griddfsEndNode,
                wallFollower: gridwallFollowerEndNode
            };

            algorithms.forEach(algorithm => {
                const promise = runAlgorithmInWorker({
                    algorithm,
                    grid: grids[algorithm],
                    startNode: startNodes[algorithm],
                    endNode: endNodes[algorithm],
                }).then(({ algorithm, time, visitedNodesInOrder, nodesInShortestPathOrder, memoryUsed, visitedPercentage }) => {
                    const visitedNodes = visitedNodesInOrder.length;
                    const pathLength = nodesInShortestPathOrder.length;
                    metricsSPOff[algorithm].time.push(time);
                    metricsSPOff[algorithm].visitedNodes.push(visitedNodes);
                    metricsSPOff[algorithm].visitedPercentage.push(visitedPercentage);
                    metricsSPOff[algorithm].pathLength.push(pathLength);
                    metricsSPOff[algorithm].memoryUsed.push(memoryUsed);
                });
                promises.push(promise);
            });

            await Promise.all(promises);
            console.log(`Completed test ${i + 1} of ${numMazes} for mazes with multiple paths`);

        }
        // Calculate averages and log results
        const averagesSPOn = calculateAverages(metricsSPOn);
        const averagesSPOff = calculateAverages(metricsSPOff);

        // Write results to a .csv file
        writeResultsToCsv('./averages.csv', averagesSPOn, averagesSPOff);
    }

    function calculateAverages(metrics) {
        const averages = {
            dijkstra: { time: 0, visitedNodes: 0, visitedPercentage: 0, pathLength: 0, memoryUsed: 0 },
            astar: { time: 0, visitedNodes: 0, visitedPercentage: 0, pathLength: 0, memoryUsed: 0 },
            bfs: { time: 0, visitedNodes: 0, visitedPercentage: 0, pathLength: 0, memoryUsed: 0 },
            dfs: { time: 0, visitedNodes: 0, visitedPercentage: 0, pathLength: 0, memoryUsed: 0 },
            wallFollower: { time: 0, visitedNodes: 0, visitedPercentage: 0, pathLength: 0, memoryUsed: 0 }
        };


        for (const algorithm in metrics) {
            const numTests = metrics[algorithm].time.length;
            for (const metric in metrics[algorithm]) {
                const sum = metrics[algorithm][metric].reduce((a, b) => a + b, 0);
                let average = sum / numTests;
                if (metric === 'time') {
                    average = average.toFixed(4); // round to 4 decimal places
                } else if (metric === 'visitedPercentage') {
                    average = (average).toFixed(2); // convert to percentage and round to 2 decimal places
                } else if (metric === 'memoryUsed') {
                    average = (average / 1024 / 1024).toFixed(2); // convert to MB and round to 2 decimal places
                }
                averages[algorithm][metric] = average;
            }
        }

        return averages;
    }

    function writeResultsToCsv(filename, averagesSPOn, averagesSPOff) {
        const header = ['Algorithm', 'SinglePath', 'Time', 'VisitedNodes', 'VisitedPercentage', 'PathLength', 'MemoryUsed'];
        const rows = [];

        for (const algorithm in averagesSPOn) {
            const row = [
                algorithm,
                'true',
                averagesSPOn[algorithm].time,
                averagesSPOn[algorithm].visitedNodes,
                averagesSPOn[algorithm].visitedPercentage,
                averagesSPOn[algorithm].pathLength,
                averagesSPOn[algorithm].memoryUsed
            ];
            rows.push(row.join(','));
        }

        for (const algorithm in averagesSPOff) {
            const row = [
                algorithm,
                'false',
                averagesSPOff[algorithm].time,
                averagesSPOff[algorithm].visitedNodes,
                averagesSPOff[algorithm].visitedPercentage,
                averagesSPOff[algorithm].pathLength,
                averagesSPOff[algorithm].memoryUsed
            ];
            rows.push(row.join(','));
        }

        const csv = [header.join(','), ...rows].join('\n');

        // Write to the CSV file
        fs.writeFile(filename, csv, (err) => {
            if (err) {
                console.error('Error writing to CSV file', err);
            } else {
                console.log('Successfully wrote to CSV file');
            }
        });
    }

    const getInitialGrid = (numOfRows, numOfCols, singlePath) => {
        const { gridDijsktra, gridAstar, gridBFS, gridDFS, gridWallFollower, gridDijkstraStartNode, gridDijkstraEndNode, gridAstarStartNode, gridAstarEndNode, gridBFSStartNode, gridBFSEndNode, gridDFSStartNode, gridDFSEndNode, gridWallFollowerStartNode, gridWallFollowerEndNode } = generateMaze(numOfRows, numOfCols, singlePath);
        return {
            griddijkstra: gridDijsktra,
            gridastar: gridAstar,
            gridbfs: gridBFS,
            griddfs: gridDFS,
            gridwallFollower: gridWallFollower,
            griddijkstraStartNode: gridDijkstraStartNode,
            griddijkstraEndNode: gridDijkstraEndNode,
            gridastarStartNode: gridAstarStartNode,
            gridastarEndNode: gridAstarEndNode,
            gridbfsStartNode: gridBFSStartNode,
            gridbfsEndNode: gridBFSEndNode,
            griddfsStartNode: gridDFSStartNode,
            griddfsEndNode: gridDFSEndNode,
            gridwallFollowerStartNode: gridWallFollowerStartNode,
            gridwallFollowerEndNode: gridWallFollowerEndNode
        };
    };

    runTest(5)
        .catch(error => {
            console.error(error);
        });
} else {
    // This code will be executed in the worker thread
    const { algorithm, grid, startNode, endNode } = workerData;
    const startTime = performance.now();
    const initialMemoryUsage = process.memoryUsage().heapUsed;
    const visitedNodesInOrder = algorithms[algorithm](grid, startNode, endNode);
    const finalMemoryUsage = process.memoryUsage().heapUsed;
    const endTime = performance.now();
    const nodesInShortestPathOrder = getNodesInShortestPathOrder(endNode);
    const time = endTime - startTime;
    const memoryUsed = finalMemoryUsage - initialMemoryUsage;

    // Calculate the total number of visitable nodes
    const totalNodes = grid.length * grid[0].length;
    const wallNodes = grid.flat().filter(node => node.isWall).length;
    const nonWallNodes = totalNodes - wallNodes;

    // Calculate the visited percentage based on the number of visitable nodes
    const visitedPercentage = (visitedNodesInOrder.length / nonWallNodes) * 100;

    parentPort.postMessage({ algorithm, time, visitedNodesInOrder, nodesInShortestPathOrder, memoryUsed, visitedPercentage });
}


function getNodesInShortestPathOrder(endNode) {
    const nodesInShortestPathOrder = [];
    let currentNode = endNode;
    while (currentNode !== null && currentNode !== undefined) {
        nodesInShortestPathOrder.unshift(currentNode);
        currentNode = currentNode.previousNode;
    }
    return nodesInShortestPathOrder;
}