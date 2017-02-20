package models

import (
	"time"
)

// Archiver is the model for a user account, which may be for a regular user
// or for an administrator, who will have permission to create new monitors.
type Archiver struct {
	id           int
	isAdmin      bool
	emailAddress string
	passwordHash string
	loggedInFrom string
	loggedInAt   time.Time
}

// ID is a getter function for an user's identifier.
func (a Archiver) ID() int {
	return a.id
}
