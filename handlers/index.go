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

// The name of the index HTML template file to serve to users.
const indexPage string = "index.html"

// IndexHandler implements http.ServeHTTP to load and serve a simple index
// page to users.
type IndexHandler struct {
	cfg *config.Config
	db  *sql.DB
}

// NewIndexHandler is the constructor function for IndexHandler.
// Arguments:
// cfg: A reference to the application's global configuration.
func NewIndexHandler(cfg *config.Config, db *sql.DB) IndexHandler {
	return IndexHandler{
		cfg: cfg,
		db:  db,
	}
}

// ServeHTTP is implemented by IndexHandler to serve an index page to users.
// Arguments:
// res: Provided by the net/http server.
// req: Provided by the net/http server.
func (h IndexHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	loggedIn := false
	isAdmin := false
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err == nil {
		activeUser, err := models.FindSessionOwner(h.db, cookie.Value)
		if err == nil {
			loggedIn = true
			isAdmin = activeUser.IsAdmin()
		}
	}
	t, loadErr := template.ParseFiles(
		path.Join(h.cfg.TemplateDir, indexPage),
		path.Join(h.cfg.TemplateDir, headTemplate),
		path.Join(h.cfg.TemplateDir, navTemplate))
	if loadErr != nil {
		fmt.Println("failed to load template", loadErr)
		InternalError(res, req, h.cfg, errTemplateLoad, loggedIn, isAdmin)
		return
	}
	t.Execute(res, struct {
		LoggedIn    bool
		UserIsAdmin bool
		Successes   []string
	}{loggedIn, isAdmin, []string{}})
}
