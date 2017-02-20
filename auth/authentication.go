package auth

import (
	"errors"

	"../models"

	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"
)

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

// GenerateUniqueSessionToken tries to read random bytes to produce a new
// session token and then checks whether the token is already in use.
// The function will attempt to check the database maxGenerateTokenAttempts
// times before producing an error.
// Arguments:
// db: A database connection.
// length: The number of random bytes to use for the token.
// Returns:
// A unique session token if one is generated and not in use, else an error
// if checking the database for the uniqueness of a token fails too many times.
func GenerateUniqueSessionToken(db *sql.DB, length uint) (string, error) {
	buffer := make([]byte, length)
	readBytes, genErr := rand.Read(buffer)
	attempts := 0
	done := false
	token := ""
	for !done && attempts < maxGenerateTokenAttempts {
		for readBytes != length || genErr != nil {
			<-time.After(generateAttemptWait)
			readBytes, genErr = rand.Read(buffer)
		}
		token = hex.EncodeToString(buffer)
		session, genErr = models.FindSession(db, token)
		attempts++
		done = session.ID() == ""
	}
	if attempts >= maxGenerateTokenAttempts {
		return "", errors.New("could not test for token uniqueness")
	}
	return token, nil
}
