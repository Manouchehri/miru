package handlers

import (
	"../config"
	"./admin"
	"./archivers"
	"./index"
	"./reports"
	"./requests"

	"github.com/gorilla/mux"

	"database/sql"
	"errors"
)

// HeadTemplate is the name of the template file that contains the HTML
// head contents for all pages.
const HeadTemplate string = "head.html"

// NavTemplate is the name of the template file that contains the HTML
// navigation contents for all pages.
const NavTemplate string = "nav.html"

// Common errors containing messages that are safe to show the user.
var (
	ErrTemplateLoad       = errors.New("failed to load a page template")
	ErrInvalidCredentials = errors.New("the provided credentials are invalid")
	ErrDatabaseOperation  = errors.New("an internal database error occurred")
	ErrNotAllowed         = errors.New("you are not allowed to do that")
	ErrGenericInvalidData = errors.New("some of the input provided is invalid")
	ErrCreateFile         = errors.New("could not create a file for the monitor script")
	ErrBadPassword        = errors.New("password and repeated password must match and contain " +
		"at least one lowercase and uppercase letter, symbol, and number")
	ErrInvalidEmail = errors.New("invalid email address")
)

// RegisterHandlers registers all of our request handlers.
func RegisterHandlers(r *mux.Router, cfg *config.Config, db *sql.DB) {
	adminRouter := r.PathPrefix("/admin").Subrouter()
	archiversRouter := r.PathPrefix("/archivers").Subrouter()
	indexRouter := r.PathPrefix("/").Subrouter()
	reportsRouter := r.PathPrefix("/reports").Subrouter()
	requestsRouter := r.PathPrefix("/requests").Subrouter()
	admin.RegisterHandlers(adminRouter, cfg, db)
	archivers.RegisterHandlers(archiversRouter, cfg, db)
	index.RegisterHandlers(indexRouter, cfg, db)
	reports.RegisterHandlers(reportsRouter, cfg, db)
	requests.RegisterHandlers(requestsRouter, cfg, db)
}
