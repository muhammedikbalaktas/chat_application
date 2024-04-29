
var token = sessionStorage.getItem('token');
sendRequest(token)
function sendRequest(token) {
    const url = 'http://localhost:8080/get_contacts';
    const data = {
        "token": token
    };

    fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        .then(response => {
            if (!response.ok) {
                console.log(response);
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            
            const dataArray = Array.isArray(data) ? data : [data]; 
            parseContacts(dataArray)
            
        })
        .catch(error => {
            console.error('Error:', error);
        });
}

function parseContacts(contacts){
    var baseDiv=document.getElementById("users_div")
    
    contacts.forEach(user => {
        const cardDiv = document.createElement('div');
        cardDiv.classList.add("card_div");
        cardDiv.id="card_div"
        const img = document.createElement('img');
        img.src="../testImages/img1.jpeg"
        img.alt="img"
        img.id="pp"
        cardDiv.appendChild(img)
        const h1 = document.createElement('h1');
        h1.id="username"
        h1.textContent=user.username
        cardDiv.appendChild(h1)
        console.log(user.id);
        console.log(user.username);
        

        cardDiv.addEventListener("click", event=>{
            sessionStorage.setItem("reciever_id",user.id)
            getMessages(token,user.id)
        })
        baseDiv.appendChild(cardDiv)
    });
}



function getMessages(token, recieverId){

const url = 'http://localhost:8080/list_messages';
    const data = {
        "token": token,
        "reciever_id":recieverId
    };

    fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        .then(response => {
            if (!response.ok) {
                console.log(response);
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            
            const dataArray = Array.isArray(data) ? data : [data]; 
            parseMessages(dataArray,recieverId)
            
        })
        .catch(error => {
            console.error('Error:', error);
        });
}

function parseMessages(dataArray,recieverId){
    
    var rootDiv=document.getElementById("bubbles_div")
    dataArray.forEach(message => {
        if(recieverId!=message.reciever_id){
            //add bubble to left
            var leftBubble=document.createElement("div")
            leftBubble.classList.add("bubble_left")
            var p=document.createElement("p")
            p.textContent=message.text
            leftBubble.appendChild(p)
            rootDiv.appendChild(leftBubble)
        }else{
            //add bubble to right
            var leftBubble=document.createElement("div")
            leftBubble.classList.add("bubble_right")
            var p=document.createElement("p")
            p.textContent=message.text
            leftBubble.appendChild(p)
            rootDiv.appendChild(leftBubble)
        }

    });
    var objDiv = document.getElementById("bubbles_div");
    objDiv.scrollTop = objDiv.scrollHeight;
    
}   


