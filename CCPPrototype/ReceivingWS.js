const http = require('http');
const WebSocket = require('ws');
const express = require('express');
const app = express();
app.set('port', process.env.PORT || 3000);
const server = http.createServer(app);

const readline = require("readline");
const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

// const server = http.createServer();

const wsServer = new WebSocket.Server({port: 3012});

app.get('/', (req, res) => {
    //res.send('Server is running on port 3000');
    res.sendFile('MapPageWS.html', { root: __dirname});
});


wsServer.on('connection', function connection(ws, request, client) {
	ws.on('message', function incoming(message) {
		ws.id = JSON.parse(message).id;
		console.log(ws.id + ' has joined');
	});
});

console.log('ReceivingWS listening on PORT 3000');

// rl.question("Enter to view list of clients", function(keyP) {
// 	if(keyP === '\r'){
// 		wsServer.clients.forEach(function each(client) {
// 			console.log('Client ' + client.id + ' exists');
// 		});
// 	}
// });

server.listen(app.get("port"));

process.stdin.on('keypress', (str, key) => {
  // console.log(str)
  // console.log(key)
    if(key.sequence === '\r') {
      	 wsServer.clients.forEach(function each(client) {
    		console.log('Client ' + client.id + ' exists');
    	});
    }
});