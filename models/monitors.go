package models

import (
	"database/sql"
	"fmt"
	"math"
	"time"
)

// Interpreter is a pseudo-enum covering the script interpreters supported.
type Interpreter string

const (
	// PythonInterpreter identifies a python script
	PythonInterpreter Interpreter = "python"

	// RubyInterpreter identifies a ruby script
	RubyInterpreter Interpreter = "ruby"

	// PerlInterpreter identifies a perl script
	PerlInterpreter Interpreter = "perl"
)

// Monitor is the model for "rules" that specify a script to run in order to
// check a website for changes.
type Monitor struct {
	id          int
	createdFor  int
	createdBy   int
	interpreter string
	scriptPath  string
	createdAt   time.Time
	lastRan     time.Time
	waitPeriod  uint
	timeToRun   uint
}

// NewMonitor is the constructor for the monitor type. When a new script is
// uploaded to monitor a website, a monitor should be created using NewMonitor.
// Arguments:
// creator: The administrator that uploaded the update checking script.
// requestedBy: The request that is being fulfilled.
// cmd: The interpreter used to run the script.
// filePath: The path to the script saved on disk.
// waitBetweenRuns: The amount of time (minutes) to wait between script runs.
// expectedRuntime: The amount of time (seconds) expected to run the script.
// Returns:
// A new Monitor containing the provided data.
func NewMonitor(
	creator Archiver,
	requestedBy Request,
	cmd Interpreter,
	filePath string,
	waitBetweenRuns time.Duration,
	expectedRuntime time.Duration,
) Monitor {
	waitMinutes := uint(math.Ceil(waitBetweenRuns.Minutes()))
	runTimeSeconds := uint(math.Ceil(expectedRuntime.Seconds()))
	return Monitor{
		id:          -1,
		createdBy:   creator.ID(),
		createdFor:  requestedBy.ID(),
		interpreter: string(cmd),
		scriptPath:  filePath,
		createdAt:   time.Now(),
		lastRan:     time.Now().Add(-1 * waitBetweenRuns),
		waitPeriod:  waitMinutes,
		timeToRun:   runTimeSeconds,
	}
}

// ListMonitors attempts to get a list of all of the monitors registered.
// Arguments:
// db: A database connection.
// Returns:
// A list of all monitors in the dabatase, or as many as can be read until an
// error occurs. The first error encountered trying to read data is returned.
func ListMonitors(db *sql.DB) ([]Monitor, error) {
	allMonitors := []Monitor{}
	rows, err := db.Query(QListMonitors)
	if err != nil {
		return allMonitors, err
	}
	for rows.Next() {
		m := Monitor{}
		err = rows.Scan(
			&m.id, &m.interpreter, &m.scriptPath, &m.createdFor, &m.createdBy,
			&m.createdAt, &m.lastRan, &m.waitPeriod, &m.timeToRun)
		if err != nil {
			break
		}
		allMonitors = append(allMonitors, m)
	}
	return allMonitors, err
}

// FindReadyMonitors finds monitors that we've waited long enough to run again.
// The function will return the first error it encounters, along with any
// monitors retrieved until that point.
// Arguments:
// db: A database connection.
// limit: The maximum number of monitors to fetch.
// Returns:
// An array of monitors that can be run, and the first error if one occurs.
func FindReadyMonitors(db *sql.DB, limit uint) ([]Monitor, error) {
	allMonitors := make([]Monitor, limit)
	monitorsFound := 0
	rows, err := db.Query(QFindReadyMonitors, limit)
	if err != nil {
		return []Monitor{}, err
	}
	for rows.Next() {
		var m Monitor
		err = rows.Scan(
			&m.id, &m.interpreter, &m.scriptPath, &m.createdFor, &m.createdBy,
			&m.createdAt, &m.lastRan, &m.waitPeriod, &m.timeToRun)
		if err != nil {
			break
		}
		allMonitors[monitorsFound] = m
		monitorsFound++
	}
	shrunkArray := []Monitor{}
	shrunkArray = append(shrunkArray, allMonitors[:monitorsFound]...)
	return shrunkArray, err
}

// Interpreter is a getter function that converts the name of the monitor's
// script type back into an Interpreter type.
func (m Monitor) Interpreter() Interpreter {
	return Interpreter(m.interpreter)
}

// ScriptPath is a getter function for the monitor's script path on disk.
func (m Monitor) ScriptPath() string {
	return m.scriptPath
}

// ID is a getter function for a monitor's unique identifier.
// Returns:
// The monitor's id.
func (m Monitor) ID() int {
	return m.id
}

// SetLastRun sets the monitor's last run time to now.
func (m *Monitor) SetLastRun() {
	m.lastRan = time.Now()
}

// Save inserts a new monitor into the database and updates the id field.
// WARNING: Save should *not* be called more than once on a model.
func (m *Monitor) Save(db *sql.DB) error {
	fmt.Println("Saving monitor for request ID", m.createdFor)
	_, err := db.Exec(QSaveMonitor,
		m.interpreter, m.scriptPath, m.createdFor, m.createdBy, m.createdAt,
		m.lastRan, m.waitPeriod, m.timeToRun)
	if err != nil {
		return err
	}
	err = db.QueryRow(QLastRowID).Scan(&m.id)
	return err
}

// Update modifies the monitor's database row to set the time the monitor was
// last run, the time we want to wait between running it, and the amount of
// time to allow the monitor to run for.
func (m *Monitor) Update(db *sql.DB) error {
	_, err := db.Exec(QUpdateMonitor,
		m.lastRan, m.waitPeriod, m.timeToRun, m.id)
	return err
}

// Delete removes the monitor from the database.
func (m *Monitor) Delete(db *sql.DB) error {
	_, err := db.Exec(QDeleteMonitor, m.id)
	return err
}
