package handlers

import (
	"../auth"
	"../config"
	"../models"

	"database/sql"
	"html/template"
	"net/http"
	"net/url"
	"path"
)

// requestPage is the name of the HTML template containing the form for archivers to
// request a site be monitored through.
const requestPage string = "request.html"

// MakeRequestPageHandler implements net/http.ServeHTTP to serve the request page.
type MakeRequestPageHandler struct {
	cfg *config.Config
	db  *sql.DB
}

// MakeRequestHandler implements net/http.ServeHTTP to handle requests to have a
// site monitored.
type MakeRequestHandler struct {
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

// NewMakeRequestHandler is the constructor function for a MakeRequestHandler.
// Arguments:
// cfg: The application's global configuration.
// db: A database connection.
// Returns:
// A new MakeRequestHandler that can be bound to a router.
func NewMakeRequestHandler(cfg *config.Config, db *sql.DB) MakeRequestHandler {
	return MakeRequestHandler{
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

// ServeHTTP handles a form upload containing a request to have a site monitored.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h MakeRequestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Check that the request is coming from an authenticated archiver.
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		BadRequest(res, req)
		return
	}
	archiver, err := models.FindSessionOwner(h.db, cookie.Value)
	if err != nil {
		BadRequest(res, req)
		return
	}
	// Extract data from the form.
	requstedURL := req.FormValue("url")
	instructions := req.FormValue("instructions")
	parsedURL, parseErr := url.Parse(requstedURL)
	if parseErr != nil {
		BadRequest(res, req)
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
		InternalError(res, req)
		return
	}
	http.Redirect(res, req, "/request", http.StatusFound)
}
