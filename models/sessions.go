package models

import (
	"../auth"

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

var errSessionExpired = errors.New("session is expired")

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
func NewSession(owner Archiver, ipAddr string) Session {
	return Session{
		id:        "",
		owner:     owner.ID(),
		createdAt: time.Now(),
		expiresAt: time.Now().Add(sessionLifetime),
		ipAddress: ipAddr,
	}
}

// FindSession attempts to find a session for an authenticated archiver.
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

// FindSessionByOwnerEmail attempts to find a session owned by an archiver.
func FindSessionByOwnerEmail(db *sql.DB, email string) (Session, error) {
	s := Session{}
	err := db.QueryRow(QFindSessionByOwnerEmail, email).Scan(
		&s.id, &s.owner, &s.createdAt, &s.expiresAt, &s.ipAddress)
	if err != nil {
		return Session{}, err
	}
	return s, nil
}

// ID is a getter function that gets the session's id/token.
func (s Session) ID() string {
	return s.id
}

// Expires is a getter function that gets the session's expire time.
func (s Session) Expires() time.Time {
	return s.expiresAt
}

// Owner is a getter function for a session's owner archiver ID.
func (s Session) Owner() int {
	return s.owner
}

// IsExpired checks if the session has expired.
func (s Session) IsExpired() bool {
	return s.expiresAt.After(time.Now())
}

// Save stores a new session token in the database after making a secure token.
func (s *Session) Save(db *sql.DB) error {
	token, genErr := auth.GenerateUniqueSessionToken(
		sessionTokenLength,
		func(generatedToken string) bool {
			session, _ := FindSession(db, generatedToken)
			return session.ID() != ""
		})
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
func (s *Session) Update(db *sql.DB) error {
	return errors.New("cannot update sessions")
}

// Delete removes a session token from the database, effectively logging an
// archiver out of their account.
func (s *Session) Delete(db *sql.DB) error {
	_, err := db.Exec(QDeleteSession, s.id)
	return err
}
