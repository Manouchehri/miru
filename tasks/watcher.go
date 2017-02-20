package tasks

import (
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
	for {
		select {
		case <-time.After(sleepPeriod):
			fmt.Println("Look for and run a task")
		case <-terminate:
			break
		}
	}
	close(errors)
}
