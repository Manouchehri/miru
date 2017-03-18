package tasks

import (
	"../models"

	"database/sql"
	"fmt"
	"time"
)

// RunMonitors runs until signalled to terminate, periodically fetching new
// monitors whose scripts are ready to be run, running them, and then manages
// their reports.
func RunMonitors(db *sql.DB, sleepPeriod time.Duration, errors chan<- error, terminate <-chan bool) {
	results := make(chan models.Report)
	terminated := false
	for !terminated {
		select {
		case <-time.After(sleepPeriod):
			monitors, err := models.FindReadyMonitors(db, 1)
			if err != nil {
				errors <- err
			}
			for _, monitor := range monitors {
				fmt.Println("+++ Found monitor", monitor)
				monitor.SetLastRun()
				updateErr := monitor.Update(db)
				if updateErr != nil {
					errors <- updateErr
				}
				lastReport, findErr := models.FindLastReportForMonitor(db, monitor)
				if findErr != nil {
					fmt.Println("Couldn't find report for monitor", findErr)
					lastReport = models.NewReport(monitor)
					saveErr := lastReport.Save(db)
					if saveErr != nil {
						fmt.Println("could not save new report", saveErr)
						errors <- saveErr
					}
				}
				go RunMonitorScript(monitor, lastReport, results, errors)
			}
		case result := <-results:
			fmt.Println("Got result", result)
			saveErr := result.Save(db)
			if saveErr != nil {
				errors <- saveErr
			}
		case <-terminate:
			terminated = true
		}
	}
	close(errors)
	close(results)
}
