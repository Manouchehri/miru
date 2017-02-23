package auth

import (
	"errors"

	"crypto/rand"
	"encoding/hex"
	"time"

	auth "github.com/StratumSecurity/scryptauth"
)

// SessionCookieName is the name of the cookie to store in the user's
// browser to identify their authenticated session with.
const SessionCookieName string = "mirusession"

// maxGenerateTokenAttempts is the maximum number of times to attempt to
// generate a unique session token.  If GenerateUniqueSessionToken fails to
// check the database this many times to see if a token generated exists,
// then it will return an error.
const maxGenerateTokenAttempts uint = 5

// generateAttemptWait is the amount of time to wait before trying to re-read
// more cryptographically secure random bytes for a session token if the
// first attempt fails. This is done to prevent our source of randomness
// from being exhausted.
const generateAttemptWait time.Duration = 50 * time.Millisecond

// CheckFn is a type alias for a function that accepts a token string and
// determines whether the token is already taken.
type CheckFn func(string) bool

// GenerateUniqueSessionToken tries to read random bytes to produce a new
// session token and then checks whether the token is already in use.
// The function will attempt to check the database maxGenerateTokenAttempts
// times before producing an error.
// Arguments:
// length: The number of random bytes to use for the token.
// taken: A CheckFn that can be used to determine if a token is taken.
// Returns:
// A unique session token if one is generated and not in use, else an error
// if checking the database for the uniqueness of a token fails too many times.
func GenerateUniqueSessionToken(length uint, taken CheckFn) (string, error) {
	buffer := make([]byte, length)
	readBytes, genErr := rand.Read(buffer)
	var attempts uint
	done := false
	token := ""
	for !done && attempts < maxGenerateTokenAttempts {
		for readBytes != int(length) || genErr != nil {
			<-time.After(generateAttemptWait)
			readBytes, genErr = rand.Read(buffer)
		}
		token = hex.EncodeToString(buffer)
		attempts++
		done = !taken(token)
	}
	if attempts >= maxGenerateTokenAttempts {
		return "", errors.New("could not test for token uniqueness")
	}
	return token, nil
}

// IsPasswordCorrect determines whether a provided password matches a stored,
// securely hashed password.
// Arguments:
// provided: A password provided by the user to login with.
// stored: The password hash belonging to the archiver attempting to log in.
// Returns:
// True if the credentials are correct, else false.
func IsPasswordCorrect(provided, stored string) bool {
	return auth.CompareHashAndPassword([]byte(stored), []byte(provided)) == nil
}
