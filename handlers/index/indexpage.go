package index

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

// The name of the index HTML template file to serve to users.
const indexPage string = "index.html"

// FrontPageHandler implements http.ServeHTTP to load and serve a simple index
// page to users.
type FrontPageHandler struct {
	cfg       *config.Config
	db        *sql.DB
	Successes []string
}

// NewFrontPageHandler is the constructor function for FrontPageHandler.
// Arguments:
// cfg: A reference to the application's global configuration.
func NewFrontPageHandler(cfg *config.Config, db *sql.DB) FrontPageHandler {
	return FrontPageHandler{
		cfg:       cfg,
		db:        db,
		Successes: []string{},
	}
}

// PushSuccessMsg adds a new message that will be displayed on the page served by the
// handler to indicate a successful operation.
// Arguments:
// msg: A success message to display to the user.
func (h *FrontPageHandler) PushSuccessMsg(msg string) {
	h.Successes = append(h.Successes, msg)
}

// ServeHTTP is implemented by FrontPageHandler to serve an index page to users.
// Arguments:
// res: Provided by the net/http server.
// req: Provided by the net/http server.
func (h FrontPageHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	loggedIn := false
	isAdmin := false
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err == nil {
		activeUser, err := models.FindSessionOwner(h.db, cookie.Value)
		if err == nil {
			fmt.Println("Found session owner", activeUser)
			loggedIn = true
			isAdmin = activeUser.IsAdmin()
		}
	}
	t, loadErr := template.ParseFiles(
		path.Join(h.cfg.TemplateDir, indexPage),
		path.Join(h.cfg.TemplateDir, common.HeadTemplate),
		path.Join(h.cfg.TemplateDir, common.NavTemplate))
	if loadErr != nil {
		fmt.Println("failed to load template", loadErr)
		fail.InternalError(res, req, h.cfg, common.ErrTemplateLoad, loggedIn, isAdmin)
		return
	}
	t.Execute(res, struct {
		LoggedIn    bool
		UserIsAdmin bool
		Successes   []string
	}{loggedIn, isAdmin, h.Successes})
}
