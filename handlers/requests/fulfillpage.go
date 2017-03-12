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
	"html/template"
	"net/http"
	"path"
	"strconv"
)

// uploadPage is the name of the HTML template containing a monitor script
// uploading form.
const uploadPage string = "scriptupload.html"

// FulfillPageHandler implements net/http.ServeHTTP to sreve a page to
// administrators that they can use to upload new monitor scripts.
type FulfillPageHandler struct {
	cfg       *config.Config
	db        *sql.DB
	Successes []string
}

// NewFulfillPageHandler is the constructor function for FulfillPageHandler.
// Arguments:
// cfg: A reference to the application's global configuration.
// Returns:
// A new FulfillPageHandler that can be bound to a router.
func NewFulfillPageHandler(cfg *config.Config, db *sql.DB) FulfillPageHandler {
	return FulfillPageHandler{
		cfg:       cfg,
		db:        db,
		Successes: []string{},
	}
}

// PushSuccessMsg adds a new message that will be displayed on the page served by the
// handler to indicate a successful operation.
// Arguments:
// msg: A success message to display to the user.
func (h *FulfillPageHandler) PushSuccessMsg(msg string) {
	h.Successes = append(h.Successes, msg)
}

// ServeHTTP serves a page that administrators can use to upload new
// monitor scripts through.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h FulfillPageHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Check that the request is coming from an authenticated archiver.
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		fail.BadRequest(res, req, h.cfg, handlers.ErrNotAllowed, false, false)
		return
	}
	activeUser, err := models.FindSessionOwner(h.db, cookie.Value)
	if err != nil {
		fail.BadRequest(res, req, h.cfg, handlers.ErrNotAllowed, false, false)
		return
	}
	requestIDs, found := req.URL.Query()["id"]
	if !found || len(requestIDs) == 0 {
		fmt.Println("Need a request id")
		fail.BadRequest(res, req, h.cfg, errors.New("missing request id url parameter"), true, activeUser.IsAdmin())
		return
	}
	requestID, parseErr := strconv.Atoi(requestIDs[0])
	if parseErr != nil {
		fmt.Println("Need a request valid id")
		fail.BadRequest(res, req, h.cfg, handlers.ErrGenericInvalidData, true, activeUser.IsAdmin())
		return
	}
	t, err := template.ParseFiles(
		path.Join(h.cfg.TemplateDir, uploadPage),
		path.Join(h.cfg.TemplateDir, handlers.HeadTemplate),
		path.Join(h.cfg.TemplateDir, handlers.NavTemplate))
	if err != nil {
		fail.InternalError(res, req, h.cfg, handlers.ErrTemplateLoad, true, true)
		return
	}
	t.Execute(res, struct {
		CreatedFor  int
		LoggedIn    bool
		UserIsAdmin bool
		Successes   []string
	}{requestID, true, activeUser.IsAdmin(), h.Successes})
}
