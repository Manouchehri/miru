package auth

import auth "github.com/StratumSecurity/scryptauth"

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
