package archivers

import (
	"../../auth"
	"../../config"
	"../../models"
	"../common"
	"../fail"

	"database/sql"
	"fmt"
	"net/http"
)

// LoginHandler implements net/http.ServeHTTP to handle archiver logins.
type LoginHandler struct {
	cfg *config.Config
	db  *sql.DB
}

// NewLoginHandler is the constructor function for a LoginHandler.
func NewLoginHandler(cfg *config.Config, db *sql.DB) LoginHandler {
	return LoginHandler{
		cfg: cfg,
		db:  db,
	}
}

// ServeHTTP handles a login form POST request from the user and attempts to
// establish a new session for them if the supplied credentials are correct.
func (h LoginHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	email := req.FormValue("email")
	password := req.FormValue("password")
	csrfToken := req.FormValue("csrfToken")
	fmt.Println("Got CSRF token", csrfToken)
	if !models.VerifyAndDeleteAntiCSRFToken(h.db, csrfToken) {
		fail.BadRequest(res, req, h.cfg, common.ErrNotAllowed, false, false)
		return
	}
	// Prevent people from trying to guess the password to an account.
	attemptedLoginsByUser, _ := models.FindLoginAttemptsBySender(h.db, req.RemoteAddr)
	attemptedLoginsForEmail, _ := models.FindLoginAttemptsByEmail(h.db, email)
	if len(attemptedLoginsByUser) >= auth.MaxLoginAttempts || len(attemptedLoginsForEmail) >= auth.MaxLoginAttempts {
		fail.BadRequest(res, req, h.cfg, common.ErrLoginAttemptsExceeded, false, false)
		return
	}
	// Check the provided credentials.
	archiver, findErr := models.FindArchiverByEmail(h.db, email)
	if findErr != nil {
		fmt.Println("Invalid email address")
		attempt := models.NewLoginAttempt(email, req.RemoteAddr)
		if err := attempt.Save(h.db); err != nil {
			fmt.Println("Error saving login attempt", err)
		}
		fail.BadRequest(res, req, h.cfg, common.ErrInvalidCredentials, false, false)
		return
	}
	if !auth.IsPasswordCorrect(password, archiver.Password()) {
		fmt.Println("Incorrect password")
		attempt := models.NewLoginAttempt(email, req.RemoteAddr)
		if err := attempt.Save(h.db); err != nil {
			fmt.Println("Error saving login attempt", err)
		}
		fail.BadRequest(res, req, h.cfg, common.ErrInvalidCredentials, false, false)
		return
	}
	// When a successful login occurs that does not exceed the maximum number of
	// allowed attempts, delete all recorded attempts so they don't count towards
	// future logins.
	go (func() {
		for _, attempt := range attemptedLoginsByUser {
			attempt.Delete(h.db)
		}
		for _, attempt := range attemptedLoginsForEmail {
			attempt.Delete(h.db)
		}
	})()
	// Delete an old session if one exists.
	oldSession, findErr := models.FindSessionByOwnerEmail(h.db, email)
	if findErr == nil {
		oldSession.Delete(h.db)
	}
	// Establish a session.
	session := models.NewSession(archiver, req.RemoteAddr)
	saveErr := session.Save(h.db)
	if saveErr != nil {
		fmt.Println("Error creating new session", saveErr)
		fail.InternalError(res, req, h.cfg, common.ErrDatabaseOperation, false, false)
		return
	}
	cookie := http.Cookie{
		Name:    auth.SessionCookieName,
		Value:   session.ID(),
		Path:    "/",
		Expires: session.Expires(),
	}
	fmt.Println("Created cookie", cookie)
	http.SetCookie(res, &cookie)
	fmt.Println("Successful login from", email)
	http.Redirect(res, req, "/", http.StatusFound)
}
