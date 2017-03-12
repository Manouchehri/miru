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
