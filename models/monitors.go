package models

import (
  "time"
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
