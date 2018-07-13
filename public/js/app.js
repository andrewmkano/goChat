var inputMessage = document.getElementById("text-box");
window.onload = function() {
  inputMessage.focus();
};
var sendBtn = document.getElementById("send-btn");
var outputMessage = document.getElementById("msg-output");
var notificationBell = document.getElementById("notification");
var chanList = document.getElementById("channel-selector");
var channelBtn = document.getElementById("channel-btn");
var chanInput = document.getElementById("channel-box");
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
  if (msg.Type == "CHANNEL") {
    console.log(msg.body);
    updateChanList(msg.Body);
    return;
  }

  // console.log(msg);
  outputMessage.innerHTML += "<pre>" + msg.Body + "\n" + "</pre>";
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
      body: channel,
      ChannelName: ""
    })
  );
  outputMessage.innerHTML +=
    "<pre>" + "Server: Welcome to" + channel + "Channel!" + "\n" + "</pre>";
};

sendBtn.addEventListener("click", function(msg) {
  msg = inputMessage.value;
  chan = getSelectedChannel();
  wSocket.send(
    JSON.stringify({ type: "MESSAGE", body: msg, ChannelName: chan })
  );
  inputMessage.value = "";
});

inputMessage.addEventListener("keyup", function(event) {
  if (event.keyCode == 13) {
    document.getElementById("send-btn").click();
  }
});

channelBtn.addEventListener("click", function(msg) {
  if (chanInput.value.trim() == "") {
    console.error("The field cannot be blank");
  }

  channelName = chanInput.value;
  wSocket.send(JSON.stringify({ Type: "NEW_CHANNEL", Body: channelName }));
  chanInput.value = "";
});

chanInput.addEventListener("keyup", function(evt) {
  if (evt.keyCode == 13) {
    chanInput.click();
  }
});

function updateChanList(channels) {
  clearChanList();
  channels.forEach(element => {
    var optionElement = document.createElement("option");
    optionElement.innerHTML = element.name;
    optionElement.value = element.name;
    chanList.appendChild(optionElement);
  });
}
function clearChanList() {
  chanList.options.length = 0;
}

function createChannel() {
  if (chanInput.value.trim() == "") {
    return;
  }
  var chanName = chanInput.value;
  wSocket.send(chanName);
  chanInput.value = "";
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
