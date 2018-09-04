# Ninja Server

[![Build Status](https://travis-ci.org/fractalbach/ninjaServer.svg?branch=master)](https://travis-ci.org/fractalbach/ninjaServer)
[![Go Report Card](https://goreportcard.com/badge/github.com/fractalbach/ninjaServer)](https://goreportcard.com/report/github.com/fractalbach/ninjaServer)
[![GoDoc](https://godoc.org/github.com/fractalbach/ninjaServer?status.svg)](https://godoc.org/github.com/fractalbach/ninjaServer)

Repository of the back end for Ninja Arena
https://thebachend.com 
(Still in early development).


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
https://thebachend.com/cookie | Cookie Handler; used for storing player ID and logins; a precursor to the login page.
https://thebachend.com/sessions | Displays the number of active sessions, and a list of usernames.

## Server Goals: I Can Haz Game?

- Logins
- GameState
- Chat
- Random Game Objects
- Player Interactions
- Arena Battles
- Wins/Losses/Etc.
- Nonvolatile Memory


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
