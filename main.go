package main

import (
	"./config"
	"./handlers"
	"./models"

	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config.MustLoad()

	db, dbErr := sql.Open("sqlite3", cfg.Database)
	if dbErr != nil {
		panic(dbErr)
	}
	initErr := models.InitializeTables(db)
	if initErr != nil {
		panic(initErr)
	}

	r := mux.NewRouter()
	index := handlers.NewIndexHandler(&cfg)
	r.Handle("/", index)
	http.Handle("/", r)
	fmt.Println("Listening on", cfg.BindAddress)
	http.ListenAndServe(cfg.BindAddress, nil)
}
