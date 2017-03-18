package auth

import (
	"strings"

	auth "github.com/StratumSecurity/scryptauth"
)

// PasswordComplexityChecker is used to test whether a given password is complex enough
// to meet the application's security requirements.
type PasswordComplexityChecker struct {
	MinLength    uint
	MinLowercase uint
	MinUppercase uint
	MinSymbols   uint
	MinNumbers   uint
}

// DefaultPasswordComplexityChecker constructs a new PasswordComplexityChecker that is
// configured to require at least 10 characters with 1 lowercase, 1 uppercase,
// 1 symbol, and 1 number.
func DefaultPasswordComplexityChecker() PasswordComplexityChecker {
	return PasswordComplexityChecker{
		MinLength:    10,
		MinLowercase: 1,
		MinUppercase: 1,
		MinSymbols:   1,
		MinNumbers:   1,
	}
}

// IsPasswordSecure checks if a given password passes the security requirements configured.
func (c PasswordComplexityChecker) IsPasswordSecure(password string) bool {
	var lc, uc, s, n uint
	for _, character := range password {
		if character >= 'a' && character <= 'z' {
			lc++
		} else if character >= 'A' && character <= 'Z' {
			uc++
		} else if character >= '0' && character <= '9' {
			n++
		} else {
			s++
		}
	}
	return lc >= c.MinLowercase &&
		uc >= c.MinUppercase &&
		s >= c.MinSymbols &&
		n >= c.MinNumbers &&
		uint(len(password)) >= c.MinLength
}

// SecurePassword applies a random salt to and then hashes a password with a
// cryptographically secure password hashing algorithm with a known-secure
// (as of 20/02/2017) configuration.
func SecurePassword(password string) string {
	hashed, err := auth.GenerateFromPassword(
		[]byte(password), auth.DefaultHashConfiguration())
	for err != nil {
		hashed, err = auth.GenerateFromPassword(
			[]byte(password), auth.DefaultHashConfiguration())
	}
	return string(hashed)
}

// IsEmailValid Does a simple check to make sure an email address is
// formatted reasonably. Determining whether every character is valid is
// infeasible, so we'll just make sure the address is formatted like:
//   <thing>@<domain>.<tld>[.<tld2>[.<tld3>...]]
func IsEmailValid(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	parts = strings.Split(parts[1], ".")
	if len(parts) < 2 {
		return false
	}
	return true
}
