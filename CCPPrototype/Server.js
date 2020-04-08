const express = require('express');
const http = require('http');
const app = express();
app.set('port', process.env.PORT || 3000);
const server = http.createServer(app);
const io = require('socket.io').listen(server);

console.log('Node app running on port 3000');
server.listen(app.get('port'));

app.get('/', (req, res) => {
	//res.send('Server is running on port 3000');
	res.sendFile('MapPage.html', { root: __dirname});
});

io.on('connection', (socket) => {

	socket.on('join', function(name) {
		console.log(name + " has connected " + socket.id);
		socket.emit('returnID', socket.id);
	});

	socket.on('send_location', function(JSONobj, name) {
		console.log("From " + JSONobj['ID'] + name + ": " + JSONobj['lattitude'] + ", " + JSONobj['longitude']);
		io.sockets.emit('updateLocation', name, JSONobj);
	});

	socket.on('clean_up', function(uniqueID) {
		console.log("Removing marker from map of ID " + uniqueID);
		io.sockets.emit('removeMarker', uniqueID);
	});

	socket.on('disconnect', function(reason) {
		console.log("user has disconnected " + reason);
		io.sockets.emit('removeMarker', socket.id);
	});
});

