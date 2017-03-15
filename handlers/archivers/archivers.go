package archivers

import (
	"../../config"

	"github.com/gorilla/mux"

	"database/sql"
)

// RegisterHandlers registers request handlers to a subrouter.
func RegisterHandlers(r *mux.Router, cfg *config.Config, db *sql.DB) {
	r.Handle("/list", NewListHandler(cfg, db)).Methods("GET")
	r.Handle("/login", NewLoginPageHandler(cfg, db)).Methods("GET")
	r.Handle("/login", NewLoginHandler(cfg, db)).Methods("POST")
	r.Handle("/logout", NewLogoutHandler(cfg, db)).Methods("GET")
	r.Handle("/register", NewRegisterPageHandler(cfg)).Methods("GET")
	r.Handle("/register", NewRegisterHandler(cfg, db)).Methods("POST")
	r.Handle("/promote", NewPromoteHandler(cfg, db)).Methods("POST")
}
