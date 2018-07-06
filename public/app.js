var inputMessage = document.getElementById("text-box");
var outputMessage = document.getElementById("msg-output");
var notificationBell = document.getElementById("notification");

function createWebSocket(path) {
  var protocolPrefix = window.location.protocol === "https:" ? "wss:" : "ws:";
  return new WebSocket(protocolPrefix + "//" + location.host + path);
}

var wSocket = createWebSocket("/g");
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

function send() {
  wSocket.send(inputMessage.value);
  inputMessage.value = "";
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
