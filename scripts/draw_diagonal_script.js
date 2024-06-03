const BASE_URL = "http://localhost:17000";
const figureWidthNorm = (300 / 800).toFixed(2);
const figureHeightNorm = (280 / 800).toFixed(2);
const offsetX = (figureWidthNorm / 2).toFixed(2);
const offsetY = (figureHeightNorm / 2).toFixed(2);

let Xstart = offsetX;
let Ystart = offsetY;
const step = (20 / 800).toFixed(2);
let direction = "down_right";

function sendPostRequest(data) {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", BASE_URL, true);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200) {
            console.log(`Request with data '${data}' sent successfully.`);
        }
    };
    xhr.send(data);
}

function updatePosition() {
    let newX, newY;

    if (direction === "down_right") {
        newX = (parseFloat(Xstart) + parseFloat(step)).toFixed(2);
        newY = (parseFloat(Ystart) + parseFloat(step)).toFixed(2);
        if (newX > 1 - offsetX || newY > 1 - offsetY) {
            direction = "up_left";
            newX = (1 - offsetX).toFixed(2);
            newY = (1 - offsetY).toFixed(2);
        }
    } else {
        newX = (parseFloat(Xstart) - parseFloat(step)).toFixed(2);
        newY = (parseFloat(Ystart) - parseFloat(step)).toFixed(2);
        if (newX < offsetX || newY < offsetY) {
            direction = "down_right";
            newX = offsetX;
            newY = offsetY;
        }
    }

    Xstart = newX;
    Ystart = newY;

    sendPostRequest("reset");
    sendPostRequest("green");
    sendPostRequest(`figure ${Xstart} ${Ystart}`);
    sendPostRequest("update");
}

setInterval(updatePosition, 1000);
