let selectedChat = 'general';
let wsCnn;

function changeChatroom() {
    let newChat = document.getElementById('chatroom');
    if (newChat && newChat.value != selectedChat) {
        console.log(newChat);
    }

    return false;
}

function sendMessage() {
    let newMessage = document.getElementById('message');
    if (newMessage) {
        wsCnn.send(newMessage.value);
        newMessage.value = '';
    }

    return false;
}

document.addEventListener('DOMContentLoaded', () => {
    console.log('frontend OK');
    document.getElementById('chatroom-selection').onsubmit = changeChatroom;
    document.getElementById('chatroom-message').onsubmit = sendMessage;

    if (window['WebSocket']) {
        console.log('Websocket suppport OK');
        wsCnn = new WebSocket(`ws://${document.location.host}/ws`);
    }
    else
        console.warn('Does NOT support Websocket');
})