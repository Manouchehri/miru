package requests

import (
	"../../auth"
	"../../config"
	"../../models"
	"../common"
	"../fail"

	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path"
)

// requestListPage is the name of the HTML template that displays a list of pending
// monitor requests.
const requestListPage string = "requestlist.html"

// ListHandler implements net/http.ServeHTTP to serve a page showing all
// pending monitor requests to administrators.
type ListHandler struct {
	cfg       *config.Config
	db        *sql.DB
	Successes []string
}

// NewListHandler is the constructor function for a ListHandler.
// Arguments:
// cfg: The application's global configuration.
// db: A database connection.
// Returns:
// A new ListHandler that can be bound to a router.
func NewListHandler(cfg *config.Config, db *sql.DB) ListHandler {
	return ListHandler{
		cfg:       cfg,
		db:        db,
		Successes: []string{},
	}
}

// PushSuccessMsg adds a new message that will be displayed on the page served by the
// handler to indicate a successful operation.
// Arguments:
// msg: A success message to display to the user.
func (h *ListHandler) PushSuccessMsg(msg string) {
	h.Successes = append(h.Successes, msg)
}

// ServeHTTP serves a page for administrators to view pending monitor requests.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h ListHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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
	// Load all pending requests into an array of structs we can display in the page.
	requests, err := models.ListPendingRequests(h.db)
	if err != nil {
		fmt.Println("Could not get requests", err)
		fail.InternalError(res, req, h.cfg, common.ErrDatabaseOperation, true, true)
		return
	}
	type Data struct {
		MadeBy       string
		URL          string
		Instructions string
		CSRFToken    string
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
		csrfToken := models.GenerateAntiCSRFToken(h.db, auth.AntiCSRFTokenLength)
		pendingRequests = append(pendingRequests, Data{
			MadeBy:       madeBy,
			URL:          request.URL(),
			Instructions: request.Instructions(),
			CSRFToken:    csrfToken.Token(),
			RequestID:    request.ID(),
		})
	}
	// Serve the listing page.
	t, err := template.ParseFiles(
		path.Join(h.cfg.TemplateDir, requestListPage),
		path.Join(h.cfg.TemplateDir, common.HeadTemplate),
		path.Join(h.cfg.TemplateDir, common.NavTemplate))
	if err != nil {
		fmt.Println("Could not load template", err)
		fail.InternalError(res, req, h.cfg, common.ErrTemplateLoad, true, true)
		return
	}
	t.Execute(res, struct {
		Requests    []Data
		LoggedIn    bool
		UserIsAdmin bool
		Successes   []string
	}{pendingRequests, true, archiver.IsAdmin(), h.Successes})
}
