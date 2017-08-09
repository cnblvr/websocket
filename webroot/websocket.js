var sock = null;
var wsuri = "ws://localhost:8080/websocket";

sock = new WebSocket(wsuri);

function randomInteger(min, max) {
    var rand = min - 0.5 + Math.random() * (max - min + 1)
    rand = Math.round(rand);
    return rand;
}

names = ["Ann Lawrence", "Casey Ryan", "Alton Hayes", "Tony Turner", "Gladys Horton", "Anita Bailey", "Lois Figueroa", "Stephanie Cunningham", "Rosa Fowler", "Hubert Weber", "Al Sharp", "Fred Sullivan", "Marlene Palmer", "Charlotte Potter", "Jermaine Gardner", "Jodi Nash", "Bridget Spencer", "Van Miller", "Matt Jensen", "Rodney Bishop", "Ronnie Manning", "Lawrence Mccormick", "Cory Stephens", "Janice Morris", "Antonia Fitzgerald", "Doug Reese", "Harry Harvey", "Salvatore Schwartz", "Noah Washington", "Rodolfo Myers", "Wm Fox", "Roderick Black", "Camille Hoffman", "Sara Knight", "Claude Davidson", "Levi Glover", "Agnes Moss", "Ernestine Collins", "Gilberto Martinez", "Dominick Colon", "Muriel Clark", "Gabriel Medina", "Jody Ramsey", "Eileen Allison", "Pearl Payne", "Adam Sanchez", "Isaac Foster", "Kristine Cruz", "Jerome Frazier", "Tomas Patton"];
colors = ["cornflowerblue", "cadetblue", "crimson", "darkcyan", "darkorange", "darkseagreen", "deeppink",
    "deepskyblue", "dodgerblue", "gold", "khaki", "palevioletred", "royalblue",
    "slateblue", "tomato", "yellowgreen"];

sock.onopen = function () {
    console.debug("connected to " + wsuri);
    setProfile(names[randomInteger(0, names.length-1)], colors[randomInteger(0, colors.length-1)]);
};

sock.onmessage = function (packet) {
    console.debug("message recieved: \"" + packet.data.trim() + "\"");
    var json = JSON.parse(packet.data);
    switch(json.command) {
        case "past_messages":
            AddMessages(json.data);
            break;
        case "new_message":
            AddMessages([json.data]);
            break;
        case "list_users":
            UsersInfo(json.data);
            break;
        default:
            console.log("неизвестная команда у сообщения")
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

function setProfile(name, color) {
    var profile = {};

    if(name !== null) {
        profile.name = name;
        document.title = "чат [" + name + "]";
        document.getElementById("changeName").innerText = name;
    }
    if(color !== null) {
        profile.color = color;
        document.getElementById("changeColor").style.backgroundColor = color;
    }

    if(profile !== {}) {
        _send(JSON.stringify({command:"set_profile", data:profile}));
    }
}

function sendMessage(message) {
    _send(JSON.stringify({command:"send", data:{message:message}}));
}