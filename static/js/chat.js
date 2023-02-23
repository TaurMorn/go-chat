let hiddenRoomNumber = document.getElementById("message-room-number");
let roomForm = document.getElementById("room-form");
let messageForm = document.getElementById("message-form");


window.addEventListener("DOMContentLoaded", (_) => {
    if (hiddenRoomNumber.value) {
        roomForm.hidden = true;
        messageForm.hidden = false;
    }
} );