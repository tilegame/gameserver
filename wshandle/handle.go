package wshandle

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"fmt"
)

var clientroom = NewClientRoom()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func Handle(w http.ResponseWriter, r *http.Request) {

	log.Println("new connection:", r.RemoteAddr)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Fprintln(clientroom, "Welcome:", r.RemoteAddr)
	
	client := NewClient(clientroom, conn)
	clientroom.add <- client

	go client.readPump()
	go client.writePump()
}

