package cookiez

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
)

const (
	MainCookieName = "yummy-cookie"
	hashKeyLen     = 32 // can be 32 or 64 bytes
	blockKeyLen    = 16 // can be 16, 24, or 32 bytes.
)

var (
	s      = GimmeCookie()
	idIter = 123
)

type userPair struct {
	ID     int
	Jumble []byte
}

func newUserPair() userPair {
	return userPair{
		ID:     nextID(),
		Jumble: securecookie.GenerateRandomKey(32),
	}
}

func (u *userPair) String() string {
	return fmt.Sprintf("id: %d, jumble: %x", u.ID, u.Jumble)
}

// increments the package's global variable "idIter", copies that new value,
// and returns the copy.  For use as a new unique player id.
func nextID() int {
	idIter++
	return idIter
}

// GimmeCookie randomly generates keys and returns a Secure Cookie.
func GimmeCookie() *securecookie.SecureCookie {
	hashKey := securecookie.GenerateRandomKey(hashKeyLen)
	blockKey := securecookie.GenerateRandomKey(blockKeyLen)
	// if hashKey == nil || blockKey == nil {
	//	return nil, fmt.Errorf("GenerateRandomKey has returned nil.")
	// }
	return securecookie.New(hashKey, blockKey)
}

// SetCookieHandler is called by the server to hand out cookies.
func SetCookieHandler(w http.ResponseWriter, r *http.Request) {
	value := newUserPair()
	encoded, err := s.Encode(MainCookieName, value)
	if err != nil {
		log.Println(err)
		return
	}
	cookie := &http.Cookie{
		Name:  MainCookieName,
		Value: encoded,
		Path:  "/",
		MaxAge: 60,
	}
	http.SetCookie(w, cookie)
	fmt.Fprintln(w, "Here, have a cookie!")
}

// ReadCookieHandler checks the client's cookies, and prints back a message
// if it's valid.  Does not check yet check to see if the id matches the
// value that it should; simply just confirms that it is a valid cookie.
func ReadCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err1 := r.Cookie(MainCookieName)
	if err1 != nil {
		log.Println(err1)
		return
	}
	value := userPair{}
	err2 := s.Decode(MainCookieName, cookie.Value, &value)
	if err2 != nil {
		fmt.Fprintln(w, "You're a strange cookie! Here's a new one.")
		log.Println(err2)
		SetCookieHandler(w, r)
		return
	}
	fmt.Fprintln(w, "You've got a valid cookie: ", value)
	fmt.Fprintln(w, "Once those 2 values are checked; you'll be logged in. ")
}
