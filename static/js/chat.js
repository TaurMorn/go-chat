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
                Message: "Hello, blah-blah-blah", 
                Ping: true
              });
            websocket.send(msg);
        });
        websocket.addEventListener("message", (event) => {
            let msg = JSON.parse(event.data);
            console.log(msg);
            if (msg.Ping === true){
                setTimeout(() => websocket.send(
                    JSON.stringify({
                        RoomNumber: hiddenRoomNumber.value,
                        UserName: nickName.value,
                        Message: "blip-blop",
                        Ping: true
                    })
                ), 10000);
            return;
            }
        });

        textArea.focus();
    }
} );