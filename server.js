import express from 'express';
import { createServer } from 'http';
import WebSocket from 'isomorphic-ws';
import { exec } from 'child_process';
import fs from 'fs';
import path from 'path';

import { fileURLToPath } from 'url';
import { dirname } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const app = express();
app.get('/', (req, res) => {
    res.send('WebSocket server is running');
});
const server = createServer(app);
const wss = new WebSocket.Server({ server });

wss.on('connection', ws => {
    ws.on('message', message => {
        const { mazeSize, numTests } = JSON.parse(message);
        const testRunner = exec(`node testRunner.js ${mazeSize} ${numTests}`);

        testRunner.stdout.on('data', data => {
            ws.send(data); // Send the data to the client
        });

        testRunner.stderr.on('data', data => {
            console.error(`stderr: ${data}`);
        });

        testRunner.on('close', code => {
            // Read the CSV data
            const dataFileName = `averages${mazeSize}x${mazeSize}.csv`;
            const data = fs.readFileSync(path.join(__dirname, 'data', dataFileName), 'utf8');
            // Send the CSV data to the client
            ws.send(JSON.stringify({ data }));
        });
    });
});

server.listen(5000, () => console.log('Server is running on port 5000'));
