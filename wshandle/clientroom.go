package wshandle

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Id   int
	Data []byte
}

// ClientRoom handles the list of active clients, and allows messages to be
// broadcast to all of them in a concurrency-safe way.
//
//     fmt.Fprintln(exampleClientRoom, "hello everybody!")
//
// Printing to the ClientRoom itself will broadcast a message to all of the
// clients in the clientroom.
type ClientRoom struct {
	Messages  chan Message
	clientmap map[int]*Client
	broadcast chan []byte
	add       chan *Client
	remove    chan *Client
	upgrader  websocket.Upgrader
}

// NewClientRoom creates the client room object and starts the goroutine
// that manage its client list.
func NewClientRoom() *ClientRoom {
	room := &ClientRoom{
		Messages:  make(chan Message),
		clientmap: map[int]*Client{},
		broadcast: make(chan []byte),
		add:       make(chan *Client),
		remove:    make(chan *Client),
		upgrader:  makeUpgrader(),
	}
	go room.run()
	return room
}

func (r *ClientRoom) run() {
	for {
		select {
		case client := <-r.add:
			r.clientmap[client.Id] = client
			log.Println("client added:", client)

		case client := <-r.remove:
			delete(r.clientmap, client.Id)
			log.Println("client removed:", client)

		case message := <-r.broadcast:
			log.Printf("broadcasting: %s", message)
			for _, client := range r.clientmap {
				select {
				case client.send <- message:
					// message was successfully sent. continue.
				default:
					// something's wrong.  close connection.
					close(client.send)
					delete(r.clientmap, client.Id)
				}
			}
		}
		log.Println("~~~ chatroom status:", len(r.clientmap), "~~~")
	}
}

func (r *ClientRoom) Write(p []byte) (int, error) {
	n := len(p)
	b := make([]byte, len(p))
	copy(b, p)
	r.broadcast <- b
	return n, nil
}

// Handle is the HTTP/WebSocket handler for a given instance of a
// ClientRoom.
func (room *ClientRoom) Handle(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection:", r.RemoteAddr)

	conn, err := room.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// fmt.Fprintln(room, "Welcome:", r.RemoteAddr)

	client := NewClient(room, conn)
	room.add <- client

	go client.readPump()
	go client.writePump()
}

func makeUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
}

// Client returns a reference to a client object by looking up their
// id number.  Useful for sending messages back to the clients.  Since
// this is basically just looking up the client in the map, it returns
// (*Client, bool), similar to the way a map would.
func (room *ClientRoom) Client(id int) (*Client, bool) {
	c, ok := room.clientmap[id]
	return c, ok
}
