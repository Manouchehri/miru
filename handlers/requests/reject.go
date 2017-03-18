package requests

import (
	"../../auth"
	"../../config"
	"../../models"
	"../common"
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
func NewRejectHandler(cfg *config.Config, db *sql.DB) RejectHandler {
	return RejectHandler{
		cfg: cfg,
		db:  db,
	}
}

// ServeHTTP deletes a pending request.
func (h RejectHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Check that the request is coming from an authenticated administrator.
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		fmt.Println("No cookie", err)
		fail.BadRequest(res, req, h.cfg, common.ErrNotAllowed, false, false)
		return
	}
	archiver, err := models.FindSessionOwner(h.db, cookie.Value)
	if err != nil || !archiver.IsAdmin() {
		fmt.Println("Not admin", err)
		fail.BadRequest(res, req, h.cfg, common.ErrNotAllowed, err == nil, false)
		return
	}
	// Extract inputs from the submitted form.
	req.ParseForm()
	csrfToken := req.FormValue("csrfToken")
	fmt.Println("Got ANTI-CSRF token", csrfToken)
	if !models.VerifyAndDeleteAntiCSRFToken(h.db, csrfToken) {
		fail.BadRequest(res, req, h.cfg, common.ErrNotAllowed, false, false)
		return
	}
	requestID := req.FormValue("requestID")
	id, parseErr := strconv.Atoi(requestID)
	if parseErr != nil {
		fail.BadRequest(res, req, h.cfg, common.ErrGenericInvalidData, true, true)
		return
	}
	request, findErr := models.FindRequest(h.db, id)
	if findErr != nil {
		fail.BadRequest(res, req, h.cfg, errors.New("no such request"), true, true)
		return
	}
	deleteErr := request.Delete(h.db)
	if deleteErr != nil {
		fail.InternalError(res, req, h.cfg, common.ErrDatabaseOperation, true, true)
		return
	}
	handler := NewListHandler(h.cfg, h.db)
	handler.PushSuccessMsg(fmt.Sprintf("Successfully rejected request with ID %d", id))
	handler.ServeHTTP(res, req)
}
