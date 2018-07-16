package combat

import (
	"fmt"
)

const (
	count_numberOfAttacks   = 3
	count_numberOfLocations = 3

	action_block  = 0
	action_attack = 1
	action_dodge  = 2

	loc_body = 0
	loc_head = 1
	loc_legs = 2
)

/*
   ActionBuffer holds all of the actions that a player wants to commit to.
   Each spot holds the actions.
   Remember, Reaction window must be less than the number of spots.
*/
type ActionBuffer struct {
	spots          []int
	reactionWindow int
}

/*
   Action is an action done by an Actor and stored in an ActionBuffer. (lol).
   Kind - the kind of action it is, such as attack, block, dodge, etc.
   Target - where the action is targeted at.
*/
type Action struct {
	Kind   int
	Target int
}

// NewActionBuffer creates a new action buffer, with default number of spots.
func NewActionBuffer() ActionBuffer {
	return ActionBuffer{
		spots:          make([]int, default_actor_buffer_spots),
		reactionWindow: default_actor_reaction_window,
	}
}

// ExtendDirectly will copy the value at (spot), and then copy each value
// between (spot) and (extendTo), regardless of what is there.
//
func (a *ActionBuffer) ExtendDirectly(spot, extendTo int) error {

	// No need to continue if spot is the same as extendTo.
	// It will always be the same as itself.
	if spot == extendTo {
		return nil
	}
	L := len(a.spots)

	// Check for valid inputs that make sense.
	if spot > L || extendTo > L {
		return fmt.Errorf("ExtendRight: Out of Range. got:(%v), max len: (%v)", spot, len(a.spots))
	}
	if spot < 0 || extendTo < 0 {
		return fmt.Errorf("ExtendRight: Out of Range. Given a key less than 0.")
	}

	// We want to find the first place in the array to start changing values.
	low := spot
	high := extendTo
	if extendTo < spot {
		low = extendTo
		high = spot
	}
	for i := low; i <= high; i++ {
		a.spots[i] = a.spots[spot]
	}
	return nil
}

// Shift removes, and returns, the first element in the action buffer.
func (a *ActionBuffer) Shift() int {
	var x int
	x, a.spots = a.spots[0], a.spots[1:]
	return x
}

// Push adds an element to the end of the action buffer.
func (a *ActionBuffer) Push(int) {
	a.spots = append(a.spots, 0)
}

// Next returns the next element in the action buffer, and adds a 0 to end.
func (a *ActionBuffer) Next() int {
	next := a.Shift()
	a.Push(0)
	return next
}
