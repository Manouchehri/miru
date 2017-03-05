package handlers

import "errors"

const headTemplate string = "head.html"
const navTemplate string = "nav.html"
const errorTemplate string = "error.html"

var (
	errTemplateLoad       = errors.New("failed to load a page template")
	errInvalidCredentials = errors.New("the provided credentials are invalid")
	errDatabaseOperation  = errors.New("an internal database error occurred")
	errNotAllowed         = errors.New("you are not allowed to do that")
	errGenericInvalidData = errors.New("some of the input provided is invalid")
	errCreateFile         = errors.New("could not create a file for the monitor script")
	errBadPassword        = errors.New("password and repeated password must match and contain " +
		"at least one lowercase and uppercase letter, symbol, and number")
)
