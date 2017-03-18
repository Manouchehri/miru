package fail

import (
	"../../config"
	"../common"

	"html/template"
	"net/http"
	"path"
)

// BadRequest is a simple net/http HandlerFunc that will write an error
// message to users if something is wrong with a request.
func BadRequest(res http.ResponseWriter, req *http.Request, cfg *config.Config, err error, loggedIn, isAdmin bool) {
	res.WriteHeader(http.StatusBadRequest)
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
