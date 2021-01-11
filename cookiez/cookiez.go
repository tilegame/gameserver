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
	mainCookieName  = "tilegame-session"
	hashKeyLen      = 32 // can be 32 or 64 bytes
	blockKeyLen     = 16 // can be 16, 24, or 32 bytes.
	sessionDuration = time.Minute * 1
	maxAgeSeconds   = 60
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

type cookieServer struct {
	reg                 *registrar.Registrar
	gorillaSecureCookie *securecookie.SecureCookie
	uniqueID            int
	secure              bool
}

// Creates a new cookie server that holds a registrar with active user sessions,
// and can be used to generate new cookies for new users. Defaults to using
// secure cookies, which will only work when using TLS, but this can be toggled.
func NewCookieServer() *cookieServer {
	return &cookieServer{
		reg:                 registrar.NewRegistrar(),
		gorillaSecureCookie: gimmeCookie(),
		secure:              true,
		uniqueID:            123,
	}
}

// Set to true in order to serve only secure cookies (which is true by default),
// or change to false to secure insecure cookies, which might be used when
// running the gameserver locally, for example.
func (c *cookieServer) SetCookieSecurity(secure bool) {
	c.secure = secure
}

// Returns a unique id that can be used to store a new player session.
// Increments ids internally.
func (c *cookieServer) nextUniqueID() int {
	c.uniqueID++
	return c.uniqueID
}

type userData struct {
	ID      int
	Name    string
	Token   []byte
	Expires time.Time
}

func (c *cookieServer) newUserData() userData {
	return userData{
		ID:      c.nextUniqueID(),
		Name:    namegen.GenerateUsername(),
		Token:   securecookie.GenerateRandomKey(32),
		Expires: time.Now().Add(sessionDuration),
	}
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
func (c *cookieServer) setCookieHandler(w http.ResponseWriter, r *http.Request) {
	v := c.newUserData()
	encoded, err := c.gorillaSecureCookie.Encode(mainCookieName, v)
	if err != nil {
		log.Println(err)
		return
	}
	cookie := &http.Cookie{
		Name:   mainCookieName,
		Value:  encoded,
		Path:   "/",
		MaxAge: maxAgeSeconds,
		Secure: c.secure,
	}
	http.SetCookie(w, cookie)
	user := registrar.User{
		Name:  v.Name,
		Token: v.Token,
	}
	session := registrar.UserSession{
		User:       user,
		Expiration: time.Now().Add(sessionDuration),
	}
	c.reg.Add(session)
	fmt.Fprintf(w, loginString, v.Name, v.ID, v.Token, sessionDuration)
}

// ReadCookieHandler checks the client's cookies, and prints back a message if
// it's valid.  Does not check yet check to see if the id matches the value that
// it should; simply just confirms that it is a valid cookie.
func (c *cookieServer) readCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err1 := r.Cookie(mainCookieName)
	if err1 != nil {
		c.setCookieHandler(w, r)
		return
	}
	v := userData{}
	err2 := c.gorillaSecureCookie.Decode(mainCookieName, cookie.Value, &v)
	if err2 != nil {
		log.Println(r.RemoteAddr, err2)
		c.setCookieHandler(w, r)
		return
	}
	user := registrar.User{
		Name:  v.Name,
		Token: v.Token,
	}
	if c.reg.Validate(user) {
		timeLeft := v.Expires.Sub(time.Now())
		fmt.Fprintf(w, validString, v.Name, v.ID, v.Token, timeLeft)
		return
	}
	c.deleteCookie(w, r)
	return
}

func (c *cookieServer) deleteCookie(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   mainCookieName,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	fmt.Fprintln(w, "You had an invalid cookie, try refreshing the page.")
}

// ServeCookies is a handler for the cookie server, called by server.go, that
// determines if the client already has the main authentication cookie.  The
// http.ResponseWriter and http.Request are handed off to the cookie handlers in
// package cookiez.  Either a new cookie will be given, or an old cookie will be
// validated.f
func (c *cookieServer) ServeCookies(w http.ResponseWriter, r *http.Request) {
	c.readCookieHandler(w, r)
}

func (c *cookieServer) HandleInfo(w http.ResponseWriter, r *http.Request) {
	c.reg.HandleInfo(w, r)
}
