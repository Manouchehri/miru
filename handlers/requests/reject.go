package requests

import (
	"../"
	"../../auth"
	"../../config"
	"../../models"
	"../fail"

	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// RejectHandler implements net/http.ServeHTTP to handle the rejection
// of monitor requests by administrators.
type RejectHandler struct {
	cfg *config.Config
	db  *sql.DB
}

// NewRejectHandler is the constructor function for a RejectHandler.
// Arguments:
// cfg: The application's global configuration.
// db: A database connection.
// Returns:
// A new RejectHandler that can be bound to a router.
func NewRejectHandler(cfg *config.Config, db *sql.DB) RejectHandler {
	return RejectHandler{
		cfg: cfg,
		db:  db,
	}
}

// ServeHTTP deletes a pending request.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h RejectHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Check that the request is coming from an authenticated administrator.
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		fmt.Println("No cookie", err)
		fail.BadRequest(res, req, h.cfg, handlers.ErrNotAllowed, false, false)
		return
	}
	archiver, err := models.FindSessionOwner(h.db, cookie.Value)
	if err != nil || !archiver.IsAdmin() {
		fmt.Println("Not admin", err)
		fail.BadRequest(res, req, h.cfg, handlers.ErrNotAllowed, err == nil, false)
		return
	}
	// Extract inputs from the submitted form.
	req.ParseForm()
	requestID := req.FormValue("requestID")
	id, parseErr := strconv.Atoi(requestID)
	if parseErr != nil {
		fail.BadRequest(res, req, h.cfg, handlers.ErrGenericInvalidData, true, true)
		return
	}
	request, findErr := models.FindRequest(h.db, id)
	if findErr != nil {
		fail.BadRequest(res, req, h.cfg, errors.New("no such request"), true, true)
		return
	}
	deleteErr := request.Delete(h.db)
	if deleteErr != nil {
		fail.InternalError(res, req, h.cfg, handlers.ErrDatabaseOperation, true, true)
		return
	}
	handler := NewListRequestsHandler(h.cfg, h.db)
	handler.PushSuccessMsg(fmt.Sprintf("Successfully rejected request with ID %d", id))
	handler.ServeHTTP(res, req)
}
