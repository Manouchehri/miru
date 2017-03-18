package archivers

import (
	"../../config"
	"../common"
	"../fail"

	"html/template"
	"net/http"
	"path"
)

// registerPage is the name of the template file containing the register page.
const registerPage string = "register.html"

// RegisterPageHandler implements net/http.ServeHTTP to handle GET requests
// made by users wishing to view the register page.
type RegisterPageHandler struct {
	cfg *config.Config
}

// NewRegisterPageHandler is the constructor function for a new
// RegisterPageHandler.
func NewRegisterPageHandler(cfg *config.Config) RegisterPageHandler {
	return RegisterPageHandler{
		cfg: cfg,
	}
}

// ServeHTTP writes the register page to the requester.
func (h RegisterPageHandler) ServeHTTP(
	res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles(
		path.Join(h.cfg.TemplateDir, registerPage),
		path.Join(h.cfg.TemplateDir, common.HeadTemplate),
		path.Join(h.cfg.TemplateDir, common.NavTemplate))
	if err != nil {
		fail.InternalError(res, req, h.cfg, common.ErrTemplateLoad, false, false)
		return
	}
	t.Execute(res, struct {
		LoggedIn    bool
		UserIsAdmin bool
		Successes   []string
	}{false, false, []string{}})
}
