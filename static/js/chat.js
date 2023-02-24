let hiddenRoomNumber = document.getElementById("message-room-number");
let roomForm = document.getElementById("room-form");
let messageForm = document.getElementById("message-form");
let textArea = document.getElementById ("text-input");
let websocket = null;

window.addEventListener("DOMContentLoaded", (_) => {
    if (hiddenRoomNumber.value) {
        let protocol = location.protocol === "http:" ? "ws:" : "wss:" ;
        websocket = new WebSocket(protocol + "//" + window.location.host + "/websocket");
        roomForm.hidden = true;
        messageForm.hidden = false;

        websocket.addEventListener("open", (event) => {
            let msg =  JSON.stringify({
                RoomNumber: hiddenRoomNumber.value,
                UserName: nickName.value,
                Message: "Hello, blah-blah-blah"
              });
            websocket.send(msg);
        });

        textArea.focus();
    }
} );