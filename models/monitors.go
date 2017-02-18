package models

import (
	"database/sql"
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
// cmd: The interpreter used to run the script.
// filePath: The path to the script saved on disk.
// waitBetweenRuns: The amount of time (minutes) to wait between script runs.
// expectedRuntime: The amount of time (seconds) expected to run the script.
// Returns:
// A new Monitor containing the provided data.
func NewMonitor(
	creator Administrator,
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
		interpreter: string(cmd),
		scriptPath:  filePath,
		createdAt:   time.Now(),
		lastRan:     time.Now().Add(-1 * waitBetweenRuns),
		waitPeriod:  waitMinutes,
		timeToRun:   runTimeSeconds,
	}
}

// Save inserts a new monitor into the database and updates the id field.
// WARNING: Save should *not* be called more than once on a model.
func (m *Monitor) Save(db *sql.DB) error {
	return nil
}
