package models

import (
	"database/sql"
	"errors"
	"time"
)

// Importance is used to rate the degree of change observed on a site.
type Importance uint

const (
	NoChange      Importance = 0 // The page hasn't changed.
	MinorUpdate   Importance = 1 // A minor textual change occurred, such as a typo fix.
	ContentChange Importance = 2 // The content has been modified in a meaningful way.
	Rewritten     Importance = 4 // The entirety of the site's content has been replaced.
	Deleted       Importance = 8 // The site has been completed deleted.
)

// Report contains information output by a monitor script informing us of any
// changes on the site being monitored. The stateData (state in JSON) field can be
// used by monitor scripts to include any extra data that might be useful to them.
type Report struct {
	id                 int                    `json:"-"`
	createdBy          int                    `json:"createdBy"`
	createdAt          time.Time              `json:"createdAt"`
	changeSignificance Importance             `json:"lastChangeSignificance"`
	messageToAdmin     string                 `json:"message"`
	checksum           string                 `json:"checksum"`
	stateData          map[string]interface{} `json:"state"`
}

// FindLastReportForMonitor looks up the last report output by a monitor script.
// Arguments:
// db: A database connection.
// monitor: The monitor that is going to be run.
// Returns:
// A Report containing the monitor script's last output, if one can be found,
// or else an error if something fails in the database.
func FindLastReportForMonitor(db *sql.DB, monitor Monitor) (Report, error) {
	r := Report{}
	err := db.QueryRow(QFindLastReportForMonitor, monitor.ID()).Scan(
		&r.id, &r.createdAt, &r.changeSignificance, &r.messageToAdmin, &r.checksum, &r.stateData)
	if err != nil {
		return Report{}, err
	}
	return r, nil
}

// Save creates a new Report in the database for an admin to view later and to be
// provided as input during the next invokation of the monitor script that
// produced it.
// Arguments:
// db: A database connection.
// Returns:
// An error if the database insertion fails.
func (r *Report) Save(db *sql.DB) error {
	_, err := db.Exec(QSaveReport,
		r.createdBy, r.createdAt, r.changeSignificance, r.messageToAdmin, r.checksum, r.stateData)
	if err != nil {
		return err
	}
	err = db.QueryRow(QLastRowID).Scan(&r.id)
	return err
}

// Update always returns an error because we don't want to allow reports to be changed.
// Arguments:
// db: A database connection.
// Returns:
// An error saying that reports cannot be updated.
func (r *Report) Update(db *sql.DB) error {
	return errors.New("cannot change a report")
}

// Delete always returns an error because we don't want to allow reports to be deleted.
// Arguments:
// db: A database connection.
// Returns:
// An error saying that reports cannot be deleted.
func (r *Report) Delete(db *sql.DB) error {
	return errors.New("cannot delete a report")
}
