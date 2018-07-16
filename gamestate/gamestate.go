package gamestate

import (
	"log"
	"time"
)

type GameMessageKind int

const (
	Example GameMessageKind = iota
	AddPlayer
	RemovePlayer
)

type Game struct {
	StartTime      time.Time
	MessageChannel chan GameMessage
	playerMap      map[int]Player
}

// GameMessage is the structure of the messages interpretted by the Game Hub,
// they are like commands to the tell the game instance what to do.  Examples
// are adding, deleting, or moving players.
type GameMessage struct {
	Kind GameMessageKind
	Data interface{}
}

func NewGame() *Game {
	g := &Game{
		StartTime:      time.Now(),
		playerMap:      make(map[int]Player),
		MessageChannel: make(chan GameMessage),
	}
	go g.runGameMessageHub()
	return g
}

// RunGameMessageHub is meant to be executed in it's own goroutine after
// a nw game is created.  This should happen internally, ensuring that
// it only happens once for each game instance.
func (g *Game) runGameMessageHub() {
	select {
	case m := <-g.MessageChannel:
		g.handleGameMessage(m)
	}
}

func (g *Game) handleGameMessage(m GameMessage) {
	switch m.Kind {
	case Example:
		log.Println("Example message received!")

	case AddPlayer:
		id, ok := m.Data.(int)
		if ok {
			g.addPlayer(id)
		} else {
			log.Println("AddPlayerMessage: data needs to be integer")
		}

	case RemovePlayer:
		id, ok := m.Data.(int)
		if ok {
			delete(g.playerMap, id)
		} else {
			log.Println("AddPlayerMessage: data needs to be integer")
		}
	}
}

// addPlayer() is not safe for concurrent execution.  Returns false if
// there is already a player by the given name.
func (g *Game) addPlayer(id int) bool {
	_, ok := g.playerMap[id]
	if !ok {
		return false
	}
	g.playerMap[id] = Player{}
	return true
}

func (g *Game) Uptime() time.Duration {
	return time.Now().Sub(g.StartTime)
}

type Player struct {
	CurrentPosition Location3
	TargetPosition  Location3
}

type Location3 struct {
	x, y, z float64
}
