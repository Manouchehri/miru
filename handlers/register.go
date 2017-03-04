package handlers

import (
	"../auth"
	"../config"
	"../models"

	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path"
)

// registerPage is the name of the template file containing the register page.
const registerPage string = "register.html"

// RegisterPageHandler implements net/http.ServeHTTP to handle GET requests
// made by users wishing to view the register page.
type RegisterPageHandler struct {
	cfg *config.Config
}

// RegisterHandler implements net/http.ServeHTTP to handle POST requests
// containing the email address and password sent in archiver the register form.
type RegisterHandler struct {
	cfg *config.Config
	db  *sql.DB
}

// NewRegisterPageHandler is the constructor function for a new
// RegisterPageHandler.
// Arguments:
// cfg: A reference to the application's global configuration.
// Returns:
// A new RegisterPageHandler that can be bound to a router.
func NewRegisterPageHandler(cfg *config.Config) RegisterPageHandler {
	return RegisterPageHandler{
		cfg: cfg,
	}
}

// NewRegisterHandler is the constructor function for a new RegisterHandler.
// Arguments:
// cfg: A reference to the application's global configuration.
// db: A reference to a database connection.
// Returns:
// A new RegisterHandler that can be bound to a router.
func NewRegisterHandler(cfg *config.Config, db *sql.DB) RegisterHandler {
	return RegisterHandler{
		cfg: cfg,
		db:  db,
	}
}

// ServeHTTP writes the register page to the requester.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h RegisterPageHandler) ServeHTTP(
	res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles(
		path.Join(h.cfg.TemplateDir, registerPage),
		path.Join(h.cfg.TemplateDir, headTemplate),
		path.Join(h.cfg.TemplateDir, navTemplate))
	if err != nil {
		InternalError(res, req)
		return
	}
	t.Execute(res, nil)
}

// ServeHTTP handles POST requests containing an archiver's registration data.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h RegisterHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	email := req.FormValue("email")
	password := req.FormValue("password")
	passwordRepeated := req.FormValue("passrepeat")

	if password != passwordRepeated {
		fmt.Println("Passwords don't match")
		BadRequest(res, req)
		return
	}
	archiver, _ := models.FindArchiverByEmail(h.db, email)
	if archiver.Email() != "" {
		fmt.Println("Email address taken")
		BadRequest(res, req)
		return
	}
	passwordHash := auth.SecurePassword(password)
	archiver = models.NewArchiver(email, passwordHash)
	saveErr := archiver.Save(h.db)
	if saveErr != nil {
		fmt.Println("Failed to save new archiver", saveErr)
		InternalError(res, req)
		return
	}
	res.Write([]byte(fmt.Sprintf("Successfully registered %s", email)))
}
