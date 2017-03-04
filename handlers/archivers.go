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
	"strconv"
)

// archiversPage is the name of the template HTML file that lists archivers.
const archiversPage string = "archivers.html"

// ArchiversListPageHandler implements net/http.ServeHTTP to serve a page containing a list
// of all archivers with buttons to make them an administrator.
type ArchiversListPageHandler struct {
	cfg *config.Config
	db  *sql.DB
}

// MakeAdminHandler handles requests from administrators to promote an archiver to
// have administrator privileges.
type MakeAdminHandler struct {
	db *sql.DB
}

// NewArchiversListPageHandler is the constructor function for a new ArchiversListPageHandler.
// Arguments:
// cfg: The application's global configuration.
// db: A database connection.
// Returns:
// A new ArchiversListPageHandler which can be bound to a router.
func NewArchiversListPageHandler(cfg *config.Config, db *sql.DB) ArchiversListPageHandler {
	return ArchiversListPageHandler{
		cfg: cfg,
		db:  db,
	}
}

// NewMakeAdminHandler is the constructor for a new MakeAdminHandler.
// Arguments:
// db: A database connection.
// Returns:
// A new MakeAdminHandler which can be bound to a router.
func NewMakeAdminHandler(db *sql.DB) MakeAdminHandler {
	return MakeAdminHandler{
		db: db,
	}
}

// ServeHTTP serves a page with a table of all archivers.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h ArchiversListPageHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Check that the request is coming from an authenticated administrator.
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		fmt.Println("Could not find cookie", err)
		BadRequest(res, req)
		return
	}
	activeUser, err := models.FindSessionOwner(h.db, cookie.Value)
	if err != nil || !activeUser.IsAdmin() {
		fmt.Println("Could not get cookie owner", err)
		BadRequest(res, req)
		return
	}
	// Load information about archivers.
	archivers, findErr := models.ListArchivers(h.db)
	if findErr != nil {
		fmt.Println("Could not get archivers", findErr)
		InternalError(res, req)
		return
	}
	type Data struct {
		ID      int
		Email   string
		IsAdmin bool
	}
	data := make([]Data, len(archivers))
	for _, archiver := range archivers {
		data = append(data, Data{
			ID:      archiver.ID(),
			Email:   archiver.Email(),
			IsAdmin: archiver.IsAdmin(),
		})
	}
	// Serve the page with the data about archivers.
	t, err := template.ParseFiles(
		path.Join(h.cfg.TemplateDir, archiversPage),
		path.Join(h.cfg.TemplateDir, headTemplate),
		path.Join(h.cfg.TemplateDir, navTemplate))
	if err != nil {
		fmt.Println("Error parsing archivers page template", err)
		InternalError(res, req)
		return
	}
	t.Execute(res, struct{ Archivers []Data }{data})
}

// ServeHTTP handles requests to have an archiver promoted to become an administrator.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h MakeAdminHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Check that the request is coming from an authenticated administrator.
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		fmt.Println("Could not find cookie", err)
		BadRequest(res, req)
		return
	}
	activeUser, err := models.FindSessionOwner(h.db, cookie.Value)
	if err != nil || !activeUser.IsAdmin() {
		fmt.Println("Could not get cookie owner", err)
		BadRequest(res, req)
		return
	}
	// Extract the data submitted in the form.
	req.ParseForm()
	archiverID := req.FormValue("archiverID")
	id, parseErr := strconv.Atoi(archiverID)
	if parseErr != nil {
		fmt.Println("Invalid archiver ID", parseErr)
		BadRequest(res, req)
		return
	}
	// Try to promote the archiver selected.
	archiver, findErr := models.FindArchiver(h.db, id)
	if findErr != nil {
		fmt.Println("No such archiver", id)
		BadRequest(res, req)
		return
	}
	archiver.MakeAdmin(activeUser)
	updateErr := archiver.Update(h.db)
	if updateErr != nil {
		fmt.Println("Could not update archiver", updateErr)
		InternalError(res, req)
		return
	}
	// Redirect back to the archivers list page.
	http.Redirect(res, req, "/archivers", http.StatusSeeOther)
}
