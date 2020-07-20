const WebSocket = require('ws');

const url = 'ws://localhost:3000';
const connection = new WebSocket(url);

connection.onopen = () => {
  var msg = {
    type: "message",
    id: 'TestID'
  }
  
  connection.send(JSON.stringify(msg)); 
}

connection.onerror = (error) => {
  console.log(`WebSocket error: ${error}`);
}

connection.onmessage = (e) => {
  console.log(e.data);
}