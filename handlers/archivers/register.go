package archivers

import (
	"../../auth"
	"../../config"
	"../../models"
	"../common"
	"../fail"
	"../index"

	"database/sql"
	"fmt"
	"net/http"
)

// RegisterHandler implements net/http.ServeHTTP to handle POST requests
// containing the email address and password sent in archiver the register form.
type RegisterHandler struct {
	cfg *config.Config
	db  *sql.DB
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
		fail.BadRequest(res, req, h.cfg, common.ErrBadPassword, false, false)
		return
	}
	if !auth.DefaultPasswordComplexityChecker().IsPasswordSecure(password) {
		fmt.Println("Password is not strong enough")
		fail.BadRequest(res, req, h.cfg, common.ErrBadPassword, false, false)
		return
	}
	if !auth.IsEmailValid(email) {
		fail.BadRequest(res, req, h.cfg, common.ErrInvalidEmail, false, false)
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
		fail.InternalError(res, req, h.cfg, common.ErrDatabaseOperation, false, false)
		return
	}
	handler := index.NewFrontPageHandler(h.cfg, h.db)
	handler.PushSuccessMsg(fmt.Sprintf("You have successfully registered and can now log in with %s.", email))
	handler.ServeHTTP(res, req)
}
