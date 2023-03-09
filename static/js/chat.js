let hiddenRoomNumber = document.getElementById("message-room-number");
let roomForm = document.getElementById("room-form");
let messageForm = document.getElementById("message-form");
let textArea = document.getElementById("text-input");
let nickName = document.getElementById("message-nick-name");
let websocket = null;
let chatDiv = document.getElementById("chat-text");

window.addEventListener("DOMContentLoaded", (_) => {
    if (hiddenRoomNumber.value) {
        let protocol = location.protocol === "http:" ? "ws:" : "wss:" ;
        websocket = new WebSocket(protocol + "//" + window.location.host + "/websocket");
        roomForm.hidden = true;
        messageForm.hidden = false;
        chatDiv.hidden = false;
        
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
                //let msgArray = msg.Message.match(/.{1,42}/g);
                /*
                for (let i = 0; i < msgArray.length; i++){
                    divInner.append(msgArray[i]);
                }
                */
                divInner.append (msg.Message);
                let currentTime = new Date().toLocaleTimeString();
                
                if (msg.UserName === nickName.value) {
                    strong.innerHTML = `<b>Me, </b><i>${currentTime}</i><br>`;
                    divInner.className = "text-black p-2 rounded-8";
                    divInner.style = "background-color:#DCEDC8; max-width: 500px";
                    divMsg.className = "d-flex flex-row-reverse";
                }else {
                    strong.innerHTML = `<b>${msg.UserName}, </b><i>${currentTime}</i><br>`;
                    divInner.className = "text-black p-2 rounded-8 bg-light";
                    divInner.style = "max-width: 500px";
                    divMsg.className = "d-flex flex-row";   
                }
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