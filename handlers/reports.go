package handlers

import (
	"../auth"
	"../config"
	"../models"

	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"time"
)

// reportsPage is the name of the template HTML file with a table for monitor
// and report information to be displayed to administrators.
const reportsPage string = "reports.html"

// ReportPageHandler implements net/http.ServeHTTP to serve a page to
// administrators containing information about monitors that miru is
// running and data the scripts are reporting.
type ReportPageHandler struct {
	cfg *config.Config
	db  *sql.DB
}

// NewReportPageHandler is the constructor function for a ReportPageHandler.
// Arguments:
// cfg: The application's global configuration.
// db: A database connection.
// Returns:
// A new ReportPageHandler that can be bound to a router.
func NewReportPageHandler(cfg *config.Config, db *sql.DB) ReportPageHandler {
	return ReportPageHandler{
		cfg: cfg,
		db:  db,
	}
}

// ServeHTTP serves the reports page to an administrator.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h ReportPageHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Check that the request is coming from an authenticated administrator.
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		fmt.Println("Could not find cookie", err)
		BadRequest(res, req)
		return
	}
	activeUser, err := models.FindSessionOwner(h.db, cookie.Value)
	if err != nil || !activeUser.IsAdmin() {
		fmt.Println("Could not get cookie owner", err)
		BadRequest(res, req)
		return
	}
	// Load information about existing monitors and the last report each generated.
	monitors, findErr := models.ListMonitors(h.db)
	if findErr != nil {
		fmt.Println("Could not get monitors", findErr)
		InternalError(res, req)
		return
	}
	type Data struct {
		MonitorID          int
		ScriptPath         string
		LastRan            time.Time
		ChangeSignificance string
		Message            string
		Checksum           string
	}
	data := make([]Data, len(monitors))
	for _, monitor := range monitors {
		report, findErr := models.FindLastReportForMonitor(h.db, monitor)
		if findErr != nil {
			fmt.Println("No report for monitor", monitor, "ERROR", findErr)
		} else {
			data = append(data, Data{
				MonitorID:          monitor.ID(),
				ScriptPath:         monitor.ScriptPath(),
				LastRan:            monitor.LastRun(),
				ChangeSignificance: report.Change().String(),
				Message:            report.Message(),
				Checksum:           report.Checksum(),
			})
		}
	}
	// Serve the page with the data about monitors and their recent reports.
	t, err := template.ParseFiles(
		path.Join(h.cfg.TemplateDir, reportsPage),
		path.Join(h.cfg.TemplateDir, headTemplate),
		path.Join(h.cfg.TemplateDir, navTemplate))
	if err != nil {
		fmt.Println("Error parsing reports page template", err)
		InternalError(res, req)
		return
	}
	t.Execute(res, struct{ Reports []Data }{data})
}
