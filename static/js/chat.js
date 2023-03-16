let hiddenRoomNumber = document.getElementById("message-room-number");
let roomForm = document.getElementById("room-form");
let messageForm = document.getElementById("message-form");
let textArea = document.getElementById("text-input");
let nickName = document.getElementById("message-nick-name");
let websocket = null;
let chatDiv = document.getElementById("chat-text");
let cooldynamicDiv = document.getElementById("dynamicDiv");

window.addEventListener("DOMContentLoaded", (_) => {
    if (hiddenRoomNumber.value) {
        let protocol = location.protocol === "http:" ? "ws:" : "wss:" ;
        websocket = new WebSocket(protocol + "//" + window.location.host + "/websocket");
        roomForm.hidden = true;
        messageForm.hidden = false;
        chatDiv.hidden = false;
        cooldynamicDiv.className = "p-3 mb-0 text-black text-center";
        
        websocket.addEventListener("open", (event) => {
            sendMessage("Hello, blah-blah-blah", true);
        });
        websocket.addEventListener("message", (event) => {
            let msg = JSON.parse(event.data);
            console.log(msg);
            if (msg.Ping === true){
                setTimeout(() => {sendMessage("blip-blop", true);}, 10000);
            } else {
                let divMsg = document.createElement("div");
                let divInner = document.createElement("div");
                let strong = document.createElement("strong");
                
                strong.className = "text-primary";
                divInner.append(strong);
                
                divInner.append (msg.Message);

                let currentTime = new Date().toLocaleTimeString();
                
                if (msg.UserName === nickName.value) {
                    strong.innerHTML = `<b style="color: #03453d;">Me, </b><i style="color: #03453d;">${currentTime}</i><br>`;
                    divInner.className = "text-black p-2 rounded-8";
                    divInner.style = "max-width: 600px; background-color:#84d9cf; word-break:break-all;";
                    divMsg.className = "d-flex flex-row-reverse";
                }else {
                    strong.innerHTML = `<b style="color: #03453d;">${msg.UserName}, </b><i style="color: #03453d;">${currentTime}</i><br>`;
                    divInner.className = "text-black p-2 rounded-8 bg-light";
                    divInner.style = "word-break:break-all; max-width:600px;";
                    divMsg.className = "d-flex flex-row";   
                }
                divMsg.style = "margin-bottom: 9px;";
                divMsg.append(divInner);
                chatDiv.append(divMsg);
                chatDiv.scrollTop = chatDiv.scrollHeight;
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