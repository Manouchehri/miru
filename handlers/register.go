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
		InternalError(res, req, h.cfg, errTemplateLoad, false, false)
		return
	}
	t.Execute(res, struct {
		LoggedIn    bool
		UserIsAdmin bool
		Successes   []string
	}{false, false, []string{}})
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
		BadRequest(res, req, h.cfg, errBadPassword, false, false)
		return
	}
	if !auth.DefaultPasswordComplexityChecker().IsPasswordSecure(password) {
		fmt.Println("Password is not strong enough")
		BadRequest(res, req, h.cfg, errBadPassword, false, false)
		return
	}
	if !auth.IsEmailValid(email) {
		BadRequest(res, req, h.cfg, errInvalidEmail, false, false)
		return
	}
	archiver, _ := models.FindArchiverByEmail(h.db, email)
	// We don't want to tell users if an email address is taken so that it is
	// impossible to enumerate registered accounts.
	// TODO - When we have confirmation emails being sent, we should say that
	// an email has been sent in both the case that the email is taken and
	// in the case that it is not.
	if archiver.Email() != "" {
		res.Write([]byte(fmt.Sprintf("Successfully registered %s", email)))
		return
	}
	passwordHash := auth.SecurePassword(password)
	archiver = models.NewArchiver(email, passwordHash)
	saveErr := archiver.Save(h.db)
	if saveErr != nil {
		fmt.Println("Failed to save new archiver", saveErr)
		InternalError(res, req, h.cfg, errDatabaseOperation, false, false)
		return
	}
	handler := NewIndexHandler(h.cfg, h.db)
	handler.PushSuccessMsg(fmt.Sprintf("You have successfully registered and can now log in with %s.", email))
	handler.ServeHTTP(res, req)
}
