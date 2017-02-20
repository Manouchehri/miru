package models

// QLastRowID is an SQL query that gets the ID of the last row created.
const QLastRowID = `select last_insert_rowid();`

// QInitMonitorsTable is an SQL query that creates the monitors table.
const QInitMonitorsTable = `
create table if not exists monitors (
  id integer primary key,
  interpreter varchar(16) not null,
  script_location varchar(255) not null,
  created_by integer,
  created_at timestamp,
  last_ran_at timestamp,
  wait_period_minutes integer,
  expected_run_time integer,
  foreign key(created_by) references archivers(id)
);`

// QInitArchiversTable is an SQL query that creates the archivers table.
const QInitArchiversTable = `
create table if not exists archivers (
  id integer primary key,
  is_administrator bool default false,
  email_address varchar(64) unique not null,
  password_hash varchar(255) not null,
  last_login_ip varchar(45),
  last_login_time timestamp
);`

// QSaveMonitor is an SQL query that saves a new monitor.
const QSaveMonitor = `
insert into monitors (
  interpreter, script_location, created_by, created_at,
  last_ran_at, wait_period_minutes, expected_run_time
) values ($1, $2, $3, $4, $5, $6, $7);`

// QUpdateMonitor is an SQL query that updates an existing monitor.
// Note that we would rather a monitor be deleted and a new one created if a new
// script is uploaded.
const QUpdateMonitor = `
update monitors set
  last_ran_at = $1,
  wait_period_minutes = $2,
  expected_run_time = $3
where id = $4;`

// QDeleteMonitor is an SQL query that deletes a monitor.
const QDeleteMonitor = `delete from monitors where id = $1;`

// QFindReadyMonitors is an SQL query that retrieves monitors whose wait period
// between runs has expired.
const QFindReadyMonitors = `
select
  id, interpreter, script_location, created_by, created_at,
  last_ran_at, wait_period_minutes, expected_run_time
from monitors
where
  ((select julianday('now')) - julianday(last_ran_at)) * (60 * 24)
  >= wait_period_minutes
limit $1;`

// QIsUserAnAdmin is an SQL query that checks if a given user has
// administrator privileges, allowing them to create monitors.
const QIsUserAnAdmin = `select is_administrator from archivers where id = $1;`
