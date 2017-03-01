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

// archiversPage is the name of the template HTML file that lists archivers.
const archiversPage string = "archivers.html"

//  ArchiversListPageHandler implements net/http.ServeHTTP to serve a page containing a list
// of all archivers with buttons to make them an administrator.
type ArchiversListPageHandler struct {
	cfg *config.Config
	db  *sql.DB
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
	t, err := template.ParseFiles(path.Join(h.cfg.TemplateDir, archiversPage))
	if err != nil {
		fmt.Println("Error parsing archivers page template", err)
		InternalError(res, req)
		return
	}
	t.Execute(res, struct{ Archivers []Data }{data})
}
