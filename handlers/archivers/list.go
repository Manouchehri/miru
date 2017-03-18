package archivers

import (
	"../../auth"
	"../../config"
	"../../models"
	"../common"
	"../fail"

	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path"
)

// archiversPage is the name of the template HTML file that lists archivers.
const archiversPage string = "archivers.html"

// ListHandler implements net/http.ServeHTTP to serve a page containing a list
// of all archivers with buttons to make them an administrator.
type ListHandler struct {
	cfg       *config.Config
	db        *sql.DB
	Successes []string
}

// NewListHandler is the constructor function for a new ListHandler.
func NewListHandler(cfg *config.Config, db *sql.DB) ListHandler {
	return ListHandler{
		cfg:       cfg,
		db:        db,
		Successes: []string{},
	}
}

// PushSuccessMsg adds a new message that will be displayed on the page served by the
// handler to indicate a successful operation.
func (h *ListHandler) PushSuccessMsg(msg string) {
	h.Successes = append(h.Successes, msg)
}

// ServeHTTP serves a page with a table of all archivers.
func (h ListHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Check that the request is coming from an authenticated administrator.
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		fmt.Println("Could not find cookie", err)
		fail.BadRequest(res, req, h.cfg, common.ErrNotAllowed, false, false)
		return
	}
	fmt.Println("Found cookie", cookie.Value)
	activeUser, err := models.FindSessionOwner(h.db, cookie.Value)
	if err != nil || !activeUser.IsAdmin() {
		fmt.Println("Could not get cookie owner", err)
		fail.BadRequest(res, req, h.cfg, common.ErrNotAllowed, err == nil, false)
		return
	}
	// Load information about archivers.
	archivers, findErr := models.ListArchivers(h.db)
	if findErr != nil {
		fmt.Println("Could not get archivers", findErr)
		fail.InternalError(res, req, h.cfg, common.ErrDatabaseOperation, true, true)
		return
	}
	type Data struct {
		ID        int
		Email     string
		CSRFToken string
		IsAdmin   bool
	}
	data := []Data{}
	for _, archiver := range archivers {
		csrfToken := models.GenerateAntiCSRFToken(h.db, auth.AntiCSRFTokenLength)
		saveErr := csrfToken.Save(h.db)
		if saveErr != nil {
			fail.InternalError(res, req, h.cfg, common.ErrDatabaseOperation, true, true)
			return
		}
		data = append(data, Data{
			ID:        archiver.ID(),
			Email:     archiver.Email(),
			CSRFToken: csrfToken.Token(),
			IsAdmin:   archiver.IsAdmin(),
		})
	}
	// Serve the page with the data about archivers.
	t, err := template.ParseFiles(
		path.Join(h.cfg.TemplateDir, archiversPage),
		path.Join(h.cfg.TemplateDir, common.HeadTemplate),
		path.Join(h.cfg.TemplateDir, common.NavTemplate))
	if err != nil {
		fmt.Println("Error parsing archivers page template", err)
		fail.InternalError(res, req, h.cfg, common.ErrTemplateLoad, true, true)
		return
	}
	t.Execute(res, struct {
		Archivers   []Data
		LoggedIn    bool
		UserIsAdmin bool
		Successes   []string
	}{data, true, true, h.Successes})
}
