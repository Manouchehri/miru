package models

import (
  "database/sql"
)

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
  return nil
}
