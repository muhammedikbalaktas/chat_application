var submitButton = document.getElementById("btn_submit");

submitButton.addEventListener("click", function() {
    var inputUsername = document.getElementById("username");
    var inputPassword = document.getElementById("password");
    const username = inputUsername.value;
    const password = inputPassword.value;
    if (username == "" || password == "") {
        alert("Some fields are empty");
    } else {
        var loading = document.getElementById("loader");
        loading.style.display = 'grid';
        sendRequest(username, password);
    }
});

function sendRequest(username, password) {
    const url = 'http://localhost:8080/get_user';
    const data = {
        "username": username,
        "password": password
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
            console.log('Success:', data);
            
            console.log(data.token);
            sessionStorage.setItem("token",data.token)
            window.location.href = "http://127.0.0.1:5500/homepage/index.html";
        })
        .catch(error => {
            console.error('Error:', error);
        });
}
