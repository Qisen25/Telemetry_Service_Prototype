const WebSocket = require('ws');

const url = 'ws://localhost:8020';
const connection = new WebSocket(url);

var msg = {
type: "message",
id: 'TestID'
}

connection.onopen = () => {  
  connection.send(JSON.stringify(msg));
  connection.isOpen = 'true';
  // console.log("Sending data " + d.toLocaleString());
}

connection.onerror = (error) => {
  console.log(`WebSocket error: ${error}`);
}

connection.onclose = () => {
  console.log('WebSocket closed');
}

connection.onmessage = (e) => {
  console.log(e.data);
}


