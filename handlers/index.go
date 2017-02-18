package handlers

import (
  "../config"

  "html/template"
  "net/http"
  "path"
)

// The name of the index HTML template file to serve to users.
const index_page string = "index.html"

// IndexHandler implements http.ServeHTTP to load and serve a simple index
// page to users.
type IndexHandler struct {
  cfg *config.Config
}

// NewIndexHandler is the constructor function for IndexHandler.
// Arguments:
// cfg: A reference to the application's global configuration.
func NewIndexHandler(cfg *config.Config) IndexHandler {
  return IndexHandler {
    cfg: cfg,
  }
}

// ServeHTTP is implemented by IndexHandler to serve an index page to users.
// Arguments:
// res: Provided by the net/http server.
// req: Provided by the net/http server.
func (h IndexHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  t, loadErr := template.ParseFiles(path.Join(h.cfg.TemplateDir, index_page))
  if loadErr != nil {
    InternalError(res, req)
    return
  }
  t.Execute(res, nil)
}
