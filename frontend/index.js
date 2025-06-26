let selectedChat = 'general';
let wsCnn;

class Event {
    constructor(type, payload) {
        this.type = type;
        this.payload = payload;
    }
}

function routeEvent(ev) {
    if (ev.type == undefined) {
        alert('no "type" field in event!');
        return;
    }

    switch (ev.type) {
        case 'new_message':
            console.log('new message');
            break;
        default:
            alert('unsupported message type');
            break;
    }
}

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
        // wsCnn.send(newMessage.value);
        sendEvent('send_message', newMessage.value);
        newMessage.value = '';
    }

    return false;
}

function sendEvent(evName, payload){
    const ev = new Event(evName, payload);
    wsCnn.send(JSON.stringify(ev));
}

document.addEventListener('DOMContentLoaded', () => {
    console.log('frontend OK');
    document.getElementById('chatroom-selection').onsubmit = changeChatroom;
    document.getElementById('chatroom-message').onsubmit = sendMessage;

    if (window['WebSocket']) {
        console.log('Websocket suppport OK');
        wsCnn = new WebSocket(`ws://${document.location.host}/ws`);

        //Add a listener to the onmessage event
        wsCnn.onmessage = function (ev) {
            console.log(ev);

            const eventData = JSON.parse(ev.data);
            const event = Object.assign(new Event, eventData);

            routeEvent(event);
        }
    }
    else
        console.warn('Does NOT support Websocket');
})