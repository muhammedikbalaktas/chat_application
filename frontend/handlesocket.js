const websocketURL = 'ws://localhost:8080/';
const authToken = sessionStorage.getItem('token');

// Create WebSocket object with authorization token in headers
const socket = new WebSocket(websocketURL, ['Authorization', authToken]);

// Event listener for when the WebSocket connection is opened
socket.addEventListener('open', function (event) {
    console.log('WebSocket connection established');
});

// Event listener for when the WebSocket connection receives a message
socket.addEventListener('message', function (event) {
    console.log('Message from server:', event.data);
});

// Event listener for when an error occurs with the WebSocket connection
socket.addEventListener('error', function (event) {
    console.error('WebSocket error:', event);
});

// Event listener for when the WebSocket connection is closed
socket.addEventListener('close', function (event) {
    console.log('WebSocket connection closed');
});

// Function to send a message using WebSocket
function sendMessage() {
    const message = 'Hello, WebSocket Server!';
    socket.send(message);
}
