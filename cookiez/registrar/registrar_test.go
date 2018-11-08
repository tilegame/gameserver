package registrar

import (
	"testing"
	"time"
)

/*

TODO
 - test that listing info works correctly in different scenarios.
 - test for concurrent adds, removes, and listing.

*/

const expireDuration = 100 * time.Millisecond

var (
	name       = "testUser123"
	token      = []byte("arbitrary sequence of bytes.")
	expiration = time.Now().Add(expireDuration)
	user       = User{name, token}
	session    = UserSession{user, expiration}
)

// Test the most simplistic form of the APIs.
func TestBasicAPI(t *testing.T) {
	r := NewRegistrar()
	r.Add(session)
	r.Remove(name)
}

func TestValidate(t *testing.T) {
	r := NewRegistrar()
	r.Add(session)

	// Try to validate a real user who was just added to the
	// registar.
	if !r.Validate(user) {
		t.Error("unable to validate a valid user!")
	}
	t.Log("Valid User was successfully reported Valid.")

	// Try to validate a fake user who wasn't added to the
	// registar.
	fake := User{"invalid user", []byte("lolzors")}
	if r.Validate(fake) {
		t.Error("A fake user was validated, but was never added!")
	}
	t.Log("Invalid User was successfully reported Invalid..")

	// Remove the user from the registrar, then try to validate
	// Expected behavior: user not valid.
	r.Remove(user.Name)
	if r.Validate(user) {
		t.Error("the user was removed, but still was validated!")
	}
	t.Log("user successfully removed, and then Not validated.")
}

// TestExpiration adds a user, then waits until the expiration of that
// token has passed.  Then, the user tries to validate.  The test
// passes if they correctly aren't reported Invalid after the
// expiration time has passed.
func TestExpiration(t *testing.T) {
	r := NewRegistrar()
	r.Add(session)
	time.Sleep(expireDuration)
	if r.Validate(user) {
		t.Error("The token has expired, but the user was still validated!")
	}
	t.Log("The token expired, and the user was succesfully reported Invalid.")
}
