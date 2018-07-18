package cookiez

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
)

const (
	MainCookieName  = "yummy-cookie"
	hashKeyLen  = 32 // can be 32 or 64 bytes
	blockKeyLen = 16 // can be 16, 24, or 32 bytes.
)

var (
	s          = GimmeCookie()
	exampleID  = []byte{1, 2, 3, 4, 5}
	exampleKey = securecookie.GenerateRandomKey(32)
)

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
	value := map[string]([]byte){
		"id":         exampleID,
		"examplekey": exampleKey,
	}
	encoded, err := s.Encode(MainCookieName, value)
	if err != nil {
		log.Println(err)
		return
	}
	cookie := &http.Cookie{
		Name:  MainCookieName,
		Value: encoded,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}

func ReadCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err1 := r.Cookie(MainCookieName)
	if err1 != nil {
		log.Println(err1)
		return
	}
	value := make(map[string]([]byte))
	err2 := s.Decode(MainCookieName, cookie.Value, &value)
	if err2 != nil {
		log.Println(err2)
		return
	}
	fmt.Fprintln(w, value)
}
