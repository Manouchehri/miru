package handlers

import (
	"net/http"
	"path"

	"../config"

	"database/sql"
	"html/template"
)

// uploadPage is the name of the HTML template containing a monitor script
// uploading form.
const uploadPage string = "scriptupload.html"

// UploadScriptHandler implements net/http.ServeHTTP to handle new monitor
// script uploads from administrators.
type UploadScriptHandler struct {
	cfg *config.Config
	db  *sql.DB
}

// UploadPageHandler implements net/http.ServeHTTP to sreve a page to
// administrators that they can use to upload new monitor scripts.
type UploadPageHandler struct {
	cfg *config.Config
}

// NewUploadScriptHandler is the constructor function for UploadScriptHandler.
// Arguments:
// c: A reference to the application's global configuration.
// db: A reference to a database connection.
func NewUploadScriptHandler(c *config.Config, db *sql.DB) UploadScriptHandler {
	return UploadScriptHandler{
		cfg: c,
		db:  db,
	}
}

// NewUploadPageHandler is the constructor function for UploadPageHandler.
// Arguments:
// cfg: A reference to the application's global configuration.
func NewUploadPageHandler(cfg *config.Config) UploadPageHandler {
	return UploadPageHandler{
		cfg: cfg,
	}
}

// ServeHTTP handles file uploads containing new monitor scripts.
func (h UploadScriptHandler) ServeHTTP(
	res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("hello world"))
}

// ServeHTTP serves a page that administrators can use to upload new
// monitor scripts through.
func (h UploadPageHandler) ServeHTTP(
	res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles(path.Join(h.cfg.TemplateDir, uploadPage))
	if err != nil {
		InternalError(res, req)
		return
	}
	t.Execute(res, nil)
}
