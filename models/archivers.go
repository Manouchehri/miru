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
// Arguments:
// email: The email address that the user wants to register with.
// passwordHash: The user's password, with scrypt applied.
// Returns:
// A new Archiver, which we can call Save() on.
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

// FindArchiver attempts to find an archiver in the database with a given id.
// Arguments:
// db: A database connection.
// id: The unique identifier of the archiver to look up.
// Returns:
// An Archiver instance if a user exists, and an error if one does not, or
// retrieving the account's information fails.
func FindArchiver(db *sql.DB, id int) (Archiver, error) {
	a := Archiver{}
	err := db.QueryRow(QFindArchiver, id).Scan(
		&a.emailAddress, &a.passwordHash, &a.madeAdminBy,
		&a.isAdmin, &a.loggedInFrom, &a.loggedInAt)
	if err != nil {
		return Archiver{}, err
	}
	return a, nil
}

// FindSessionOwner attempts to find the archiver that owns a session token.
// Arguments:
// db: A database connection.
// sessionToken: A session token read from a cookie.
// Returns:
// The Archiver owning the session provided or any error that occurs trying
// to find them.
func FindSessionOwner(db *sql.DB, sessionToken string) (Archiver, error) {
	s, err := FindSession(db, sessionToken)
	if err != nil {
		return Archiver{}, err
	}
	return FindArchiver(db, s.Owner())
}

// FindArchiverByEmail attempts to find an Archiver in the database who has
// registered with the provided email address.
// Arguments:
// db: A database connection.
// email: The email address to look for an account associated with.
// Returns:
// An Archiver instance if a user exists, and an error if one does not, or
// retrieving the account's information fails.
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
// Returns:
// The archiver's id in the archivers table.
func (a Archiver) ID() int {
	return a.id
}

// Email is a getter function that gets the archiver's email address.
// Returns:
// The archiver's email address.
func (a Archiver) Email() string {
	return a.emailAddress
}

// Password is a getter function that gets the archiver's hashed password,
// for use during login.
// Returns:
// The archiver's hashed password.
func (a Archiver) Password() string {
	return a.passwordHash
}

// IsAdmin is a getter function that determines whether the archiver is an
// administrator.
// Returns:
// True if the archiver has admin privileges, else false.
func (a Archiver) IsAdmin() bool {
	return a.isAdmin
}

// MakeAdmin is a setter function that allows one administrator to give
// administrator privileges to another archiver.
// Arguments:
// authorizedBy: The archiver granting the permissions.
// Returns:
// An error if the archiver granting permission is not themselves an admin.
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
// Arguments:
// db: A database connection
// Returns:
// True if it is okay for the archiver to have admin privileges, else false.
func (a Archiver) canBeMadeAdmin(db *sql.DB) bool {
	var canBeAdmin bool
	err := db.QueryRow(QIsUserAnAdmin, a.madeAdminBy).Scan(&canBeAdmin)
	return err == nil && canBeAdmin
}

// Save inserts a new user account into the archivers table. This function also
// double checks that, if the archiver to create has been given administrator
// privileges, that the one who granted them is also an administrator.
// Arguments:
// db: A database connection.
// Returns:
// An error if the user who made this archiver an admin is not an admin
// themselves, or any error that occurs in the database.
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
// Arguments:
// db: A database connection.
// Returns:
// Any errors from the database.
func (a *Archiver) Update(db *sql.DB) error {
	_, err := db.Exec(QUpdateArchiver,
		a.madeAdminBy, a.isAdmin,
		a.emailAddress, a.passwordHash,
		a.loggedInFrom, a.loggedInAt,
		a.id)
	return err
}

// Delete completely removes a user account from the database.
// Arguments:
// db: A database connection.
// Returns:
// Any errors from the database.
func (a *Archiver) Delete(db *sql.DB) error {
	_, err := db.Exec(QDeleteArchiver, a.id)
	return err
}
