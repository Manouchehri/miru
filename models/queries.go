package models

// QLastRowID is an SQL query that gets the ID of the last row created.
const QLastRowID = `select last_insert_rowid();`

// QInitMonitorsTable is an SQL query that creates the monitors table.
const QInitMonitorsTable = `
create table if not exists monitors (
  id integer primary key,
  interpreter varchar(16) not null,
  script_location varchar(255) not null,
	created_for integer,
  created_by integer,
  created_at timestamp,
  last_ran_at timestamp,
  wait_period_minutes integer,
  expected_run_time integer,
	foreign key(created_for) references requests(id),
  foreign key(created_by) references archivers(id)
);`

// QInitArchiversTable is an SQL query that creates the archivers table.
const QInitArchiversTable = `
create table if not exists archivers (
  id integer primary key,
  made_admin_by integer,
  is_administrator bool default false,
  email_address varchar(64) unique not null,
  password_hash varchar(255) not null,
  last_login_ip varchar(45),
  last_login_time timestamp
);`

// QInitSessionsTable is an SQL query that creates the sessions table.
const QInitSessionsTable = `
create table if not exists sessions (
  id varchar(32) primary key,
  owner int unique,
  created_at timestamp not null,
  expires_at timestamp not null,
  ip_address varchar(45) not null,
  foreign key(owner) references archivers(id)
);`

// QInitRequestsTable is an SQL query that creates the requests table.
const QInitRequestsTable = `
create table if not exists requests (
	id integer primary key,
	created_by integer,
	created_at timestamp not null,
	url text not null,
	instructions text,
	foreign key(created_by) references archivers(id)
);`

// QSaveMonitor is an SQL query that saves a new monitor.
const QSaveMonitor = `
insert into monitors (
  interpreter, script_location, created_for, created_by, created_at,
  last_ran_at, wait_period_minutes, expected_run_time
) values ($1, $2, $3, $4, $5, $6, $7, $8);`

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
  id, interpreter, script_location, created_for, created_by, created_at,
  last_ran_at, wait_period_minutes, expected_run_time
from monitors
where
  ((select julianday('now')) - julianday(last_ran_at)) * (60 * 24)
  >= wait_period_minutes
limit $1;`

// QIsUserAnAdmin is an SQL query that checks if a given user has
// administrator privileges, allowing them to create monitors.
const QIsUserAnAdmin = `select is_administrator from archivers where id = $1;`

// QSaveArchiver is an SQL query that creates a new archiver account.
const QSaveArchiver = `
insert into archivers (
  made_admin_by, is_administrator, email_address,
  password_hash, last_login_ip, last_login_time
) values ($1, $2, $3, $4, $5, $6);`

// QUpdateArchiver is an SQL query that updates an existing archiver account.
const QUpdateArchiver = `
update archivers set
  made_admin_by = $1,
  is_administrator = $2,
  email_address = $3,
  password_hash = $4,
  last_login_ip = $5,
  last_login_time = $6
where id = $7;`

// QDeleteArchiver is an SQL query that deletes a user account entirely.
const QDeleteArchiver = `delete from archivers where id = $1;`

// QFindArchiver is an SQL query that looks for an archiver given their ID.
const QFindArchiver = `
select
  email_address, password_hash, made_admin_by,
  is_administrator, last_login_ip, last_login_time
from archivers
where id = $1;`

// QFindArchiverByEmail is an SQL query that attempts to find a user account
// associated with a given email address.
const QFindArchiverByEmail = `
select
  id, made_admin_by, is_administrator, password_hash,
  last_login_ip, last_login_time
from archivers
where email_address = $1;`

// QSaveSession is an SQL query that creates a new session for an
// authenticated archiver.
const QSaveSession = `
insert into sessions (
  id, owner, created_at, expires_at, ip_address
) values ($1, $2, $3, $4, $5);`

// QDeleteSession is an SQL query that deletes a session.
const QDeleteSession = `delete from sessions where id = $1;`

// QFindSession is an SQL query that finds a session given its token (id).
const QFindSession = `
select
  owner, created_at, expires_at, ip_address
from sessions
where id = $1;`

// QSaveRequest is an SQL query that inserts a new request.
const QSaveRequest = `
insert into requests (
	created_by, created_at, url, instructions
) values ($1, $2, $3, $4);`

// QRejectRequest is an SQL query that deletes a request.
const QRejectRequest = `delete from requests where id = $1;`

// QIsRequestFulfilled is an SQL query that determines whether a request has been
// fulfilled by looking for a monitor that was created to fulfill it.
const QIsRequestFulfilled = `
select exists(
	select id
	from monitors 
	where created_for = $1
);`

// QFindRequest is an SQL query that attempts to find a monitor request.
const QFindRequest = `
select created_by, created_at, url, instructions
from requests
where id = $1;`

// QListPendingRequests is an SQL query that finds all requests for which no
// monitor has yet been created to fulfill.
const QListPendingRequests = `
select R.id, R.created_by, R.created_at, R.url, R.instructions
from requests R
where not exists(
	select M.id
	from monitors M
	where M.created_for = R.id
);`
