package admin

import (
	"../"
	"../../auth"
	"../../config"
	"../../models"
	"../fail"

	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path"
)

// adminPanelPage is the name of the template HTML file that links to
// important pages for admins.
const adminPanelPage string = "adminpanel.html"

// PanelPageHandler implements net/http.ServeHTTP to serve administrators a
// panel page that links them to other pages where they can carry out
// administrative tasks.
type PanelPageHandler struct {
	cfg *config.Config
	db  *sql.DB
}

// NewPanelPageHandler is the constructor function for an
// Arguments:
// cfg: The application's global configuration.
// db: A database connection.
// Returns:
// A new PanelPageHandler that can be bound to a router.
func NewPanelPageHandler(cfg *config.Config, db *sql.DB) PanelPageHandler {
	return PanelPageHandler{
		cfg: cfg,
		db:  db,
	}
}

func (h PanelPageHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Check that the request is coming from an authenticated administrator.
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		fmt.Println("Could not find cookie", err)
		fail.BadRequest(res, req, h.cfg, errNotAllowed, false, false)
		return
	}
	activeUser, err := models.FindSessionOwner(h.db, cookie.Value)
	if err != nil || !activeUser.IsAdmin() {
		fmt.Println("Could not get cookie owner", err)
		fail.BadRequest(res, req, h.cfg, errNotAllowed, err == nil, false)
		return
	}
	// Serve the admin panel page.
	t, err := template.ParseFiles(
		path.Join(h.cfg.TemplateDir, adminPanelPage),
		path.Join(h.cfg.TemplateDir, handlers.HeadTemplate),
		path.Join(h.cfg.TemplateDir, handlers.NavTemplate))
	if err != nil {
		fmt.Println("Failed to parse templates", err)
		fail.InternalError(res, req, h.cfg, errTemplateLoad, true, true)
		return
	}
	t.Execute(res, struct {
		UserIsAdmin bool
		LoggedIn    bool
		Successes   []string
	}{true, true, []string{}})
}
