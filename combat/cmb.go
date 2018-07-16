package combat

import (
	"time"
)

const (
	default_actor_health          = 100
	default_actor_buffer_spots    = 100
	default_actor_reaction_window = 10
	default_instance_actorCount   = 2
	default_instance_tickInterval = time.Second
	minimum_actor_health          = 1
)

/*
   Instance is the main Combat Data Structure.  It contains all of the
   information relevant to combat:  the players involved, the amount of time
   between each attack, and the rules that define the combat.

    Some Definitions

        Tick Interval:  The amount of time between each iteration of logic.
                        The game is not actually turn-based, but if it were:
                        this interval sets the amount of time between turns.
*/
type Instance struct {
	tickInterval time.Duration
	players      []Actor
}

/*
   NewInstance Creates a new combat instance.  The defaults are:
           - 1 second per tick.
           - 2 players
           - 100 total spots
           - 10 spots in reaction window
*/
func NewInstance() *Instance {
	return &Instance{
		tickInterval: default_instance_tickInterval,
		players:      MakeMultipleNewActors(default_instance_actorCount),
	}
}

func (i *Instance) SetInterval(d time.Duration) {
	i.tickInterval = d
}

// AddDefaultPlayer adds a new player to the combat instance, with default attributes and buffers.
func (i *Instance) AddDefaultPlayer() {
	i.players = append(i.players, NewActor())
}

/*
func (i *Instance) String() string {
	playerString := ""
	for _, v := range i.players {
		playerString += "\n\t\t" + v.String()
	}
	return fmt.Sprintf("Instance:{\n\tTickInterval:(%v)\n\tPlayers:[%v\n\t]\n}", i.tickInterval, playerString)
}


func (a Actor) String() string {
    return fmt.Sprintf("Actor{ Attributes:%v, %v", a.attr, a.buffer)
}


func (a ActionBuffer) String() string {
	return fmt.Sprintf("ActionBuffer{ #spots:(%+v), #window:(%v) } ", len(a.spots), a.reactionWindow)
}
*/
