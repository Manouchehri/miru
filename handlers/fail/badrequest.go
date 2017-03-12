package fail

import (
	"../"
	"../../config"

	"html/template"
	"net/http"
	"path"
)

// BadRequest is a simple net/http HandlerFunc that will write an error
// message to users if something is wrong with a request.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
// cfg: The application's global configuration.
// err: The error to write.
// loggedIn: True if the user is logged in.
// isAdmin: True if the user is an admin.
func BadRequest(res http.ResponseWriter, req *http.Request, cfg *config.Config, err error, loggedIn, isAdmin bool) {
	res.WriteHeader(http.StatusBadRequest)
	t, _ := template.ParseFiles(
		path.Join(cfg.TemplateDir, errorTemplate),
		path.Join(cfg.TemplateDir, handlers.HeadTemplate),
		path.Join(cfg.TemplateDir, handlers.NavTemplate))
	t.Execute(res, struct {
		LoggedIn    bool
		UserIsAdmin bool
		Errors      []string
		Successes   []string
	}{loggedIn, isAdmin, []string{err.Error()}, []string{}})
}
