package handlers

import (
	"../config"

	"html/template"
	"net/http"
	"path"
)

// InternalError is a simple net/http HandlerFunc that will write a simple error
// page with an error message.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
// err: The error to write.
// loggedIn: True if the user is logged in.
// isAdmin: True if the user is an admin.
func InternalError(res http.ResponseWriter, req *http.Request, cfg *config.Config, err error, loggedIn, isAdmin bool) {
	res.WriteHeader(http.StatusInternalServerError)
	t, _ := template.ParseFiles(
		path.Join(cfg.TemplateDir, errorTemplate),
		path.Join(cfg.TemplateDir, headTemplate),
		path.Join(cfg.TemplateDir, navTemplate))
	t.Execute(res, struct {
		LoggedIn    bool
		UserIsAdmin bool
		Errors      []string
		Successes   []string
	}{loggedIn, isAdmin, []string{err.Error()}, []string{}})
}
