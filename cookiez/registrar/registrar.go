package registrar

import (
	"bytes"
	"time"
)

type user struct {
	id    int
	token []byte
}
type query struct {
	user     user
	sendback chan (bool)
}

type registrar struct {
	queryChannel chan query
	entering     chan user
	leaving      chan int
	userMap      map[int]([]byte)
}

func newRegistrar() *registrar {
	return &registrar{
		queryChannel: make(chan query),
		entering:     make(chan user),
		leaving:      make(chan int),
	}
}

// running the registrar waits along its channels for a message.
// When one is received, the registrar will either add, delete, or lookup
// an entry.  Use this function in its own go-routine.
func (r *registrar) run() {
	select {
	case id := <-r.leaving:
		delete(r.userMap, id)

	case user := <-r.entering:
		r.userMap[user.id] = user.token

	case query := <-r.queryChannel:
		token, ok := r.userMap[query.user.id]
		answer := (ok && bytes.Equal(token, query.user.token))
		query.sendback <- answer
	}
}

// validate is safe for concurrent use.  Returns true or false.
// If the player is "logged-in" based on the given credentials, return true.
func (r *registrar) validate(user user) bool {
	answer := make(chan bool)
	r.queryChannel <- query{user: user, sendback: answer}
	return <-answer
}

// add is safe for concurrent use.  Add a new (id, hash) pair to the registrar.
func (r *registrar) add(user user) {
	r.entering <- user
}

// remove is safe for concurrent use. Removes a player from the registrar.
func (r *registrar) remove(id int) {
	r.leaving <- id
}

// userTimeBomb is safe for concurrent use.
func (r *registrar) userTimeBomb(id int, t time.Duration) {
	defer r.remove(id)
	time.Sleep(t)
}

// --------------------------------------------------------------------

var soloReg = newRegistrar()

func init() {
	go soloReg.run()
}

// AddUser registers a user by their id and token, and starts a timer
// that automatically removes them from the registrar after the timeout.
func AddUser(id int, token []byte, timeout time.Duration) {
	soloReg.add(user{id, token})
	go soloReg.userTimeBomb(id, timeout)
}

// Validate returns true if (id, token) are in the registrar.  Returns false
// if they are not registered. Returns false if the tokens do not match.
func Validate(id int, token []byte) bool {
	return soloReg.validate(user{id, token})
}

// UserTimeBomb creates a ticker countdown that will eventually trigger a
// deletion of the player.  This time cannot be reset, because its intended
// use is for removing old session tokens from the registrar.
//
// When UserTimeBomb is called, a new goroutine is launched, and UserTimeBomb
// returns immediately after.  The goroutine sleeps for the specified amount
// of time, and then sends a "delete user with <id>" message to the registrar.
// Thus, UserTimeBomb is safe for concurrent use: because its underlying
// mechanism is also safe for concurrent use.
//
// Note:
//
// There is probably no need to call this function, because it will be called
// automatically by AddUser(), since in this implementation of the registrar,
// all registered users must have a finite time before being removed.
func UserTimeBomb(id int, timeout time.Duration) {
	go soloReg.userTimeBomb(id, timeout)
}

// --------------------------------------------------------------------
