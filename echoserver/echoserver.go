// package echoserver uses websockets and is currently under construction.
package echoserver

import (
	"encoding/json"
	"fmt"
	"time"
	//"github.com/fractalbach/ninjaServer/gamestate"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type client struct {
	conn *websocket.Conn
}

var clientlist = map[*client]bool{}

func broadcast(data []byte) {
	for client := range clientlist {
		client.conn.WriteMessage(websocket.TextMessage, data)
	}
}

// HandleWs is called by ninjaServer.go and converts all received
// messages into uppercase, and sends it back to the original source.
func HandleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// create the client object and add it to the list.
	me := &client{
		conn: conn,
	}
	clientlist[me] = true

	// Remove self from the client list when function returns.
	defer func() {
		delete(clientlist, me)
		conn.Close()
	}()

	// wait for new incoming messages.
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		// Process the message and replace it with a response.
		message = send(message)

		// Send the response back.
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println(err)
			break
		}
	}
	log.Println(err)
}

// ~~~~~~~~~~~~~~~~~~ Game ~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TODO: move these out of this file, and make a seperate file for the
// state of the game.  Variables related to the game.
var (
	nextPlayerId  = 136
	playerlist    = map[string]*Player{}
	activeplayers = map[string]bool{}

	TICK_DURATION               = time.Millisecond * 500
	PLAYERLIST_REFRESH_DURATION = 3 * time.Minute
	StartTickerChan             = make(chan bool)
	StopTickerChan              = make(chan bool)
)

func init() {
	go runPostOffice()
	go runGameTicker()
}

func getNextPlayerId() int {
	nextPlayerId++
	return nextPlayerId
}

type IncomingMessage struct {
	ID     int
	Method string
	Params []interface{}
}

