package models

import (
	"database/sql"
	"errors"
	"time"
)

// LoginAttempt is a record of an attempted login that failed.
// We record login attempts for both accounts that do and don't exist.
// This is to prevent enumeration attacks and regular bruteforce attacks.
type LoginAttempt struct {
	id           int
	emailAddress string
	senderIP     string
	madeAt       time.Time
}

// NewLoginAttempt is the constructor function for a new LoginAttempt that
// records an attempt to log into a given email from a given IP.
func NewLoginAttempt(email, from string) LoginAttempt {
	return LoginAttempt{
		emailAddress: email,
		senderIP:     from,
		madeAt:       time.Now(),
	}
}

// FindLoginAttemptsBySender finds all login attempts made from a given IP address.
func FindLoginAttemptsBySender(db *sql.DB, senderIP string) ([]LoginAttempt, error) {
	attempts := []LoginAttempt{}
	rows, err := db.Query(QFindLoginAttemptsBySender, senderIP)
	if err != nil {
		return attempts, err
	}
	for rows.Next() {
		a := LoginAttempt{}
		rows.Scan(&a.id, &a.emailAddress, &a.madeAt)
		a.senderIP = senderIP
		attempts = append(attempts, a)
	}
	return attempts, nil
}

// FindLoginAttemptsByEmail finds all login attempts made for a given email address.
func FindLoginAttemptsByEmail(db *sql.DB, emailRequested string) ([]LoginAttempt, error) {
	attempts := []LoginAttempt{}
	rows, err := db.Query(QFindLoginAttemptsByEmail, emailRequested)
	if err != nil {
		return attempts, err
	}
	for rows.Next() {
		a := LoginAttempt{}
		rows.Scan(&a.id, &a.senderIP, &a.madeAt)
		a.emailAddress = emailRequested
		attempts = append(attempts, a)
	}
	return attempts, nil
}

// AccountEmail is a getter function that retrieves the email address that a login
// attempt was made to access.
func (a LoginAttempt) AccountEmail() string {
	return a.emailAddress
}

// SenderIP is a getter function that retrieves the IP address from which the
// login attempt was made.
func (a LoginAttempt) SenderIP() string {
	return a.senderIP
}

// CreatedAt is a getter function that retrieves the time that the login attempt
// was made at.
func (a LoginAttempt) CreatedAt() time.Time {
	return a.madeAt
}

// Save records a new login attempt.
func (a *LoginAttempt) Save(db *sql.DB) error {
	_, err := db.Exec(QSaveLoginAttempt, a.emailAddress, a.senderIP, a.madeAt)
	if err != nil {
		return err
	}
	err = db.QueryRow(QLastRowID).Scan(&a.id)
	return err
}

// Update always returns an error.
func (a *LoginAttempt) Update(db *sql.DB) error {
	return errors.New("cannot update a login attempt")
}

// Delete destroys a record of a login attempt so that it does not
// count against future login attempts.
func (a *LoginAttempt) Delete(db *sql.DB) error {
	_, err := db.Exec(QDeleteLoginAttempt, a.id)
	return err
}
