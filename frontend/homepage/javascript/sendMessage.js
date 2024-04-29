var authToken = sessionStorage.getItem('token');
console.log(authToken);

var sendIcon=document.getElementById("send_icon");
sendIcon.addEventListener("click",sendMessage);
document.getElementById("text_input").addEventListener("keypress", function(event) {
    
    if (event.keyCode === 13) {
      
      event.preventDefault();
      
      sendMessage()
    }
  });
function sendMessage() {
    var inputField = document.getElementById("text_input");
    var inputString = inputField.value;

    if (inputString === "") {
        alert("Please write something");
        return;
    }

    var divElement = document.createElement('div');
    divElement.setAttribute('class', 'bubble_right'); // Use a class instead of an ID
    var paragraph = document.createElement('p');
    paragraph.textContent = inputString;
    divElement.appendChild(paragraph);

    var rootDiv = document.getElementById("bubbles_div");
    rootDiv.appendChild(divElement);
    
    

    var objDiv = document.getElementById("bubbles_div");
    objDiv.scrollTop = objDiv.scrollHeight;
    var recieverId = sessionStorage.getItem('reciever_id');
    sendMessageToServer(inputString,recieverId)
    inputField.value = '';

}


function recieveMessage(input){
    var divElement = document.createElement('div');
    divElement.setAttribute('class', 'bubble_left'); // Use a class instead of an ID

    var paragraph = document.createElement('p');
    paragraph.textContent = input;

    divElement.appendChild(paragraph);

    var rootDiv = document.getElementById("bubbles_div");
    rootDiv.appendChild(divElement);

}

const websocketURL = 'ws://localhost:8080/?token=' + encodeURIComponent(authToken);


const socket = new WebSocket(websocketURL);


socket.addEventListener('open', function (event) {
    console.log('WebSocket connection established');
});


socket.addEventListener('message', function (event) {
    
    
    parseSingleMessage(event.data)
});


socket.addEventListener('error', function (event) {
    console.error('WebSocket error:', event);
});


socket.addEventListener('close', function (event) {
    console.log('WebSocket connection closed');
});


function sendMessageToServer(message, recieverId) {
    var objDiv = document.getElementById("bubbles_div");
    objDiv.scrollTop = objDiv.scrollHeight;
    const messageObj = {
        "text": message,
        "reciever_id": parseInt(recieverId)
    };


    const messageStr = JSON.stringify(messageObj);


    socket.send(messageStr);

    
}
function parseSingleMessage(data){
    
    var rootDiv=document.getElementById("bubbles_div")
    //add bubble to left
    // Message from server: {"text":"merhaba","reciever_id":1}
    const jsonObject = JSON.parse(data);
    var leftBubble=document.createElement("div")
    leftBubble.classList.add("bubble_left")
    var p=document.createElement("p")
    p.textContent=jsonObject.text
    leftBubble.appendChild(p)
    rootDiv.appendChild(leftBubble)
    var objDiv = document.getElementById("bubbles_div");
    objDiv.scrollTop = objDiv.scrollHeight;
    
    
} 

