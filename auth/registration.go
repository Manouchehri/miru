package auth

import (
	"../models"

	"database/sql"

	auth "github.com/StratumSecurity/scryptauth"
)

// DoPasswordsMatch determines whether a password and its repeat are equal.
// DO NOT USE THIS FOR LOGIN.
// A user's password will be hashed before being inserted in the
// database and must be checked with a constant time comparison
// function to avoid timing attacks. Use CredentialsAreCorrect from
// auth/authentication.go for this instead.
// Arguments:
// p1: A password.
// p2: Another password, expected to be p1 repeated.
// Returns:
// True if the passwords are the same, else false.
func DoPasswordsMatch(p1, p2 string) bool {
	return p1 == p2
}

// IsEmailAddressTaken determines whether a user account associated with a
// given email address exists.
// Arguments:
// db: A database connection.
// email: The email address to lookup.
// Returns:
// True if the account exists, false otherwise.
func IsEmailAddressTaken(db *sql.DB, email string) bool {
	archiver, _ := models.FindArchiverByEmail(db, email)
	return archiver.Email() == ""
}

// SecurePassword applies a random salt to and then hashes a password with a
// cryptographically secure password hashing algorithm with a known-secure
// (as of 20/02/2017) configuration.
// Arguments:
// password: A plaintext password.
// Returns:
// The hashes password, suitable to be stored in a database.
func SecurePassword(password string) string {
	hashed, err := auth.GenerateFromPassword(
		[]byte(password), auth.DefaultHashConfiguration())
	for err != nil {
		hashed, err = auth.GenerateFromPassword(
			[]byte(password), auth.DefaultHashConfiguration())
	}
	return string(hashed)
}
