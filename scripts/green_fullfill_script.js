function sendPostRequest(data) {
    const xhr = new XMLHttpRequest();
    xhr.open("POST", "http://localhost:17000", true);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
            console.log(`Request with data '${data}' sent successfully.`);
        }
    };
    xhr.send(data);
}

function sendWhiteCommand() {
    sendPostRequest("white");
}

function sendGreenCommand() {
    sendPostRequest("green");
}

function sendUpdateCommand() {
    sendPostRequest("update");
}

// Send the commands
sendWhiteCommand();
sendGreenCommand();
sendUpdateCommand();
