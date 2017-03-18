package archivers

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

// loginPage is the name of the template file containing a login form.
const loginPage string = "login.html"

// LoginPageHandler implements net/http.ServeHTTP to serve a login page.
type LoginPageHandler struct {
	cfg *config.Config
	db  *sql.DB
}

// NewLoginPageHandler is the constructor function for a LoginPageHandler.
func NewLoginPageHandler(cfg *config.Config, db *sql.DB) LoginPageHandler {
	return LoginPageHandler{
		cfg: cfg,
		db:  db,
	}
}

// ServeHTTP writes the login page to the requester.
func (h LoginPageHandler) ServeHTTP(
	res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles(
		path.Join(h.cfg.TemplateDir, loginPage),
		path.Join(h.cfg.TemplateDir, common.HeadTemplate),
		path.Join(h.cfg.TemplateDir, common.NavTemplate))
	if err != nil {
		fail.InternalError(res, req, h.cfg, common.ErrTemplateLoad, false, false)
		return
	}
	csrfToken := models.GenerateAntiCSRFToken(h.db, auth.AntiCSRFTokenLength)
	saveErr := csrfToken.Save(h.db)
	if saveErr != nil {
		fail.InternalError(res, req, h.cfg, common.ErrDatabaseOperation, false, false)
		return
	}
	t.Execute(res, struct {
		LoggedIn    bool
		UserIsAdmin bool
		CSRFToken   string
		Successes   []string
	}{false, false, csrfToken.Token(), []string{}})
}
