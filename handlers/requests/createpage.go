package requests

import (
	"../../auth"
	"../../config"
	"../../models"
	"../common"
	"../fail"

	"database/sql"
	"html/template"
	"net/http"
	"path"
)

// requestPage is the name of the HTML template containing the form for archivers to
// request a site be monitored through.
const requestPage string = "request.html"

// CreatePagePageHandler implements net/http.ServeHTTP to serve the request page.
type CreatePagePageHandler struct {
	cfg       *config.Config
	db        *sql.DB
	Successes []string
}

// NewCreatePagePageHandler is the constructor function for a CreatePagePageHandler.
// Arguments:
// cfg: The application's global configuration.
// db: A database connection.
// Returns:
// A new CreatePagePageHandler that can be bound to a router.
func NewCreatePagePageHandler(cfg *config.Config, db *sql.DB) CreatePagePageHandler {
	return CreatePagePageHandler{
		cfg:       cfg,
		db:        db,
		Successes: []string{},
	}
}

// PushSuccessMsg adds a new message that will be displayed on the page served by the
// handler to indicate a successful operation.
// Arguments:
// msg: A success message to display to the user.
func (h *CreatePagePageHandler) PushSuccessMsg(msg string) {
	h.Successes = append(h.Successes, msg)
}

// ServeHTTP serves a page that archivers can use to make requests to have a site monitored.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h CreatePagePageHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Check that the request is coming from an authenticated archiver.
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		fail.BadRequest(res, req, h.cfg, common.ErrNotAllowed, false, false)
		return
	}
	activeUser, err := models.FindSessionOwner(h.db, cookie.Value)
	if err != nil {
		fail.BadRequest(res, req, h.cfg, common.ErrNotAllowed, false, false)
		return
	}
	// Serve the page.
	t, err := template.ParseFiles(
		path.Join(h.cfg.TemplateDir, requestPage),
		path.Join(h.cfg.TemplateDir, common.HeadTemplate),
		path.Join(h.cfg.TemplateDir, common.NavTemplate))
	if err != nil {
		fail.InternalError(res, req, h.cfg, common.ErrTemplateLoad, true, activeUser.IsAdmin())
		return
	}
	t.Execute(res, struct {
		LoggedIn    bool
		UserIsAdmin bool
		Successes   []string
	}{true, activeUser.IsAdmin(), h.Successes})
}
