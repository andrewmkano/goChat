var inputMessage = document.getElementById("text-box");
window.onload = function() {
  inputMessage.focus();
};
var userName;
var sendBtn = document.getElementById("send-btn");
var outputMessage = document.getElementById("msg-output");
var notificationBell = document.getElementById("notification");
var chanList = document.getElementById("channel-selector");
var channelBtn = document.getElementById("channel-btn");
var chanInput = document.getElementById("channel-box");
var nameInput = document.getElementById("name-box");

function createWebSocket(path) {
  var protocolPrefix = window.location.protocol === "https:" ? "wss:" : "ws:";
  return new WebSocket(protocolPrefix + "//" + location.host + path);
}

var wSocket = createWebSocket("/ws");
var req = new XMLHttpRequest();
wSocket.addEventListener("message", function(e) {
  console.log(e);
});

wSocket.onopen = function() {
  outputMessage.innerHTML += "<pre>" + "Status: Connected" + "\n" + "</pre>";
};

wSocket.onmessage = function(e) {
  var msg = JSON.parse(e.data);
  console.log(msg.Type);
  if (msg.Type == "CHANNEL") {
    updateChanList(msg.Body);
    return;
  }
  if (msg.Type == "PMESSAGE") {
    outputMessage.innerHTML +=
      "<pre>" +
      " -Private- " +
      msg.From +
      " says: " +
      msg.Body +
      "\n" +
      "</pre>";
    var colorToApply = randomColor();
    notificationBell.style.color = colorToApply;
    return;
  }
  if (msg.Type == "NEW_USER") {
    console.log(msg.Body);
    return;
  }
  outputMessage.innerHTML +=
    "<pre>" + msg.Username + " says: " + msg.Text + "\n" + "</pre>";
  var colorToApply = randomColor();
  notificationBell.style.color = colorToApply;
};

wSocket.onclose = function() {
  outputMessage.innerHTML +=
    "<pre>" + "Status: Connection Closed" + "\n" + "</pre>";
};
function getSelectedChannel() {
  var selectedChannel = "";
  selectedChannel = chanList.value;
  return selectedChannel;
}

chanList.onchange = function() {
  channel = getSelectedChannel();
  wSocket.send(
    JSON.stringify({
      type: "CHANGE",
      body: channel
    })
  );
  outputMessage.innerHTML +=
    "<pre>" + "Server: Welcome to " + channel + " Channel!" + "\n" + "</pre>";
};

sendBtn.addEventListener("click", function(msg) {
  msg = inputMessage.value;
  if (inputMessage.value.search("/B") == 0) {
    var BroadcastPrefix = "-BROADCAST- ";
    broadcastMessage = BroadcastPrefix + msg.split("/B").join("");
    wSocket.send(JSON.stringify({ type: "BROADCAST", body: broadcastMessage }));
    inputMessage.value = "";
  } else {
    if (inputMessage.value.search("@") == 0) {
      var inputArray = inputMessage.value.split(/([@])\w+/);
      var Msg = inputArray[2];
      var tmpusr = inputMessage.value.split(" ", 1);
      var TargetUser = tmpusr[0].split("@").join("");
      console.log(inputArray);
      console.log(Msg, TargetUser);
      wSocket.send(
        JSON.stringify({ Type: "PMESSAGE", Body: Msg, Enduser: TargetUser })
      );
      inputMessage.value = "";
    } else {
      chan = getSelectedChannel();
      wSocket.send(JSON.stringify({ type: "MESSAGE", body: msg }));
      inputMessage.value = "";
    }
  }
});

inputMessage.addEventListener("keyup", function(event) {
  if (event.keyCode == 13) {
    document.getElementById("send-btn").click();
  }
});

function captureName() {
  userName = window.prompt("Enter Your nickname", "");
  if (userName != null) {
    wSocket.send(JSON.stringify({ Type: "NEW_USER", Body: userName }));
    console.log(userName);
    return;
  }
}

inputMessage.addEventListener("focus", captureName, { once: true });

channelBtn.addEventListener("click", function() {
  channelName = chanInput.value;
  wSocket.send(JSON.stringify({ Type: "NEW_CHANNEL", Body: channelName }));
  chanInput.value = "";
});

chanInput.addEventListener("keyup", function(evt) {
  if (evt.keyCode == 13) {
    document.getElementById("channel-btn").click();
  }
});

function updateChanList(channels) {
  clearChanList();
  channels.forEach(element => {
    var optionElement = document.createElement("option");
    optionElement.innerHTML = element.ChannelName;
    optionElement.value = element.ChannelName;
    chanList.appendChild(optionElement);
  });
}
function clearChanList() {
  chanList.options.length = 0;
}

function randomColor() {
  //Base string for the Hex color
  var baseColor = "#";
  //All the possible values
  var colorValues = "0123456789ABCDEF";
  for (index = 0; index < 6; index++) {
    baseColor += colorValues[Math.floor(Math.random() * 16)];
  }
  return baseColor;
}
