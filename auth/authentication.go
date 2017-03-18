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

// MaxLoginAttempts is the maximum number of attempts we want to allow
// for a user to attempt to login to any account.
const MaxLoginAttempts int = 5

// AntiCSRFTokenLength is the number of bytes of random data to read in
// order to generate an anti-csrf token.
const AntiCSRFTokenLength uint = 32

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
func IsPasswordCorrect(provided, stored string) bool {
	return auth.CompareHashAndPassword([]byte(stored), []byte(provided)) == nil
}
