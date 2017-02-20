package models

import (
	"../auth"
	"../models"

	"database/sql"
	"errors"
	"time"
)

// sessionLifetime is the amount of time that a user's authenticated session
// should be kept alive for.
const sessionLifetime time.Duration = 1 * time.Hour

// sessionTokenLength is the number of random bytes to generate for new
// session IDs.
const sessionTokenLength uint = 16

// Session contains information about a user's authenticated session.
// Unlike other database entities, a session's ID is a string of
// cryptographically secure random bytes, encoded as hex.
type Session struct {
	id        string
	owner     int
	createdAt time.Time
	expiresAt time.Time
	ipAddress string
}

// NewSession is the constructor function for a new authenticated session,
// which should only be created after verifying that a user's login
// credentials are correct.
// Arguments:
// owner: The archiver on whose behalf the session is being created.
// ipAddr: The IP address that the archiver is connecting from.
// Returns:
// A new Session, which, when saved, will have a token generated for it and
// then be safe to store in a secure cookie.
func NewSession(owner models.Archiver, ipAddr string) Session {
	return Session{
		id:        "",
		owner:     owner.ID(),
		createdAt: time.Now(),
		expiresAt: time.Now().Add(sessionLifetime),
		ipAddress: ipAddr,
	}
}

// FindSession attempts to find a session for an authenticated archiver.
// Arguments:
// db: A database connection.
// id: The session token to look for.
// Returns:
// A session if one with the provided id exists, else an error.
func FindSession(db *sql.DB, id string) (Session, error) {
	s := Session{}
	err := db.QueryRow(QFindSession, id).Scan(
		&s.owner, &s.createdAt, &s.expiresAt, &s.ipAddress)
	if err != nil {
		return Session{}, err
	}
	s.id = id
	return s, nil
}

// Save stores a new session token in the database after making a secure token.
// Arguments:
// db: A database connection.
// Returns:
// An error if one occurs trying to generate a token or save the session.
func (s *Session) Save(db *sql.DB) error {
	token, genErr := auth.GenerateUniqueSessionToken(db, sessionTokenLength)
	if genErr != nil {
		return genErr
	}
	s.id = token
	_, err := db.Exec(QSaveSession,
		s.id, s.owner, s.createdAt, s.expiresAt, s.ipAddress)
	if err != nil {
		s.id = ""
	}
	return err
}

// Update always produces an error.
// Arguments:
// db: A database connection.
// Returns:
// An error- cannot update sessions.
func (s *Session) Update(db *sql.DB) error {
	return errors.New("cannot update sessions")
}

// Delete removes a session token from the database, effectively logging an
// archiver out of their account.
// Arguments:
// db: A database connection.
// Returns:
// An error if the session cannot be deleted.
func (s *Session) Delete(db *sql.DB) error {
	_, err := db.Exec(QDeleteSession, s.id)
	return err
}
