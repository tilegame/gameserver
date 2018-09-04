// package echoserver uses websockets and turns all received messages into CAPS.
package echoserver

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

// HandleWs is called by ninjaServer.go and converts all received messages
// into uppercase, and sends it back to the original source.
func HandleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		// Process the message and replace it with a response.
		message = handleMessage(message)

		// Send the response back.
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println(err)
			break
		}
	}
	log.Println(err)
}

type IncomingMessage struct {
	ID  int
	Msg string
}

type ResultMessage struct {
	ID     int         `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

func handleMessage(data []byte) []byte {
	j := new(IncomingMessage)
	json.Unmarshal(data, j)
	result := handleCommand(j.Msg)
	result.ID = j.ID
	outbytes, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
	}
	return outbytes
}

// count is just for testing purposes, remove it later.
var somenumber int = 888

func nextCount() int {
	somenumber++
	return somenumber
}

// Where all the magic happens.  Accepts a command, does something,
// then returns a string as a response.  The boolean indicates an
// error.
func handleCommand(cmd string) ResultMessage {
	out := ResultMessage{}
	switch cmd {
	case "hello":
		out.Result = "well hello to you too!"
	case "gimme":
		out.Result = nextCount()
	default:
		out.Error = "command not found."
	}
	return out
}
