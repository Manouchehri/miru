package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"
)

// tokenLifetime is the amount of time a token should be considered valid for
// from the time of its creation.
const tokenLifetime = 1 * time.Hour

// AntiCSRFToken contains a non-secret token that is embedded into a
// hidden input field in forms that are used to conduct sensitive actions.
// See OWASP's information about Cross-Site Request Forgery for more information.
// https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
type AntiCSRFToken struct {
	token     string
	createdAt time.Time
}

// GenerateAntiCSRFToken produces a new AntiCSRFToken with a token of a given
// number of bytes. Note that the token itself will be hex-encoded, and thus
// occupy 2 * tokenLen characters.
func GenerateAntiCSRFToken(db *sql.DB, tokenLen uint) AntiCSRFToken {
	return AntiCSRFToken{
		token:     generateToken(db, tokenLen),
		createdAt: time.Now(),
	}
}

// FindAntiCSRFToken attempts to find an anti-csrf token in the database.
func FindAntiCSRFToken(db *sql.DB, token string) (AntiCSRFToken, error) {
	t := AntiCSRFToken{}
	err := db.QueryRow(QFindAntiCSRFToken, token).Scan(&t.createdAt)
	return t, err
}

// Token is a getter for the anti-csrf token's actual token string, to
// be embedded into a form.
func (t AntiCSRFToken) Token() string {
	return t.token
}

// IsExpired determines whether a token should be considered void for
// having been submitted too long after its creation. Expired tokens
// should be deleted immediately.
func (t AntiCSRFToken) IsExpired() bool {
	return time.Now().After(t.createdAt.Add(tokenLifetime))
}

// Save inserts a newly generated token into the database.
func (t *AntiCSRFToken) Save(db *sql.DB) error {
	_, err := db.Exec(QSaveAntiCSRFToken, t.token, t.createdAt)
	return err
}

// Update always returns an error.
func (t *AntiCSRFToken) Update(db *sql.DB) error {
	return errors.New("cannot update an anti-csrf token")
}

// Delete removes the token from the database. Once received in a
// request, a token should be deleted immediately to prevent
// reuse.
func (t *AntiCSRFToken) Delete(db *sql.DB) error {
	_, err := db.Exec(QDeleteAntiCSRFToken, t.token)
	return err
}

func generateToken(db *sql.DB, tokenLen uint) string {
	buffer := make([]byte, tokenLen)
	token := ""
	for {
		bytesRead, err := rand.Read(buffer)
		if err != nil || uint(bytesRead) < tokenLen {
			continue
		}
		token = hex.EncodeToString(buffer)
		_, err = FindAntiCSRFToken(db, token)
		if err != nil {
			break
		}
	}
	return token
}
