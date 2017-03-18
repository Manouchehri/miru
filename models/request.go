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
	rejected     bool
}

// NewRequest is the constructor function for a new request to have a site monitored.
func NewRequest(creator Archiver, url, instructions string) Request {
	return Request{
		id:           -1,
		createdBy:    creator.ID(),
		createdAt:    time.Now(),
		url:          url,
		instructions: instructions,
		rejected:     false,
	}
}

// ID is a getter function for a request's unique identifier.
func (r Request) ID() int {
	return r.id
}

// FindRequest attempts to find an existing monitor request given its ID.
func FindRequest(db *sql.DB, id int) (Request, error) {
	r := Request{}
	err := db.QueryRow(QFindRequest, id).Scan(
		&r.createdBy, &r.createdAt, &r.url, &r.instructions, &r.rejected)
	if err != nil {
		return Request{}, err
	}
	r.id = id
	return r, nil
}

// ListPendingRequests attempts to find all requests that have not had monitors
// created for them yet.
func ListPendingRequests(db *sql.DB) ([]Request, error) {
	requests := []Request{}
	rows, err := db.Query(QListPendingRequests)
	if err != nil {
		return requests, err
	}
	for rows.Next() {
		r := Request{}
		err = rows.Scan(&r.id, &r.createdBy, &r.createdAt, &r.url, &r.instructions)
		if err != nil {
			return []Request{}, err
		}
		r.rejected = false
		requests = append(requests, r)
	}
	return requests, nil
}

// URL is a getter function for the URL that a request was made to monitor.
func (r Request) URL() string {
	return r.url
}

// Instructions is a getter function for the instructions provided to help
// write a monitor script for the site.
func (r Request) Instructions() string {
	return r.instructions
}

// Creator is a getter function for the ID of the archiver that created
// the request.
func (r Request) Creator() int {
	return r.createdBy
}

// Save inserts a new request into the requests table.
func (r *Request) Save(db *sql.DB) error {
	_, err := db.Exec(QSaveRequest, r.createdBy, r.createdAt, r.url, r.instructions)
	if err != nil {
		return err
	}
	err = db.QueryRow(QLastRowID).Scan(&r.id)
	return err
}

// Update always returns an error, as requests cannot be changed once made.
func (r *Request) Update(db *sql.DB) error {
	return errors.New("cannot update a monitor request")
}

// Delete removes a request from the database if it has not already been fulfilled
// and had a monitor script created for it. It is a way to reject requests only.
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
