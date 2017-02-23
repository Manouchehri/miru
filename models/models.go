package models

import (
	"database/sql"
)

// Model is an interface for all model types. We expect at least the following
// basic database manipulations to be implemented for each model.
type Model interface {
	Save(*sql.DB) error
	Update(*sql.DB) error
	Delete(*sql.DB) error
}

// InitializeTables creates all of the database tables used by miru.
// Arguments:
// db: A connection to the database.
// Returns:
// The first error that is encountered creating a table.
func InitializeTables(db *sql.DB) error {
	_, err := db.Exec(QInitArchiversTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(QInitMonitorsTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(QInitSessionsTable)
	if err != nil {
		return err
	}
	return nil
}
