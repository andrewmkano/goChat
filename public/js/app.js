var inputMessage = document.getElementById("text-box");
window.onload = function() {
  inputMessage.focus();
};
var outputMessage = document.getElementById("msg-output");
var notificationBell = document.getElementById("notification");

function createWebSocket(path) {
  var protocolPrefix = window.location.protocol === "https:" ? "wss:" : "ws:";
  return new WebSocket(protocolPrefix + "//" + location.host + path);
}

var wSocket = createWebSocket("/ws");

wSocket.addEventListener("message", function(e) {
  console.log(e);
});

wSocket.onopen = function() {
  outputMessage.innerHTML += "<pre>" + "Status: Connected" + "\n" + "</pre>";
};

wSocket.onmessage = function(e) {
  outputMessage.innerHTML += "<pre>" + e.data + "\n" + "</pre>";
  var colorToApply = randomColor();
  notificationBell.style.color = colorToApply;
};
wSocket.onclose = function() {
  outputMessage.innerHTML +=
    "<pre>" + "Status: Connection Closed" + "\n" + "</pre>";
};

function send() {
  wSocket.send(inputMessage.value);
  inputMessage.value = "";
}
//For the enter button as well
inputMessage.addEventListener("keyup", function(event) {
  if (event.keyCode == 13) {
    document.getElementById("send-btn").click();
  }
});

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
