package handlers

import (
	"../config"

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
	t, err := template.ParseFiles(path.Join(h.cfg.TemplateDir, loginPage))
	if err != nil {
		InternalError(res, req)
		return
	}
	t.Execute(res, nil)
}

// ServeHTTP handles a login form POST request from the user and attempts to
// establish a new session for them if the supplied credentials are correct.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h LoginHandler) ServeHTTP(
	res http.ResponseWriter, req *http.Request) {
}
