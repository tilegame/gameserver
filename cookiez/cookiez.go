package cookiez

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fractalbach/ninjaServer/cookiez/registrar"
	"github.com/gorilla/securecookie"
)

const (
	MainCookieName  = "yummy-cookie"
	hashKeyLen      = 32 // can be 32 or 64 bytes
	blockKeyLen     = 16 // can be 16, 24, or 32 bytes.
	SessionDuration = time.Minute * 1
)

const loginString = `
You have logged in!

Username: (todo)
PlayerID: %d
Token:    %x
Duration: %s

Try Refreshing the page to see if you stay logged in!
`

const validString = `
You have a validated cookie!  Commands sent with this cookie will be accepted.

Username: (todo)
PlayerID: %d
Token:    %x
TimeLeft: %s
`

var (
	s      = gimmeCookie()
	idIter = 123
)

type userData struct {
	ID    int
	Token []byte
	Expires time.Time
}

func newUserData() userData {
	return userData{
		ID:    nextID(),
		Token: securecookie.GenerateRandomKey(32),
		Expires: time.Now().Add(SessionDuration),
	}
}

func (u *userData) String() string {
	return fmt.Sprintf("id: %d, jumble: %x", u.ID, u.Token)
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
	value := newUserData()
	encoded, err := s.Encode(MainCookieName, value)
	if err != nil {
		log.Println(err)
		return
	}
	cookie := &http.Cookie{
		Name:   MainCookieName,
		Value:  encoded,
		Path:   "/",
		MaxAge: int(SessionDuration.Seconds()),
	}
	http.SetCookie(w, cookie)
	registrar.AddUser(value.ID, value.Token, SessionDuration)
	fmt.Fprintf(w, loginString, value.ID, value.Token, SessionDuration)
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
	value := userData{}
	err2 := s.Decode(MainCookieName, cookie.Value, &value)
	if err2 != nil {
		fmt.Fprintln(w, "You're a strange cookie! Here's a new one.")
		log.Println(err2)
		setCookieHandler(w, r)
		return
	}
	timeLeft := value.Expires.Sub(time.Now())
	fmt.Fprintf(w, validString, value.ID, value.Token, timeLeft)
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
