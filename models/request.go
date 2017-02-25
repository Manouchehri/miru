package models

import (
	"database/sql"
	"errors"
	"time"
)

// Request is the model for a request, created by an archiver, to have
// a site monitored.
type Request struct {
	id           int
	createdBy    int
	createdAt    time.Time
	url          string
	instructions string
}

// NewRequest is the constructor function for a new request to have a site monitored.
// Arguments:
// creator: The archiver who issued the request.
// url: The URL of the site to monitor.
// instructions: Any further instructions needed to write a monitor script.
// Returns:
// A new Request instance, which we can call Save() on.
func NewRequest(creator Archiver, url, instructions string) Request {
	return Request{
		id:           -1,
		createdBy:    creator.ID(),
		createdAt:    time.Now(),
		url:          url,
		instructions: instructions,
	}
}

// ID is a getter function for a request's unique identifier.
// Returns:
// The request's id in the database.
func (r Request) ID() int {
	return r.id
}

// Save inserts a new request into the requests table.
// Arguments:
// db: A database connection.
// Returns:
// An error if the database insertion fails.
func (r *Request) Save(db *sql.DB) error {
	_, err := db.Exec(QSaveRequest, r.createdBy, r.createdAt, r.url, r.instructions)
	if err != nil {
		return err
	}
	err = db.QueryRow(QLastRowID).Scan(&r.id)
	return err
}

// Update always returns an error, as requests cannot be changed once made.
// Arguments:
// db: A database connection.
// Returns:
// An error saying that the operation is not allowed.
func (r *Request) Update(db *sql.DB) error {
	return errors.New("cannot update a monitor request")
}

// Delete removes a request from the database if it has not already been fulfilled
// and had a monitor script created for it. It is a way to reject requests only.
// Arguments:
// db: A database connection.
// Returns:
// An error if the request is already fulfilled or a database error occurs.
func (r *Request) Delete(db *sql.DB) error {
	isFulfilled := false
	err := db.QueryRow(QIsRequestFulfilled, r.id).Scan(&isFulfilled)
	if err != nil {
		return err
	}
	if isFulfilled {
		return errors.New("cannot delete a fulfilled request")
	}
	_, err = db.Exec(QRejectRequest, r.id)
	return err
}