type ResultMessage struct {
	ID      int         `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Kind    string      `json:"kind,omitempty"`
	Comment string      `json:"comment,omitempty"`
}

type postcard struct {
	data []byte
	ret  chan []byte
}

func send(data []byte) []byte {
	myMailbox := make(chan []byte)
	card := postcard{
		data: data,
		ret:  myMailbox,
	}
	postoffice <- card
	answer := <-myMailbox
	return answer
}

// postoffice is the name of the channel that you can send postcards
// to!  It only works if you call runPostOffice() in a goroutine.
var postoffice = make(chan postcard)

// run this in it's own goroutine.
func runPostOffice() {
	for {
		select {
		case m := <-postoffice:
			m.ret <- handleMessage(m.data)
		}
	}
}

func handleMessage(data []byte) []byte {

	// Incoming Bytes -> Incoming Json
	j := new(IncomingMessage)
	err := json.Unmarshal(data, j)

	// TODO: Handle badly formatted json by sending back an error.
	if err != nil {
		return []byte{}
	}

	// Pass the Incoming Json to the Command Handler.  That
	// command handler returns a Result Json.
	result := handleCommand(j.Method, j.Params)

	// Add the message id number to the Result Json.  This allows
	// the client to identify their Json again.
	result.ID = j.ID

	// Convert Result Json back into bytes so that it can be sent
	// back to the client.
	outbytes, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
	}
	return outbytes
}

// Where all the magic happens.  Accepts a command, does something,
// then returns a string as a response.  The boolean indicates an
// error.
func handleCommand(cmd string, params []interface{}) ResultMessage {
	out := ResultMessage{}
	switch cmd {
	case "hello":
		out.Result = "well hello to you too!"

	case "params":
		out.Comment = "type info about the given parameters."
		s := make([]string, len(params))
		for i, _ := range params {
			q := params[i]
			s[i] = fmt.Sprintf("(%T, %v)", q, q)
		}
		out.Result = s

	case "add":
		return doAddCmd(params)

	case "remove":
		handleRemove(params, &out)

	case "list":
		out.Kind = "playerlist"
		out.Result = playerlist

	case "chat":
		doChatCmd(params, &out)

	case "move":
		return doMoveCmd(params)

	case "update":
		for _, p := range playerlist {
			p.UpdatePosition()
		}
		out.Result = true

	default:
		out.Error = "command not found."
	}
	return out
}

// Adds a player to the playerlist.
func doAddCmd(params []interface{}) ResultMessage {
	out := ResultMessage{}

	// throw error if parameters number is wrong.
	if len(params) != 1 {
		out.Error = "expected 1 parameter."
		return out
	}

	// convert the first param to  type string.
	name, ok := params[0].(string)
	if !ok {
		out.Error = "Parameter Type error"
		return out
	}

	// add the player to the playerlist.
	playerlist[name] = &Player{
		PlayerId:   getNextPlayerId(),
		CurrentPos: Loc{5, 5},
		TargetPos:  Loc{5, 5},
	}
	out.Result = true
	setPlayerToActive(name)
	return out
}

func handleRemove(params []interface{}, out *ResultMessage) {
	if len(params) != 1 {
		out.Error = "expected 1 param: (username)"
		return
	}
	name, ok := params[0].(string)
	if !ok {
		out.Error = "type error"
		return
	}
	delete(playerlist, name)
	out.Result = true
}

func doMoveCmd(params []interface{}) ResultMessage {
	out := ResultMessage{}

	// param length check.
	if len(params) != 3 {
		out.Error = "expected 3 params: (username, x, y)"
		return out
	}

	// convert types.
	name, ok1 := params[0].(string)
	x, ok2 := params[1].(float64)
	y, ok3 := params[2].(float64)
	if !(ok1 && ok2 && ok3) {
		out.Error = "type error: expected (string, int, int)"
		return out
	}

	// retrieve player pointer.
	p, ok4 := playerlist[name]
	if !ok4 {
		out.Error = "player does not exist."
		return out
	}

	// update positions and return successfull.
	p.TargetPos.X = int(x)
	p.TargetPos.Y = int(y)
	out.Result = true
	setPlayerToActive(name)
	return out
}

func doChatCmd(params []interface{}, out *ResultMessage) {
	if len(params) != 2 {
		out.Error = "expected 2 params: (username, chatmessage)"
		return
	}
	name, ok0 := params[0].(string)
	chatmessage, ok1 := params[1].(string)
	if !(ok0 && ok1) {
		out.Error = "type error: expected (string, string)"
		return
	}
	if _, ok := playerlist[name]; !ok {
		out.Error = "player does not exist."
		return
	}	
	msgtoall := ResultMessage{
		Kind: "chat",
		Result: map[string]string{
			"User":    name,
			"Message": chatmessage,
		},
	}
	data, err := json.Marshal(msgtoall)
	if err != nil {
		log.Println(err)
		return
	}
	broadcast(data)
	return
}

// converts the playerlist into an array of usernames, friendly for
// displaying the list of added players.
func sprintPlayerList() []string {
	s := make([]string, len(playerlist))
	i := 0
	for key, _ := range playerlist {
		s[i] = key
		i++
	}
	return s
}

type Loc struct {
	X int
	Y int
}

type Player struct {
	PlayerId   int
	CurrentPos Loc
	TargetPos  Loc
}

func (p *Player) UpdatePosition() {
	next := p.CurrentPos

	if p.CurrentPos.X == p.TargetPos.X {
		goto Skip1
	}
	if p.CurrentPos.X < p.TargetPos.X {
		next.X = p.CurrentPos.X + 1
	} else {
		next.X = p.CurrentPos.X - 1
	}
Skip1:
	if p.CurrentPos.Y == p.TargetPos.Y {
		goto Skip2
	}
	if p.CurrentPos.Y < p.TargetPos.Y {
		next.Y = p.CurrentPos.Y + 1
	} else {
		next.Y = p.CurrentPos.Y - 1
	}
Skip2:
	// Check for collisions at new x-position
	if NoCollisionAt(next.X, p.CurrentPos.Y) {
		p.CurrentPos.X = next.X
	}

	// Check for collisions at new y-position
	if NoCollisionAt(p.CurrentPos.X, next.X) {
		p.CurrentPos.Y = next.Y
	}
}

// TODO: add collision checking; currently there are no collisions.
//
// NoCollisionAt checks the tile at (x,y) to see if there is something
// that might prevent movement to that tile.
func NoCollisionAt(x, y int) bool {
	return true
}

// The Game Ticker continuously updates the game, by checking if
// players have moved, and updating their positions at fixed
// intervals.
func runGameTicker() {
	// note: the ability to create additional tickers might result
	// in goroutine leaks, because the docs say that the
	// ticker.Stop() does not actually close the channel.
	ticker := time.NewTicker(TICK_DURATION)
	refresher := time.NewTicker(PLAYERLIST_REFRESH_DURATION)
	for {
		select {
		case <-StartTickerChan:
			ticker = time.NewTicker(TICK_DURATION)
		case <-StopTickerChan:
			ticker.Stop()
		case <-ticker.C:
			doGameTick()
		case <-refresher.C:
			doPlayerlistRefresh()
		}
	}
}

func doGameTick() {
	if len(playerlist) == 0 || len(clientlist) == 0 {
		return
	}
	send([]byte(`{"method": "update"}`))
	out := send([]byte(`{"method": "list"}`))
	broadcast(out)
}

// looks through the list of players, and compares them to the
// activeplayer set.  If there are any players who AREN'T in the
// activeplayer set, then they are considered "inactive" and will be
// logged out.
func doPlayerlistRefresh() {
	for name := range playerlist {
		_, ok := activeplayers[name]
		if !ok {
			delete(playerlist, name)
		}
	}
	// clear the activeplayer set.
	activeplayers = make(map[string]bool)
}

// call setPlayerToActive when a player does some action that allows
// it to stay logged in.
func setPlayerToActive(name string) {
	activeplayers[name] = true
}
