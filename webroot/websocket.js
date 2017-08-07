var sock = null;
var wsuri = "ws://localhost:8080/websocket";

sock = new WebSocket(wsuri);

sock.onopen = function () {
    console.debug("connected to " + wsuri);
};

sock.onmessage = function (message) {
    console.debug("message recieved: " + message.data);
    if(message == "qwe") {
        send("qwe");
    }
};

sock.onclose = function () {
    console.debug("websocket closed");
};
sock.onerror = function (e) {
    console.debug("error occured: " + e);
};

function _send(msg) {
    if(sock) {
        console.debug("message send: " + msg);
        sock.send(msg);
    }
    else
        console.log("websocket is closed")
};

function sendMessage(message) {
    _send(JSON.stringify({message:message}));
}