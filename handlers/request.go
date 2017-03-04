package handlers

import (
	"../auth"
	"../config"
	"../models"

	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
)

// requestPage is the name of the HTML template containing the form for archivers to
// request a site be monitored through.
const requestPage string = "request.html"

// requestListPage is the name of the HTML template that displays a list of pending
// monitor requests.
const requestListPage string = "requestlist.html"

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

// ListRequestsHandler implements net/http.ServeHTTP to serve a page showing all
// pending monitor requests to administrators.
type ListRequestsHandler struct {
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

// ListRequestsHandler is the constructor function for a ListRequestsHandler.
// Arguments:
// cfg: The application's global configuration.
// db: A database connection.
// Returns:
// A new ListRequestsHandler that can be bound to a router.
func NewListRequestsHandler(cfg *config.Config, db *sql.DB) ListRequestsHandler {
	return ListRequestsHandler{
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
	t, err := template.ParseFiles(
		path.Join(h.cfg.TemplateDir, requestPage),
		path.Join(h.cfg.TemplateDir, headTemplate),
		path.Join(h.cfg.TemplateDir, navTemplate))
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
	fmt.Println("Archiver making a request:", archiver)
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
	fmt.Println("Created request", request)
	http.Redirect(res, req, "/request", http.StatusFound)
}

// ServeHTTP serves a page for administrators to view pending monitor requests.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h ListRequestsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Check that the request is coming from an authenticated administrator.
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		fmt.Println("No cookie", err)
		BadRequest(res, req)
		return
	}
	archiver, err := models.FindSessionOwner(h.db, cookie.Value)
	if err != nil || !archiver.IsAdmin() {
		fmt.Println("Not admin", err)
		BadRequest(res, req)
		return
	}
	// Load all pending requests into an array of structs we can display in the page.
	requests, err := models.ListPendingRequests(h.db)
	if err != nil {
		fmt.Println("Could not get requests", err)
		InternalError(res, req)
		return
	}
	type Data struct {
		MadeBy       string
		URL          string
		Instructions string
		RequestID    int
	}
	pendingRequests := []Data{}
	for _, request := range requests {
		fmt.Println("found request", request)
		reqCreator, err := models.FindArchiver(h.db, request.Creator())
		fmt.Println("created by", reqCreator)
		madeBy := "deleted"
		if err == nil {
			madeBy = reqCreator.Email()
		} else {
			fmt.Println("Error finding request creator", err)
		}
		pendingRequests = append(pendingRequests, Data{
			MadeBy:       madeBy,
			URL:          request.URL(),
			Instructions: request.Instructions(),
			RequestID:    request.ID(),
		})
	}
	// Serve the listing page.
	t, err := template.ParseFiles(
		path.Join(h.cfg.TemplateDir, requestListPage),
		path.Join(h.cfg.TemplateDir, headTemplate),
		path.Join(h.cfg.TemplateDir, navTemplate))
	if err != nil {
		fmt.Println("Could not load template", err)
		InternalError(res, req)
		return
	}
	t.Execute(res, struct{ Requests []Data }{pendingRequests})
}
