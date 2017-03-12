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
	"strconv"
)

// PromoteHandler handles requests from administrators to promote an archiver to
// have administrator privileges.
type PromoteHandler struct {
	cfg *config.Config
	db  *sql.DB
}

// NewPromoteHandler is the constructor for a new PromoteHandler.
// Arguments:
// db: A database connection.
// Returns:
// A new PromoteHandler which can be bound to a router.
func NewPromoteHandler(cfg *config.Config, db *sql.DB) PromoteHandler {
	return PromoteHandler{
		cfg: cfg,
		db:  db,
	}
}

// ServeHTTP handles requests to have an archiver promoted to become an administrator.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h PromoteHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Check that the request is coming from an authenticated administrator.
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		fmt.Println("Could not find cookie", err)
		fail.BadRequest(res, req, h.cfg, common.ErrNotAllowed, false, false)
		return
	}
	activeUser, err := models.FindSessionOwner(h.db, cookie.Value)
	if err != nil || !activeUser.IsAdmin() {
		fmt.Println("Could not get cookie owner", err)
		fail.BadRequest(res, req, h.cfg, common.ErrNotAllowed, err == nil, false)
		return
	}
	// Extract the data submitted in the form.
	req.ParseForm()
	archiverID := req.FormValue("archiverID")
	id, parseErr := strconv.Atoi(archiverID)
	if parseErr != nil {
		fmt.Println("Invalid archiver ID", parseErr)
		fail.BadRequest(res, req, h.cfg, common.ErrGenericInvalidData, true, true)
		return
	}
	// Try to promote the archiver selected.
	archiver, findErr := models.FindArchiver(h.db, id)
	if findErr != nil {
		fmt.Println("No such archiver", id)
		fail.BadRequest(res, req, h.cfg, common.ErrDatabaseOperation, true, true)
		return
	}
	archiver.MakeAdmin(activeUser)
	updateErr := archiver.Update(h.db)
	if updateErr != nil {
		fmt.Println("Could not update archiver", updateErr)
		fail.InternalError(res, req, h.cfg, common.ErrDatabaseOperation, true, true)
		return
	}
	// Redirect back to the archivers list page.
	handler := NewListHandler(h.cfg, h.db)
	handler.PushSuccessMsg(fmt.Sprintf("Successfully made %s an administrator.", archiver.Email()))
	handler.ServeHTTP(res, req)
}
