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
	"net/url"
)

// CreateHandler implements net/http.ServeHTTP to handle requests to have a
// site monitored.
type CreateHandler struct {
	cfg *config.Config
	db  *sql.DB
}

// NewCreateHandler is the constructor function for a CreateHandler.
func NewCreateHandler(cfg *config.Config, db *sql.DB) CreateHandler {
	return CreateHandler{
		cfg: cfg,
		db:  db,
	}
}

// ServeHTTP handles a form upload containing a request to have a site monitored.
func (h CreateHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Check that the request is coming from an authenticated archiver.
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		fail.BadRequest(res, req, h.cfg, common.ErrNotAllowed, false, false)
		return
	}
	archiver, err := models.FindSessionOwner(h.db, cookie.Value)
	fmt.Println("Archiver making a request:", archiver)
	if err != nil {
		fail.BadRequest(res, req, h.cfg, common.ErrNotAllowed, false, false)
		return
	}
	// Extract data from the form.
	requstedURL := req.FormValue("url")
	instructions := req.FormValue("instructions")
	parsedURL, parseErr := url.Parse(requstedURL)
	if parseErr != nil {
		fail.BadRequest(res, req, h.cfg, errors.New("please specify a valid url to monitor"), true, archiver.IsAdmin())
		return
	}
	// Create a new request. We strip the query string from the URL since it:
	// 1. Shouldn't be necessary.
	// 2. Could contain sensitive info about the archiver.
	parsedURL.ForceQuery = false
	parsedURL.RawQuery = ""
	request := models.NewRequest(archiver, parsedURL.String(), instructions)
	saveErr := request.Save(h.db)
	if saveErr != nil {
		fail.InternalError(res, req, h.cfg, common.ErrDatabaseOperation, true, archiver.IsAdmin())
		return
	}
	fmt.Println("Created request", request)
	handler := NewCreatePageHandler(h.cfg, h.db)
	handler.PushSuccessMsg("Successfully sent your request. An administrator will review it soon.")
	handler.ServeHTTP(res, req)
}
