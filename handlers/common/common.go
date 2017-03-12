package common

import "errors"

// HeadTemplate is the name of the template file that contains the HTML
// head contents for all pages.
const HeadTemplate string = "head.html"

// NavTemplate is the name of the template file that contains the HTML
// navigation contents for all pages.
const NavTemplate string = "nav.html"

const passwordRules = `
the submitted password and repeated password must match and contain
at least one of each of the following: a lowercase letter, an uppercase
letter, a symbol, and a number.`

// Common errors containing messages that are safe to show the user.
var (
	ErrTemplateLoad          = errors.New("failed to load a page template")
	ErrInvalidCredentials    = errors.New("the provided credentials are invalid")
	ErrDatabaseOperation     = errors.New("an internal database error occurred")
	ErrNotAllowed            = errors.New("you are not allowed to do that")
	ErrGenericInvalidData    = errors.New("some of the input provided is invalid")
	ErrCreateFile            = errors.New("could not create a file for the monitor script")
	ErrBadPassword           = errors.New(passwordRules)
	ErrInvalidEmail          = errors.New("invalid email address")
	ErrLoginAttemptsExceeded = errors.New("you have exceeded the maximum number of login attempts")
)
