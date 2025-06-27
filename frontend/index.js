let selectedChat = 'general';
let wsCnn;

let ui = {};

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

function sendEvent(evName, payload) {
    const ev = new Event(evName, payload);
    wsCnn.send(JSON.stringify(ev));
}

document.addEventListener('DOMContentLoaded', () => {
    console.log('frontend OK');

    ui = {
        connectionHeader: document.getElementById('connection-header'),
        chatroomSelection: document.getElementById('chatroom-selection'),
        chatroomMessage: document.getElementById('chatroom-message'),
    };

    ui.chatroomSelection.onsubmit = changeChatroom;
    ui.chatroomMessage.onsubmit = sendMessage;
})

function login() {
    let frm = new FormData(document.getElementById('login-form'));
    let payload = {};
    // let payload = { user: frm.get('username'), password: frm.get('password') };

    for (const [key, value] of frm) {
        payload[key] = value;
    }

    console.log(payload);

    fetch('login', {
        method: 'post',
        body: JSON.stringify(payload),
        mode: 'cors',
    }).then((res) => {
        if (res.ok)
            return res.json();

        throw 'unauthorized';
    }).then((data) => {
        //Now we hava a OTP, send a Request to Connect to Websocket
        connectWebsocket(data.otp);
    }).catch((ex) => {
        alert(ex);
    });

    return false;
}

function connectWebsocket(otp) {
    if (window['WebSocket']) {
        console.log('👌 Websockets supported');

        //Connect to websocket using OTP as a GET parameter
        wsCnn = new WebSocket(`ws://${document.location.host}/ws?otp=${otp}`);

        //onOpen
        wsCnn.onopen = (ev) => {
            ui.connectionHeader.innerHTML = '😁 Connected to Websocket: true';
        };

        wsCnn.onclose = (ev) => {
            ui.connectionHeader.innerHTML = '😢 Connected to Websocket: FALSE';
        };

        wsCnn.onmessage = (ev) => {
            console.log(ev);
            const evData = JSON.parse(ev.data);
            const event = Object.assign(new Event, evData);

            routeEvent(event);
        };

        return;
    }

    alert('Not supporting websockets');
}