package index

import (
	"../../config"

	"github.com/gorilla/mux"

	"database/sql"
)

// RegisterHandlers registers request handlers to a subrouter.
func RegisterHandlers(r *mux.Router, cfg *config.Config, db *sql.DB) {
	r.Handle("/", NewFrontPageHandler(cfg, db)).Methods("GET")
}
