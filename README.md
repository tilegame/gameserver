# Ninja Server

Repository of the back end for Ninja Arena: https://thebachend.com .
Still in early development.


## Speed Tests

To test your websocket connection on the echo server, go to
https://thebachend.com/speed
which sends 100 messages back and forth, testing the average Round Trip Time
of the web socket messages, including both client and server message processing
times.


## Endpoints

URL | Description
----|-------------
wss://thebachend.com/ws | Main websocket endpoint for the game
wss://thebachend.com/ws/echo | The echo server, used for speed tests and other experiments

## Websockets

*Note: Use the javascript websocket API*:
https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API

Example:
~~~javascript
let conn = new WebSocket("wss://thebachend.com/ws");
conn.onmessage = function(event) {
	 let messages = event.data.split('\n');
	 for (let i = 0; i < messages.length; i++) {
	     console.log(messages[i])
	 }
}
conn.send("hello!");
~~~
