package handlers

import (
	"../auth"
	"../config"
	"../models"

	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path"
)

// adminPanelPage is the name of the template HTML file that links to
// important pages for admins.
const adminPanelPage string = "adminpanel.html"

// AdminPanelPageHandler is
type AdminPanelPageHandler struct {
	cfg *config.Config
	db  *sql.DB
}

// NewAdminPanelPageHandler is the constructor function for an
// Arguments:
// cfg: The application's global configuration.
// db: A database connection.
// Returns:
// A new AdminPanelPageHandler that can be bound to a router.
func NewAdminPanelPageHandler(cfg *config.Config, db *sql.DB) AdminPanelPageHandler {
	return AdminPanelPageHandler{
		cfg: cfg,
		db:  db,
	}
}

func (h AdminPanelPageHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Check that the request is coming from an authenticated administrator.
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		fmt.Println("Could not find cookie", err)
		BadRequest(res, req, h.cfg, errNotAllowed, false, false)
		return
	}
	activeUser, err := models.FindSessionOwner(h.db, cookie.Value)
	if err != nil || !activeUser.IsAdmin() {
		fmt.Println("Could not get cookie owner", err)
		BadRequest(res, req, h.cfg, errNotAllowed, err == nil, false)
		return
	}
	// Serve the admin panel page.
	t, err := template.ParseFiles(
		path.Join(h.cfg.TemplateDir, adminPanelPage),
		path.Join(h.cfg.TemplateDir, headTemplate),
		path.Join(h.cfg.TemplateDir, navTemplate))
	if err != nil {
		fmt.Println("Failed to parse templates", err)
		InternalError(res, req, h.cfg, errTemplateLoad, true, true)
		return
	}
	t.Execute(res, struct {
		UserIsAdmin bool
		LoggedIn    bool
		Successes   []string
	}{true, true, []string{}})
}
