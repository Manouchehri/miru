package main

import (
  "./config"
  "./handlers"

  "fmt"
  "net/http"

  "github.com/gorilla/mux"
)

func main() {
  cfg := config.MustLoad()
  r := mux.NewRouter()
  index  := handlers.NewIndexHandler(&cfg)
  r.Handle("/", index)
  http.Handle("/", r)
  fmt.Println("Listening on", cfg.BindAddress)
  http.ListenAndServe(cfg.BindAddress, nil)
}
