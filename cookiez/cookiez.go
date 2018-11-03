package cookiez

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fractalbach/fractalnet/namegen"
	"github.com/gorilla/securecookie"
	"github.com/tilegame/gameserver/cookiez/registrar"
)

const (
	MainCookieName  = "yummy-cookie"
	hashKeyLen      = 32 // can be 32 or 64 bytes
	blockKeyLen     = 16 // can be 16, 24, or 32 bytes.
	SessionDuration = time.Minute * 1
	MaxAgeSeconds   = 60
)

const loginString = `
You have logged in!

Username: %s
PlayerID: %d
Token:    %x
Duration: %s

Try Refreshing the page to see if you stay logged in!
`

const validString = `
You have a validated cookie!  Commands sent with this cookie will be accepted.

Username: %s
PlayerID: %d
Token:    %x
TimeLeft: %s
`

var (
	reg    = registrar.NewRegistrar()
	s      = gimmeCookie()
	idIter = 123
)

type userData struct {
	ID      int
	Name    string
	Token   []byte
	Expires time.Time
}

func newUserData() userData {
	return userData{
		ID:      nextID(),
		Name:    namegen.GenerateUsername(),
		Token:   securecookie.GenerateRandomKey(32),
		Expires: time.Now().Add(SessionDuration),
	}
}

// increments the package's global variable "idIter", copies that new value,
// and returns the copy.  For use as a new unique player id.
func nextID() int {
	idIter++
	return idIter
}

// gimmeCookie randomly generates keys and returns a Secure Cookie.
func gimmeCookie() *securecookie.SecureCookie {
	hashKey := securecookie.GenerateRandomKey(hashKeyLen)
	blockKey := securecookie.GenerateRandomKey(blockKeyLen)
	// if hashKey == nil || blockKey == nil {
	//	return nil, fmt.Errorf("GenerateRandomKey has returned nil.")
	// }
	return securecookie.New(hashKey, blockKey)
}

// SetCookieHandler is called by the server to hand out cookies.
func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	v := newUserData()
	encoded, err := s.Encode(MainCookieName, v)
	if err != nil {
		log.Println(err)
		return
	}
	cookie := &http.Cookie{
		Name:   MainCookieName,
		Value:  encoded,
		Path:   "/",
		MaxAge: MaxAgeSeconds,
		Secure: true,
	}
	http.SetCookie(w, cookie)
	user := registrar.User{v.Name, v.Token}
	reg.Add(registrar.UserSession{
		user, time.Now().Add(SessionDuration),
	})
	fmt.Fprintf(w, loginString, v.Name, v.ID, v.Token, SessionDuration)
}

// ReadCookieHandler checks the client's cookies, and prints back a message if
// it's valid.  Does not check yet check to see if the id matches the value that
// it should; simply just confirms that it is a valid cookie.
func readCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err1 := r.Cookie(MainCookieName)
	if err1 != nil {
		log.Println(err1)
		return
	}
	v := userData{}
	err2 := s.Decode(MainCookieName, cookie.Value, &v)
	if err2 != nil {
		fmt.Fprintln(w, "Invalid cookie! Here's a new one.")
		log.Println(r.RemoteAddr, err2)
		setCookieHandler(w, r)
		return
	}
	user := registrar.User{
		Name:  v.Name,
		Token: v.Token,
	}
	if reg.Validate(user) {
		timeLeft := v.Expires.Sub(time.Now())
		fmt.Fprintf(w, validString, v.Name, v.ID, v.Token, timeLeft)
		return
	}
	fmt.Fprintln(w, "You have an invalid cookie!!")
	return
}

// ServeCookies is a handler for the cookie server, called by server.go, that
// determines if the client already has the main authentication cookie.  The
// http.ResponseWriter and http.Request are handed off to the cookie handlers in
// package cookiez.  Either a new cookie will be given, or an old cookie will be
// validated.f
func ServeCookies(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(MainCookieName)
	if err != nil {
		setCookieHandler(w, r)
		return
	}
	readCookieHandler(w, r)
}

// HandleInfo displays the information about the registrar created by
// this cookie handler.
var HandleInfo = reg.HandleInfo

// func HandleInfo(w http.ResponseWriter, r *http.Request) {
// 	reg.HandleInfo(w, r)
// }
