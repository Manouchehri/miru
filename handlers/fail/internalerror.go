package fail

import (
	"../../config"
	"../common"

	"html/template"
	"net/http"
	"path"
)

// InternalError is a simple net/http HandlerFunc that will write a simple error
// page with an error message.
func InternalError(res http.ResponseWriter, req *http.Request, cfg *config.Config, err error, loggedIn, isAdmin bool) {
	res.WriteHeader(http.StatusInternalServerError)
	t, _ := template.ParseFiles(
		path.Join(cfg.TemplateDir, errorTemplate),
		path.Join(cfg.TemplateDir, common.HeadTemplate),
		path.Join(cfg.TemplateDir, common.NavTemplate))
	t.Execute(res, struct {
		LoggedIn    bool
		UserIsAdmin bool
		Errors      []string
		Successes   []string
	}{loggedIn, isAdmin, []string{err.Error()}, []string{}})
}
