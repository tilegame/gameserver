package registrar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Registrar is the main object contains a hash map to store the user
// sessions, and several channels which allow it to be accessed
// concurrently.
type Registrar struct {
	userMap map[string]savedSession
	mutex   sync.Mutex
}

// savedSession is used internally as the "value" in the registrar,
// where the "key" is a username string.
type savedSession struct {
	token      []byte
	expiration time.Time
}

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
	SessionDetails map[string]time.Duration
	UserList       []string
}

// NewRegistrar creates a new registrar instance.  Each Registrar has
// a map that stores savedSessions, accessible via username.  The map
// is protected by a Mutex, so all of the calls can be made
// concurrently.
func NewRegistrar() *Registrar {
	return &Registrar{
		userMap: make(map[string]savedSession),
	}
}

// Validate returns true if the username matches the registered token
// and its expiration time has not yet passed.  Otherwise, Validate()
// will return false.  Safe for concurrent use.
func (r *Registrar) Validate(user User) bool {
	// clean the registrar to make sure no outdated tokens get
	// incorrectly validated
	r.Clean()
	r.mutex.Lock()
	defer r.mutex.Unlock()
	s, ok := r.userMap[user.Name]
	return (ok &&
		bytes.Equal(s.token, user.Token) &&
		time.Now().Before(s.expiration))
}

// Add creates a new UserSession in the registrar.  It is save for
// concurrent use.  Any existing entries under the same username will
// be overwritten, so it can be used to update a user's session token
// and session duration.
func (r *Registrar) Add(session UserSession) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.userMap[session.Name] = savedSession{
		token:      session.Token,
		expiration: session.Expiration,
	}
}

// Remove deletes an entry from the registrar map.  Safe for
// concurrent usage.
func (r *Registrar) Remove(name string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.userMap, name)
}

// Clean iterates through the users stored in the registrar, checks to
// see if their tokens have expired, and deletes any users that have
// expired tokens.  This function gets called AUTOMATICALLY whenever
// GenerateInfo or Validate are called, so it SHOULD NOT NEED to be
// called manually.
func (r *Registrar) Clean() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for i, v := range r.userMap {
		if v.expiration.Before(time.Now()) {
			delete(r.userMap, i)
		}
	}
}

// GenerateInfo returns an Info object with information about the
// registrar.  This information can be used in a webpage, turned into
// a JSON, etc.  Does not include the token of any user.
func (r *Registrar) GenerateInfo() *Info {
	// clean the usermap before generating any info.
	r.Clean()
	r.mutex.Lock()
	defer r.mutex.Unlock()
	info := &Info{
		ActiveSessions: len(r.userMap),
		SessionDetails: make(map[string]time.Duration),
	}
	for name, sesh := range r.userMap {
		info.SessionDetails[name] = sesh.expiration.Sub(time.Now())
	}
	return info
}

// HandleInfo returns a webpage with information about the currently
// active sesions
func (r *Registrar) HandleInfo(w http.ResponseWriter, req *http.Request) {
	msg, err := json.MarshalIndent(r.GenerateInfo(), "", "\t")

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
