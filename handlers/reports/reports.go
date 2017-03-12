package reports

import (
	"../../config"

	"github.com/gorilla/mux"

	"database/sql"
)

// RegisterHandlers registers request handlers to a subrouter.
func RegisterHandlers(r *mux.Router, cfg *config.Config, db *sql.DB) {
	r.Handle("/list", NewListHandler(cfg, db)).Methods("GET")
}
