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
// Arguments:
// db: A database connection.
// sleepPeriod: The duration of time to wait between reading the DB for tasks.
// errors: A channel through which any errors can be written.
// terminate: A channel through which a termination signal can be read.
func RunMonitors(
	db *sql.DB, sleepPeriod time.Duration,
	errors chan<- error, terminate <-chan bool) {
	results := make(chan Result)
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
				go RunMonitorScript(monitor, results, errors)
			}
		case result := <-results:
			fmt.Println("Got result", result)
		case <-terminate:
			terminated = true
		}
	}
	close(errors)
	close(results)
}
