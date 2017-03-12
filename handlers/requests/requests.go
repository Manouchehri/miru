package requests

import (
	"../../config"

	"github.com/gorilla/mux"

	"database/sql"
)

// RegisterHandlers registers request handlers to a subrouter.
func RegisterHandlers(r *mux.Router, cfg *config.Config, db *sql.DB) {
	r.Handle("/list", NewListHandler(cfg, db)).Methods("GET")
	r.Handle("/create", NewCreatePagePageHandler(cfg, db)).Methods("GET")
	r.Handle("/create", NewCreateHandler(cfg, db)).Methods("POST")
	r.Handle("/fulfill", NewFulfillPageHandler(cfg, db)).Methods("GET")
	r.Handle("/fulfill", NewFulfillHandler(cfg, db)).Methods("POST")
	r.Handle("/reject", NewRejectHandler(cfg, db)).Methods("POST")
}
