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
	index := handlers.NewIndexHandler(&cfg)
	monitorPage := handlers.NewUploadPageHandler(&cfg)
	uploadScript := handlers.NewUploadScriptHandler(&cfg, db)
	registerPage := handlers.NewRegisterPageHandler(&cfg)
	register := handlers.NewRegisterHandler(&cfg, db)
	loginPage := handlers.NewLoginPageHandler(&cfg)
	login := handlers.NewLoginHandler(&cfg, db)
	requestPage := handlers.NewMakeRequestPageHandler(&cfg, db)
	request := handlers.NewMakeRequestHandler(&cfg, db)
	listRequests := handlers.NewListRequestsHandler(&cfg, db)
	reports := handlers.NewReportPageHandler(&cfg, db)
	listArchivers := handlers.NewArchiversListPageHandler(&cfg, db)
	makeAdmin := handlers.NewMakeAdminHandler(db)
	r.Handle("/", index)
	r.Handle("/monitor", monitorPage).Methods("GET")
	r.Handle("/monitor", uploadScript).Methods("POST")
	r.Handle("/register", registerPage).Methods("GET")
	r.Handle("/register", register).Methods("POST")
	r.Handle("/login", loginPage).Methods("GET")
	r.Handle("/login", login).Methods("POST")
	r.Handle("/request", requestPage).Methods("GET")
	r.Handle("/request", request).Methods("POST")
	r.Handle("/listrequests", listRequests).Methods("GET")
	r.Handle("/reports", reports).Methods("GET")
	r.Handle("/archivers", listArchivers).Methods("GET")
	r.Handle("/makeadmin", makeAdmin).Methods("POST")

	r.PathPrefix("/js/").Handler(
		http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	r.PathPrefix("/css/").Handler(
		http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/", r)
	fmt.Println("Listening on", cfg.BindAddress)
	http.ListenAndServe(cfg.BindAddress, nil)
}
