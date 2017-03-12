package common

import "errors"

// HeadTemplate is the name of the template file that contains the HTML
// head contents for all pages.
const HeadTemplate string = "head.html"

// NavTemplate is the name of the template file that contains the HTML
// navigation contents for all pages.
const NavTemplate string = "nav.html"

// Common errors containing messages that are safe to show the user.
var (
	ErrTemplateLoad       = errors.New("failed to load a page template")
	ErrInvalidCredentials = errors.New("the provided credentials are invalid")
	ErrDatabaseOperation  = errors.New("an internal database error occurred")
	ErrNotAllowed         = errors.New("you are not allowed to do that")
	ErrGenericInvalidData = errors.New("some of the input provided is invalid")
	ErrCreateFile         = errors.New("could not create a file for the monitor script")
	ErrBadPassword        = errors.New("password and repeated password must match and contain " +
		"at least one lowercase and uppercase letter, symbol, and number")
	ErrInvalidEmail = errors.New("invalid email address")
)
