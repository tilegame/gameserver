/* 
wshandle implements the HTTP/WebSockets handler for player connections.





Websocket Handler

All players use this websocket handler while logged into the game.
It is similar to a chatroom, and is heavily based on the example from: 
https://github.com/gorilla/websocket/tree/master/examples/chat


To use this on the server, use the HTTP handler, which satisfies the 
http.HandlerFunc interface.
	wshandle.Handle






Sending Messages to Clients

Both the Client and ClientRoom types satisfy the io.Writer interface.
They are both safe for concurrent use!  So feel free to use them without
worrying.

Messages can be broadcast to all the Clients in the ClientRoom:
	fmt.Fprintln(clientroom, "hello, everybody!")

Individual messages can be sent to a specific Client:
	fmt.Fprintln(client, "hello, you!")




Under Construction

There are still some API's to work out, and make it a bit easier to use,
and more clear about what is going on.  Ideally, dealing with the websockets,
and importing gorilla/websockets should _only_ be done in this package.

	TODO:
	- List active clients.
	- Match Client with PlayerSessions 
	- Send incoming messages to GameCommandCenter
	- make it obvious what the main ClientRoom object is called, and how
	  it will be publicly accessible from the rest of the game.


*/
package wshandle
