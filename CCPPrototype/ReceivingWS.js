const https = require('https');
const http = require('http');
const fs = require('fs');
const WebSocket = require('ws');
const express = require('express');
const app = express();
app.set('port', process.env.PORT || 3000);

//https cert and keys
const options = {
    key: fs.readFileSync('server-key.pem'), 
    cert: fs.readFileSync('server-crt.pem'), 
    ca: fs.readFileSync('ca-crt.pem'), 
};


const server = http.createServer(app);

// const server = https.createServer(options, app);

// const readline = require("readline");
// const rl = readline.createInterface({
//     input: process.stdin,
//     output: process.stdout
// });

// const server = http.createServer();

const wsServer = new WebSocket.Server({server});

app.get('/', (req, res) => {
    //res.send('Server is running on port 3000');
    res.sendFile('MapPageWS.html', { root: __dirname});
});

//when client establishes client to server
wsServer.on('connection', function connection(ws, request, client) {
	// ws.on('message', function incoming(message) {
	// 	ws.id = JSON.parse(message).id;
	// 	console.log(ws.id + ' has joined');
	// });
    ws.id = request.socket.remoteAddress;
    console.log(ws.id + ' has joined');

    ws.on('close', function close() {//when client socket closes
      console.log(ws.id + ' has disconnected');
    });
});

// rl.question("Enter to view list of clients", function(keyP) {
// 	if(keyP === '\r'){
// 		wsServer.clients.forEach(function each(client) {
// 			console.log('Client ' + client.id + ' exists');
// 		});
// 	}
// });

//http
server.listen(app.get("port"), () => {
    console.log('ReceivingWS http listening on PORT ' + app.get("port"));
});
//https
// server.listen(3001, () => {
//     console.log('ReceivingWS https listening on PORT 3001');
// });

// process.stdin.on('keypress', (str, key) => {
//   // console.log(str)
//   // console.log(key)
//     if(key.sequence === '\r') {
//       	 wsServer.clients.forEach(function each(client) {
//     		console.log('Client ' + client.id + ' exists');
//     	});
//     }
// });