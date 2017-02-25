package handlers

import (
	"../auth"
	"../config"
	"../models"

	"database/sql"
	"html/template"
	"net/http"
	"path"
)

// requestPage is the name of the HTML template containing the form for archivers to
// request a site be monitored through.
const requestPage string = "request.html"

// MakeRequestPageHandler is implements net/http.ServeHTTP to serve the request page.
type MakeRequestPageHandler struct {
	cfg *config.Config
	db  *sql.DB
}

// NewMakeRequestPageHandler is the constructor function for a MakeRequestPageHandler.
// Arguments:
// cfg: The application's global configuration.
// db: A database connection.
// Returns:
// A new MakeRequestPageHandler that can be bound to a router.
func NewMakeRequestPageHandler(cfg *config.Config, db *sql.DB) MakeRequestPageHandler {
	return MakeRequestPageHandler{
		cfg: cfg,
		db:  db,
	}
}

// ServeHTTP serves a page that archivers can use to make requests to have a site monitored.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h MakeRequestPageHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Check that the request is coming from an authenticated archiver.
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		BadRequest(res, req)
		return
	}
	_, err = models.FindSessionOwner(h.db, cookie.Value)
	if err != nil {
		BadRequest(res, req)
		return
	}
	// Serve the page.
	t, err := template.ParseFiles(path.Join(h.cfg.TemplateDir, requestPage))
	if err != nil {
		InternalError(res, req)
		return
	}
	t.Execute(res, nil)
}
