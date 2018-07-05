var inputMessage = document.getElementById("text-box");
var outputMessage = document.getElementById("msg-output");

var wSocket = new WebSocket("ws://localhost:8080/echo");

wSocket.onopen = function() {
  outputMessage.innerHTML = "Status: Connected\n";
};

wSocket.onmessage = function(e) {
  outputMessage.innerHTML = "Server :" + e.data + "\n";
};

function send() {
  wSocket.send(inputMessage.value);
  inputMessage.value = "";
}
