package archivers

import (
	"../../config"
	"../common"

	"html/template"
	"net/http"
	"path"
)

// loginPage is the name of the template file containing a login form.
const loginPage string = "login.html"

// LoginPageHandler implements net/http.ServeHTTP to serve a login page.
type LoginPageHandler struct {
	cfg *config.Config
}

// NewLoginPageHandler is the constructor function for a LoginPageHandler.
// Arguments:
// cfg: The application's global configuration.
// Returns:
// A new LoginPageHandler, which can be bound to a router.
func NewLoginPageHandler(cfg *config.Config) LoginPageHandler {
	return LoginPageHandler{
		cfg: cfg,
	}
}

// ServeHTTP writes the login page to the requester.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h LoginPageHandler) ServeHTTP(
	res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles(
		path.Join(h.cfg.TemplateDir, loginPage),
		path.Join(h.cfg.TemplateDir, common.HeadTemplate),
		path.Join(h.cfg.TemplateDir, common.NavTemplate))
	if err != nil {
		InternalError(res, req, h.cfg, errTemplateLoad, false, false)
		return
	}
	t.Execute(res, struct {
		LoggedIn    bool
		UserIsAdmin bool
		Successes   []string
	}{false, false, []string{}})
}
