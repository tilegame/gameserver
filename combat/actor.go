package combat

/*
   Actor is the representation of an entity like a human player or NPC in
   the context of a combat instance.  Each actor has Attributes that are
   relevant to the combat, and an action buffer that contains their plan
   for combat.
*/
type Actor struct {
	attr   Attributes
	buffer ActionBuffer
}

// NewActor creates the default actor object, which has 100 health, a buffer of
// 100 spots, and 10 spots in the reaction window.
func NewActor() Actor {
	return Actor{
		attr:   Attributes{health: default_actor_health},
		buffer: NewActionBuffer(),
	}
}

// Buffer returns a pointer to this Actor's combat action buffer.
func (a *Actor) Buffer() *ActionBuffer {
	return &a.buffer
}

// MakeMultipleNewActors returns an array of identical default actors.
func MakeMultipleNewActors(n int) []Actor {
	out := []Actor{}
	for i := 0; i < n; i++ {
		out = append(out, NewActor())
	}
	return out
}
