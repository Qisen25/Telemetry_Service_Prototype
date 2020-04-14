const express = require('express');
const http = require('http');
const app = express();
app.set('port', process.env.PORT || 3000);
const server = http.createServer(app);
const io = require('socket.io').listen(server);

const redis = require('redis');
const redisClient = redis.createClient();

console.log('Node app running on port 3000');
server.listen(app.get('port'));

app.get('/', (req, res) => {
	//res.send('Server is running on port 3000');
	res.sendFile('MapPage.html', { root: __dirname});
});

io.on('connection', (socket) => {

	console.log("User has connected " + socket.id);
	socket.on('join', function(name) {
		console.log(name + " has connected " + socket.id);
		socket.emit('returnID', socket.id);
	});

	socket.on('send_location', function(JSONobj, name) {
		console.log("From " + JSONobj['ID'] + " " + name + ": " + JSONobj['latitude'] + ", " + JSONobj['longitude']);
		io.sockets.emit('updateLocation', name, JSONobj);
		redisClient.geoadd('publishers', JSONobj['longitude'], JSONobj['latitude'], JSONobj['ID'], function(er, reply){
			if(er)
			{
				console.log(er);
			}
		});
	});

	// socket.on('clean_up', function(uniqueID) {
	// 	console.log("Removing marker from map of ID " + uniqueID);
	// 	io.sockets.emit('removeMarker', uniqueID);
	// });

	socket.on('disconnect', function(reason) {
		console.log("user: "+ socket.id + " has disconnected " + reason);
		redisClient.zrem('publishers', socket.id);
		redisClient.sadd('offline', socket.id);
		io.sockets.emit('removeMarker', socket.id);
	});
});