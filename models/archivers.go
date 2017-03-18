package models

import (
	"database/sql"
	"errors"
	"time"
)

// Archiver is the model for a user account, which may be for a regular user
// or for an administrator, who will have permission to create new monitors.
type Archiver struct {
	id           int
	madeAdminBy  int
	isAdmin      bool
	emailAddress string
	passwordHash string
	loggedInFrom string
	loggedInAt   time.Time
}

// NewArchiver is the constructor function for a new Archiver, which will be
// created without admin privileges.
func NewArchiver(email string, passwordHash string) Archiver {
	return Archiver{
		id:           -1,
		madeAdminBy:  -1,
		isAdmin:      false,
		emailAddress: email,
		passwordHash: passwordHash,
		loggedInFrom: "",
		loggedInAt:   time.Now(),
	}
}

// ListArchivers obtains a list of all archivers registered in the system.
func ListArchivers(db *sql.DB) ([]Archiver, error) {
	archivers := []Archiver{}
	rows, err := db.Query(QListArchivers)
	if err != nil {
		return archivers, err
	}
	for rows.Next() {
		a := Archiver{}
		err = rows.Scan(
			&a.id, &a.madeAdminBy, &a.isAdmin, &a.emailAddress,
			&a.passwordHash, &a.loggedInFrom, &a.loggedInAt)
		if err != nil {
			break
		}
		archivers = append(archivers, a)
	}
	return archivers, err
}

// FindArchiver attempts to find an archiver in the database with a given id.
// retrieving the account's information fails.
func FindArchiver(db *sql.DB, id int) (Archiver, error) {
	a := Archiver{}
	err := db.QueryRow(QFindArchiver, id).Scan(
		&a.emailAddress, &a.passwordHash, &a.madeAdminBy,
		&a.isAdmin, &a.loggedInFrom, &a.loggedInAt)
	if err != nil {
		return Archiver{}, err
	}
	a.id = id
	return a, nil
}

// FindSessionOwner attempts to find the archiver that owns a session token.
func FindSessionOwner(db *sql.DB, sessionToken string) (Archiver, error) {
	s, err := FindSession(db, sessionToken)
	if err != nil {
		return Archiver{}, err
	}
	return FindArchiver(db, s.Owner())
}

// FindArchiverByEmail attempts to find an Archiver in the database who has
// registered with the provided email address.
func FindArchiverByEmail(db *sql.DB, email string) (Archiver, error) {
	a := Archiver{}
	err := db.QueryRow(QFindArchiverByEmail, email).Scan(
		&a.id, &a.madeAdminBy, &a.isAdmin, &a.passwordHash,
		&a.loggedInFrom, &a.loggedInAt)
	if err != nil {
		return Archiver{}, err
	}
	a.emailAddress = email
	return a, nil
}

// ID is a getter function for an user's identifier.
func (a Archiver) ID() int {
	return a.id
}

// Email is a getter function that gets the archiver's email address.
func (a Archiver) Email() string {
	return a.emailAddress
}

// Password is a getter function that gets the archiver's hashed password,
// for use during login.
func (a Archiver) Password() string {
	return a.passwordHash
}

// IsAdmin is a getter function that determines whether the archiver is an
// administrator.
func (a Archiver) IsAdmin() bool {
	return a.isAdmin
}

// MakeAdmin is a setter function that allows one administrator to give
// administrator privileges to another archiver.
func (a *Archiver) MakeAdmin(authorizedBy Archiver) error {
	if !authorizedBy.isAdmin {
		return errors.New("only administrators can make other users administrators")
	}
	a.madeAdminBy = authorizedBy.id
	a.isAdmin = true
	return nil
}

// canBeMadeAdmin determines whether an archiver is allowed to be given admin
// privileges, which equates to checking if the user who granted the permission
// is themselves an administrator.
func (a Archiver) canBeMadeAdmin(db *sql.DB) bool {
	var canBeAdmin bool
	err := db.QueryRow(QIsUserAnAdmin, a.madeAdminBy).Scan(&canBeAdmin)
	return err == nil && canBeAdmin
}

// Save inserts a new user account into the archivers table. This function also
// double checks that, if the archiver to create has been given administrator
// privileges, that the one who granted them is also an administrator.
func (a *Archiver) Save(db *sql.DB) error {
	if a.isAdmin && !a.canBeMadeAdmin(db) {
		a.madeAdminBy = -1
		a.isAdmin = false
		return errors.New("only administrators can make other archivers an admin")
	}
	_, err := db.Exec(QSaveArchiver,
		a.madeAdminBy, a.isAdmin,
		a.emailAddress, a.passwordHash,
		a.loggedInFrom, a.loggedInAt)
	if err != nil {
		return err
	}
	err = db.QueryRow(QLastRowID).Scan(&a.id)
	return err
}

// Update modifies the existing archiver to change the values of fields which
// may change over the course of a user's existence.
func (a *Archiver) Update(db *sql.DB) error {
	_, err := db.Exec(QUpdateArchiver,
		a.madeAdminBy, a.isAdmin,
		a.emailAddress, a.passwordHash,
		a.loggedInFrom, a.loggedInAt,
		a.id)
	return err
}

// Delete completely removes a user account from the database.
func (a *Archiver) Delete(db *sql.DB) error {
	_, err := db.Exec(QDeleteArchiver, a.id)
	return err
}
