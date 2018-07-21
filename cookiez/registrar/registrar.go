package registrar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// User consists of a username and it's secret token, it can be passed
// to the Validate() function to confirm that the username and token
// match each other in the registrar.
type User struct {
	Name  string
	Token []byte
}

// UserSession fully defines a new entry into the registar.  It can be
// passed to the AddUser() function.  When Validate() is called, the
// expiration time is checked.  If the current time is past the
// expiration time, then Validate() returns false.
type UserSession struct {
	User
	Expiration time.Time
}

// Info is for checking on the status of the registrar.  It will
// usually be converted into JSON format and displayed publicly to see
// who's logged in.  Since it's public, it will only show
// non-revealing information like the Usernames, but won't contain
// private things like the session tokens.
//
// ActiveSessions is a count of the number of sessions in the
// registrar, which might be different from the number users are
// actually online.
type Info struct {
	ActiveSessions int
	UserList       []string
}

// savedSession is used internally as the "value" in the registrar,
// where the "key" is a username string.
type savedSession struct {
	token      []byte
	expiration time.Time
}

// query is used internally, adding a sendback channel to the User
// structure.  This enabls the registrar to send a response message
// back with the answer "true" or "false".
type query struct {
	User
	sendback chan bool
}

// registrar is the main object contains a hash map to store the user
// sessions, and several channels which allow it to be accessed
// concurrently.
type registrar struct {
	queryChannel chan query
	entering     chan UserSession
	leaving      chan string
	userMap      map[string]savedSession
}

func newRegistrar() *registrar {
	return &registrar{
		queryChannel: make(chan query),
		entering:     make(chan UserSession),
		leaving:      make(chan string),
		userMap:      make(map[string]savedSession),
	}
}

// running the registrar waits along its channels for a message.
// When one is received, the registrar will either add, delete, or lookup
// an entry.  Use this function in its own go-routine.
func (r *registrar) run() {
	for {
		select {
		case name := <-r.leaving:
			delete(r.userMap, name)

		case s := <-r.entering:
			r.userMap[s.Name] = savedSession{
				token:      s.Token,
				expiration: s.Expiration,
			}

		case q := <-r.queryChannel:
			s, ok := r.userMap[q.Name]
			q.sendback <- (ok &&
				bytes.Equal(s.token, q.Token) &&
				time.Now().Before(s.expiration))
		}
	}
}

func (r *registrar) validate(user User) bool {
	answer := make(chan bool)
	r.queryChannel <- query{user, answer}
	return <-answer
}

func (r *registrar) add(session UserSession) {
	r.entering <- session
}

func (r *registrar) remove(name string) {
	r.leaving <- name
}

func (r *registrar) userTimeBomb(name string, t time.Duration) {
	defer r.remove(name)
	time.Sleep(t)
}

// not for concurrent use.  Call internally from within the registrar.
func (r *registrar) clean() {}

// maybe okay for concurrent use.  Not sure yet.  Copies the whole
// map, then extracts the useful information from the copy.  Remake
// this function when there are more players, and it becomes
// impractical to copy so much data.
func (r *registrar) generateInfo() *Info {
	mapCopy := r.userMap
	info := &Info{
		ActiveSessions: len(mapCopy),
		UserList:       make([]string, len(mapCopy)),
	}
	i := 0
	for name := range mapCopy {
		info.UserList[i] = name
		i++
	}
	return info
}

// --------------------------------------------------------------------

var soloReg = newRegistrar()

func init() {
	go soloReg.run()
}

// Add creates a new UserSession in the registrar.  It is save for
// concurrent use.  Any existing entries under the same username will
// be overwritten, so it can be used to update a user's session token
// and session duration.  This Add() function is blocking, so once it
// returns, you know that the userSession has been successfully added
// (or overwritten) to the registrar.
func Add(session UserSession) {
	soloReg.add(session)
}

// Validate returns true if the username matches the registered token
// and its expiration time has not yet passed.  Otherwise, Validate()
// will return false.  Safe for concurrent use.
func Validate(user User) bool {
	return soloReg.validate(user)
}

// UserTimeBomb creates a ticker countdown that will eventually
// trigger a deletion of the player.  This time cannot be reset,
// because its intended use is for removing old session tokens from
// the registrar.
//
// When UserTimeBomb is called, a new goroutine is launched, and
// UserTimeBomb returns immediately after.  The goroutine sleeps for
// the specified amount of time, and then sends a "delete user with
// <id>" message to the registrar.  Thus, UserTimeBomb is safe for
// concurrent use: because its underlying mechanism is also safe for
// concurrent use.
//
// Note:
//
// There is probably no need to call this function, because the
// registrar will periodically clear itself of old entries.

func UserTimeBomb(name string, timeout time.Duration) {
	go soloReg.userTimeBomb(name, timeout)
}

// --------------------------------------------------------------------

func HandleInfo(w http.ResponseWriter, r *http.Request) {
	msg, err := json.Marshal(soloReg.generateInfo())
	if err != nil {
		log.Println("Registrar:", err)
		fmt.Fprint(w, "error generating registrar info")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(msg)
	if err != nil {
		log.Println("Registrar:", "error sending json message")
	}
}
