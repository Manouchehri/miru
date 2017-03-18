package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Importance is used to rate the degree of change observed on a site.
type Importance uint

// Levels of website content change significances.
const (
	NoChange      Importance = 0 // The page hasn't changed.
	MinorUpdate   Importance = 1 // A minor textual change occurred, such as a typo fix.
	ContentChange Importance = 2 // The content has been modified in a meaningful way.
	Rewritten     Importance = 3 // The entirety of the site's content has been replaced.
	Deleted       Importance = 4 // The site has been completed deleted.
)

// String produces a human-readable representation for each level of Report Importance.
func (i Importance) String() string {
	switch i {
	case NoChange:
		return "No Change"
	case MinorUpdate:
		return "Updated"
	case ContentChange:
		return "Changed"
	case Rewritten:
		return "Rewritten"
	case Deleted:
		return "Deleted"
	default:
		return "Uknown"
	}
}

// Report contains information output by a monitor script informing us of any
// changes on the site being monitored. The stateData (state in JSON) field can be
// used by monitor scripts to include any extra data that might be useful to them.
type Report struct {
	id                 int
	createdBy          int
	createdAt          time.Time
	changeSignificance Importance
	messageToAdmin     string
	checksum           string
	stateData          map[string]interface{}
}

// encodableReport is a private struct that contains public elements, which allows
// us to JSON encode the data that we need to write as input to monitor scripts.
type encodableReport struct {
	Change   Importance             `json:"lastChangeSignificance"`
	Message  string                 `json:"message"`
	Checksum string                 `json:"checksum"`
	State    map[string]interface{} `json:"state"`
}

// NewReport is the constructor function used to produce a monitor's first report.
func NewReport(monitor Monitor) Report {
	return Report{
		id:                 -1,
		createdBy:          monitor.ID(),
		createdAt:          time.Now(),
		changeSignificance: NoChange,
		messageToAdmin:     "first run",
		checksum:           "",
		stateData:          map[string]interface{}{},
	}
}

// FindLastReportForMonitor looks up the last Report output by a monitor script.
func FindLastReportForMonitor(db *sql.DB, monitor Monitor) (Report, error) {
	r := Report{}
	stateData := ""
	err := db.QueryRow(QFindLastReportForMonitor, monitor.ID()).Scan(
		&r.id, &r.createdAt, &r.changeSignificance, &r.messageToAdmin, &r.checksum, &stateData)
	if err != nil {
		return Report{}, err
	}
	fmt.Println("----- READ STATE DATA", stateData)
	decodeErr := json.Unmarshal([]byte(stateData), &r.stateData)
	if decodeErr != nil {
		return Report{}, decodeErr
	}
	r.createdBy = monitor.ID()
	return r, nil
}

// String converts the report to a JSON string.
func (r Report) String() string {
	encodable := encodableReport{
		Change:   r.changeSignificance,
		Message:  r.messageToAdmin,
		Checksum: r.checksum,
		State:    r.stateData,
	}
	encoded, encodeErr := json.Marshal(encodable)
	if encodeErr != nil {
		fmt.Println("couldn't encode report to json", encodeErr)
		return "{}"
	}
	return string(encoded)
}

// ID is a getter function for the Report's unique identifier.
func (r Report) ID() int {
	return r.id
}

// Change is a getter for the Report's states significance of the change since its
// last inspection.
func (r Report) Change() Importance {
	return r.changeSignificance
}

// Message is a getter for the Report's message to the admins about the state of
// the site that it is monitoring.
func (r Report) Message() string {
	return r.messageToAdmin
}

// Checksum is a getter for the Report's checksum of the site's data.
func (r Report) Checksum() string {
	return r.checksum
}

// SetChange is a setter function that sets the recorded significance of a site's last change.
func (r *Report) SetChange(change Importance) {
	r.changeSignificance = change
}

// SetMessage is a setter function that sets the message to an administrator in a report.
func (r *Report) SetMessage(message string) {
	r.messageToAdmin = message
}

// SetChecksum is a setter function that sets the checksum of the data inspected.
func (r *Report) SetChecksum(checksum string) {
	r.checksum = checksum
}

// SetState is a setter function for the report's state, which is data that it wants to communicate
// back to itself in successive runs.
func (r *Report) SetState(state map[string]interface{}) {
	r.stateData = state
}

// Save creates a new Report in the database for an admin to view later and to be
// provided as input during the next invokation of the monitor script that
// produced it.
func (r *Report) Save(db *sql.DB) error {
	stateData, encodeErr := json.Marshal(r.stateData)
	if encodeErr != nil {
		return encodeErr
	}
	_, err := db.Exec(QSaveReport,
		r.createdBy, r.createdAt, r.changeSignificance, r.messageToAdmin, r.checksum, string(stateData))
	if err != nil {
		return err
	}
	err = db.QueryRow(QLastRowID).Scan(&r.id)
	return err
}

// Update always returns an error because we don't want to allow Reports to be changed.
func (r *Report) Update(db *sql.DB) error {
	return errors.New("cannot change a report")
}

// Delete always returns an error because we don't want to allow Reports to be deleted.
func (r *Report) Delete(db *sql.DB) error {
	return errors.New("cannot delete a report")
}
