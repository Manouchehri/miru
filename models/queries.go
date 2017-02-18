package models

/// QInitAdministratorsTable is an SQL query that creates the administrators
// table.
const QInitAdministratorsTable = `
create table if not exists administrators (
  id integer primary key,
  email_address varchar(64) unique not null,
  password_hash varchar(255) not null,
  last_login_ip varchar(45),
  last_login_time timestamp
);`

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
  email_address varchar(64) unique not null,
  password_hash varchar(255) not null,
  last_login_ip varchar(45),
  last_login_time timestamp
);`
