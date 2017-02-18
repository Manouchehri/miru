package models

import (
	"time"
)

// Administrator is the model for an admin account, belonging to people who
// are allowed to view monitor requests and upload scripts to check for
// updates on sites.
type Administrator struct {
	id           int
	emailAddress string
	passwordHash string
	loggedInFrom string
	loggedInAt   time.Time
}

// ID is a getter function for an Administrator's identifier.
func (a Administrator) ID() int {
	return a.id
}
