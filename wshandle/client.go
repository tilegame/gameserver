package wshandle

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var idnum = 123

func nextId() int {
	idnum++
	return idnum
}

// Client represents a client connection, and the means of communicating
// with that client.  Client satisfies the io.Writer interface, so you
// can send a concurency-safe message to the client by simply using:
//
//  	fmt.Fprintln(exampleClient, "hello there!")
//
type Client struct {
	Id   int
	room *ClientRoom
	conn *websocket.Conn
	send chan []byte
}

// Unregister informs the Clientroom that this client is leaving.
// Then, it closes the websocket connection.
func (c *Client) Unregister() {
	c.room.remove <- c
	c.conn.Close()
}

// Write to the Client is safe for concurrent use, because it sends the byte array
// through a channel instead of writing it directly to the socket.
func (c *Client) Write(p []byte) (int, error) {
	n := len(p)
	b := make([]byte, len(p))
	copy(b, p)
	c.send <- b
	return n, nil
}

// NewClient creates a new client and run its respective goroutines.
// Pass a reference to the ClientRoom that this client will join,
// and a reference to the websocket connection itself.
func NewClient(room *ClientRoom, conn *websocket.Conn) *Client {
	client := &Client{
		Id:   nextId(),
		room: room,
		conn: conn,
		send: make(chan []byte),
	}
	return client
}

/*
	___________________________________
	              Internals
	===================================
*/

// HandlePong is called whenever a client recieves a websocket "pong"
// from the server.
func (c *Client) handlePong(s string) error {
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	return nil
}

// isUnexpected is a helper function to check for unexpected socket errors.
func isUnexpected(err error) bool {
	return websocket.IsUnexpectedCloseError(
		err,
		websocket.CloseGoingAway,
		websocket.CloseAbnormalClosure)
}

// ReadPump is run once in it's own goroutine.
func (c *Client) readPump() {
	defer c.Unregister()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(c.handlePong)
	for {
		_, message, err := c.conn.ReadMessage()
		// errors are expected when the websocket connection is closing.
		// Only create a log message if there is an Unexpected Error.
		// Afterwords, exit the Read loop.
		if err != nil {
			if isUnexpected(err) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Write the message to all clients in the room.
		// c.room.broadcast <- message

		// send a message down the admin channel.
		c.room.Messages <- Message{
			Id:   c.Id,
			Data: message,
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return
			}

		}
	}

}
