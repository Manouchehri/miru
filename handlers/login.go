package handlers

import (
	"fmt"

	"../auth"
	"../config"
	"../models"

	"database/sql"
	"html/template"
	"net/http"
	"path"
)

// loginPage is the name of the template file containing a login form.
const loginPage string = "login.html"

// LoginPageHandler implements net/http.ServeHTTP to serve a login page.
type LoginPageHandler struct {
	cfg *config.Config
}

// LoginHandler implements net/http.ServeHTTP to handle archiver logins.
type LoginHandler struct {
	cfg *config.Config
	db  *sql.DB
}

// NewLoginPageHandler is the constructor function for a LoginPageHandler.
// Arguments:
// cfg: The application's global configuration.
// Returns:
// A new LoginPageHandler, which can be bound to a router.
func NewLoginPageHandler(cfg *config.Config) LoginPageHandler {
	return LoginPageHandler{
		cfg: cfg,
	}
}

// NewLoginHandler is the constructor function for a LoginHandler.
// Arguments:
// cfg: The application's global configuration.
// db: A database connection.
// Returns:
// A new LoginHandler, which can be bound to a router.
func NewLoginHandler(cfg *config.Config, db *sql.DB) LoginHandler {
	return LoginHandler{
		cfg: cfg,
		db:  db,
	}
}

// ServeHTTP writes the login page to the requester.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h LoginPageHandler) ServeHTTP(
	res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles(
		path.Join(h.cfg.TemplateDir, loginPage),
		path.Join(h.cfg.TemplateDir, headTemplate),
		path.Join(h.cfg.TemplateDir, navTemplate))
	if err != nil {
		InternalError(res, req, h.cfg, errTemplateLoad, false, false)
		return
	}
	t.Execute(res, struct {
		LoggedIn    bool
		UserIsAdmin bool
		Successes   []string
	}{false, false, []string{}})
}

// ServeHTTP handles a login form POST request from the user and attempts to
// establish a new session for them if the supplied credentials are correct.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h LoginHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	email := req.FormValue("email")
	password := req.FormValue("password")
	// Check the provided credentials.
	archiver, findErr := models.FindArchiverByEmail(h.db, email)
	if findErr != nil {
		fmt.Println("Invalid email address")
		BadRequest(res, req, h.cfg, errInvalidCredentials, false, false)
		return
	}
	if !auth.IsPasswordCorrect(password, archiver.Password()) {
		fmt.Println("Incorrect password")
		BadRequest(res, req, h.cfg, errInvalidCredentials, false, false)
		return
	}
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
		InternalError(res, req, h.cfg, errDatabaseOperation, false, false)
		return
	}
	cookie := http.Cookie{
		Name:    auth.SessionCookieName,
		Value:   session.ID(),
		Expires: session.Expires(),
	}
	fmt.Println("Created cookie", cookie)
	http.SetCookie(res, &cookie)
	fmt.Println("Successful login from", email)
	http.Redirect(res, req, "/", http.StatusFound)
}
