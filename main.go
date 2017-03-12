package main

import (
	"./config"
	"./handlers"
	"./models"
	"./tasks"

	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	// Start the task runner so that it will periodically run a monitor script
	// to check for changes to sites, and shut everything down if a terminate
	// signal is sent by the user.
	errors := make(chan error)
	terminate := make(chan bool, 2)
	go tasks.RunMonitors(db, 1*time.Second, errors, terminate)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill)
	go func() {
		<-signals
		// Write one terminate signal for RunMonitors, and another for
		// the error handling code below.
		terminate <- true
		terminate <- true
	}()

	// Read any errors encountered trying to run monitor scripts.
	go func() {
		for {
			select {
			case err := <-errors:
				fmt.Println("[---] Error: ", err.Error())
			case <-terminate:
				break
			}
		}
	}()

	r := mux.NewRouter()
	handlers.RegisterHandlers(r, &cfg, db)
	r.PathPrefix("/js/").Handler(
		http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	r.PathPrefix("/css/").Handler(
		http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/", r)
	fmt.Println("Listening on", cfg.BindAddress)
	http.ListenAndServe(cfg.BindAddress, nil)
}
