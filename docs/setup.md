# Development Setup

This document explains how to build and run Miru, assuming you have already [forked and cloned the repository](https://github.com/zsck/miru/blob/master/CONTRIBUTING.md).

## Tools and Dependencies

### The Go Compiler and Tools

The first thing you'll need is the [Go compiler](https://golang.org/dl/), which includes the [go fmt](https://golang.org/cmd/gofmt/) and [go vet](https://golang.org/cmd/vet/) tools.  Once installed, you should be able to run the following command in your terminal and see a version number greater than or equal to `go1.8`:

```
go version
```

### SQLite

The next thing you'll need is a copy of [SQLite](https://sqlite.org/), the database management system that Miru uses. It is best to install SQLite using your operating system's package manager if you use Linux, or the [Homebrew](https://brew.sh/) package manager if you use macOS.  To confirm that your installation worked, you should be able to start sqlite in your terminal.

```
sqlite3
.q
```

Entering `.q` into the sqlite prompt should cause the program to terminate and bring you back to your terminal.

### Go Library Dependencies

Miru relies on a few third-party libraries which must be installed before it can be built. If you already have the Go compiler and related tools installed, you can install all of these dependencies by running the following commands in your terminal.

```
go get github.com/gorilla/mux
go get github.com/mattn/go-sqlite3
go get github.com/StratumSecurity/scryptauth
```

## Building Miru

Once you've installed the Go compiler and all of Miru's dependencies, you should be able to build Miru in one step by running the following command in your terminal, inside of the `miru/` directory created when you cloned your fork of the repository.

```
go build
```

If this command completes successfully with no errors, you should see a new executable file in the `miru/` directory, called `miru`.

## Configuration

If you are running Miru locally for development purposes, you should not have to change any of these options before proceeding to the next section of this guide.

The configuration information used to run Miru can be found in the `miru/config/config.json` file, and has the following data in it.

```json
{
  "bindAddress": "127.0.0.1:3000",
  "templateDir": "templates",
  "database": "miru.db",
  "scriptDir": "monitorscripts"
}
```

* `"bindAddress"`  is an IP address and port number for the Miru web server to bind to in order to listen for requests from users.
* `"templateDir"` is the path to the directory containing Miru's HTML template files.
* `"database"` is the name of the database file to store Miru's SQLite data in and will be created by Miru the first time it's run.
* `"scriptDir"` is the path to the directory that you would like to have Miru save uploaded monitoring scripts to. Note that this directory **must exist before Miru is run**.

## Running Miru

Once compiled, starting Miru is as simple as executing the binary produced by the compiler by running the following command from your terminal in the `miru/` directory.

```
./miru
```

You should immediately see a message informing you that the Miru web server is running and listening for requests, like so:

> Listening on 127.0.0.1:3000

Now, if you enter `127.0.0.1:3000` (or whatever value you supplied to the configuration file) into your browser's address bar, you should see the Miru index page.

### Creating the first administrator

At the time of this writing, Miru allows administrator users to promote other registered users to administrators, however the first administrator account must be created manually. First, register a user account for yourself by clicking `Register` at the top right of the index page, and fill out your account information.  Once you've successfully registered, go to your terminal and issue the following commands.

NOTE that you should only run these commands **before making your Miru instance publicly available on the Internet!**

```
sqlite3 miru.db
update archivers set is_administrator = 1;
.q
```

Assuming you did not change the `"database"` field in your configuration file, this will connect you to Miru's database and then, once in the SQLite prompt, update **all** archivers (users) to make them administrators. 

## Conclusion

And that's it!  You're now ready to begin [using Miru](https://github.com/zsck/miru/blob/master/docs/using-miru.md)!