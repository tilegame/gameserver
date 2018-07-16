package gamestate

import (
	"log"
	"time"
)

const (
	Example GameMessageKind = iota
	AddPlayer
	RemovePlayer
)

type Game struct {
	StartTime      time.Time
	MessageChannel chan<- GameMessage
	playerMap      map[int]Player
}

// GameMessage is the structure of the messages interpretted by the Game Hub,
// they are like commands to the tell the game instance what to do.  Examples
// are adding, deleting, or moving players.
type GameMessage struct {
	Kind GameMessageKind
	Data string
}

func NewGame() *Game {
	g := &Game{
		StartTime: time.Now(),
		playerMap: make(map[int]Player),
	}
	go g.runGameMessageHub()
	return g
}

// RunGameMessageHub is meant to be executed in it's own goroutine after
// a nw game is created.  This should happen internally, ensuring that
// it only happens once for each game instance.
func (g *Game) runGameMessageHub() {
	select {
	case r <- request:
		g.handleRequest(r)
	}
}

func (g *Game) handleGameMessage(m GameMessage) {
	switch GameMessage.Kind {
	case Example:
		log.Println("Example message received!")
	case AddPlayer:
		addPlayer(m.data)
	case RemovePlayer:
		delete(g.playerMap, m.data)
	}
}

// addPlayer() is not safe for concurrent execution.  Returns false if
// there is already a player by the given name.
func (g *Game) addPlayer(s string) bool {
	_, ok := m[s]
	if !ok {
		return false
	}
	g.playerMap[s] = Player{}
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
