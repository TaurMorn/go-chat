let hiddenRoomNumber = document.getElementById("message-room-number");
let roomForm = document.getElementById("room-form");
let messageForm = document.getElementById("message-form");
let textArea = document.getElementById("text-input");
let nickName = document.getElementById("message-nick-name");
let websocket = null;

window.addEventListener("DOMContentLoaded", (_) => {
    if (hiddenRoomNumber.value) {
        let protocol = location.protocol === "http:" ? "ws:" : "wss:" ;
        websocket = new WebSocket(protocol + "//" + window.location.host + "/websocket");
        roomForm.hidden = true;
        messageForm.hidden = false;

        websocket.addEventListener("open", (event) => {
            sendMessage("Hello, blah-blah-blah", true);
        });
        websocket.addEventListener("message", (event) => {
            let msg = JSON.parse(event.data);
            console.log(msg);
            if (msg.Ping === true){
                setTimeout(() => {sendMessage("blip-blop", true);}, 10000);
            } else {
                console.log(msg);
            }
        });
        textArea.focus();
    }
} );

messageForm.addEventListener("submit", (event) => {
    event.preventDefault();
    sendMessage(textArea.value, false)
    textArea.value = "";
    textArea.focus();
});

messageForm.addEventListener("keypress", (event) => {
    if (event.key === "Enter") {
        event.preventDefault();
        sendMessage(textArea.value, false)
        textArea.value = "";
        textArea.focus();
    }
});

function sendMessage(text, ping) {
    websocket.send(
        JSON.stringify({
            RoomNumber: hiddenRoomNumber.value,
            UserName: nickName.value,
            Message: text,
            Ping: ping
        }));
}