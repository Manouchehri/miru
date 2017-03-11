package handlers

import (
	"../auth"
	"../config"
	"../models"

	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

// uploadPage is the name of the HTML template containing a monitor script
// uploading form.
const uploadPage string = "scriptupload.html"

// filenameLength is the number of random bytes to get to generate names
// for uploaded scripts with.
const filenameLength int = 16

// UploadScriptHandler implements net/http.ServeHTTP to handle new monitor
// script uploads from administrators.
type UploadScriptHandler struct {
	cfg *config.Config
	db  *sql.DB
}

// UploadPageHandler implements net/http.ServeHTTP to sreve a page to
// administrators that they can use to upload new monitor scripts.
type UploadPageHandler struct {
	cfg       *config.Config
	db        *sql.DB
	Successes []string
}

// NewUploadScriptHandler is the constructor function for UploadScriptHandler.
// Arguments:
// c: A reference to the application's global configuration.
// db: A reference to a database connection.
// Returns:
// A new UploadScriptHandler that can be bound to a router.
func NewUploadScriptHandler(c *config.Config, db *sql.DB) UploadScriptHandler {
	return UploadScriptHandler{
		cfg: c,
		db:  db,
	}
}

// NewUploadPageHandler is the constructor function for UploadPageHandler.
// Arguments:
// cfg: A reference to the application's global configuration.
// Returns:
// A new UploadPageHandler that can be bound to a router.
func NewUploadPageHandler(cfg *config.Config, db *sql.DB) UploadPageHandler {
	return UploadPageHandler{
		cfg:       cfg,
		db:        db,
		Successes: []string{},
	}
}

// PushSuccessMsg adds a new message that will be displayed on the page served by the
// handler to indicate a successful operation.
// Arguments:
// msg: A success message to display to the user.
func (h *UploadPageHandler) PushSuccessMsg(msg string) {
	h.Successes = append(h.Successes, msg)
}

// ServeHTTP handles file uploads containing new monitor scripts.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h UploadScriptHandler) ServeHTTP(
	res http.ResponseWriter, req *http.Request) {
	// Check that the request is being made by an authenticated administrator.
	fmt.Println(req.Cookies())
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		fmt.Println("Could not find cookie", err)
		BadRequest(res, req, h.cfg, errNotAllowed, false, false)
		return
	}
	activeUser, err := models.FindSessionOwner(h.db, cookie.Value)
	if err != nil || !activeUser.IsAdmin() {
		fmt.Println("Could not get cookie owner", err)
		BadRequest(res, req, h.cfg, errNotAllowed, err == nil, false)
		return
	}
	// Extract inputs from the form.
	waitPeriod, parseErr1 := strconv.Atoi(req.FormValue("waitPeriod"))
	expectedRuntime, parseErr2 := strconv.Atoi(req.FormValue("expectedRuntime"))
	requestID, parseErr3 := strconv.Atoi(req.FormValue("satisfiedRequest"))
	filetype := req.FormValue("filetype")
	ext, ftErr := filetypeExtension(filetype)
	if ftErr != nil || parseErr1 != nil || parseErr2 != nil || parseErr3 != nil {
		fmt.Println(ftErr)
		BadRequest(res, req, h.cfg, errGenericInvalidData, true, true)
		return
	}
	// Find the request that is being fulfilled to establish relational data.
	request, findErr := models.FindRequest(h.db, requestID)
	if findErr != nil {
		fmt.Println("no such request", requestID, "ERROR", findErr)
		BadRequest(res, req, h.cfg, errGenericInvalidData, true, true)
		return
	}
	file, _, openErr := req.FormFile("script")
	if openErr != nil {
		fmt.Printf("Error: %v\n", openErr)
		BadRequest(res, req, h.cfg, errCreateFile, true, true)
		return
	}
	defer file.Close()
	// Find a place to save the file to on disk.
	filename := generateUniqueFilename(h.cfg.ScriptDir, ext)
	toDisk, openErr := os.Create(filename)
	if openErr != nil {
		fmt.Printf("Error: %v\n", openErr)
		InternalError(res, req, h.cfg, errCreateFile, true, true)
		return
	}
	defer toDisk.Close()
	io.Copy(toDisk, file)
	// Create a new Monitor in the database.
	monitor := models.NewMonitor(
		activeUser,
		request,
		models.Interpreter(filetype),
		filename,
		time.Duration(waitPeriod)*time.Minute,
		time.Duration(expectedRuntime)*time.Second)
	fmt.Println("Creating monitor for", request.ID())
	fmt.Println("Monitor", monitor)
	saveErr := monitor.Save(h.db)
	if saveErr != nil {
		fmt.Println(saveErr)
		InternalError(res, req, h.cfg, errDatabaseOperation, true, true)
		return
	}
	handler := NewUploadPageHandler(h.cfg, h.db)
	handler.PushSuccessMsg(fmt.Sprintf("Successfully created a new monitor script with ID %d", monitor.ID()))
	handler.ServeHTTP(res, req)
}

// ServeHTTP serves a page that administrators can use to upload new
// monitor scripts through.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h UploadPageHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Check that the request is coming from an authenticated archiver.
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		BadRequest(res, req, h.cfg, errNotAllowed, false, false)
		return
	}
	activeUser, err := models.FindSessionOwner(h.db, cookie.Value)
	if err != nil {
		BadRequest(res, req, h.cfg, errNotAllowed, false, false)
		return
	}
	requestIDs, found := req.URL.Query()["id"]
	if !found || len(requestIDs) == 0 {
		fmt.Println("Need a request id")
		BadRequest(res, req, h.cfg, errors.New("missing request id url parameter"), true, activeUser.IsAdmin())
		return
	}
	requestID, parseErr := strconv.Atoi(requestIDs[0])
	if parseErr != nil {
		fmt.Println("Need a request valid id")
		BadRequest(res, req, h.cfg, errGenericInvalidData, true, activeUser.IsAdmin())
		return
	}
	t, err := template.ParseFiles(
		path.Join(h.cfg.TemplateDir, uploadPage),
		path.Join(h.cfg.TemplateDir, headTemplate),
		path.Join(h.cfg.TemplateDir, navTemplate))
	if err != nil {
		InternalError(res, req, h.cfg, errTemplateLoad, true, true)
		return
	}
	t.Execute(res, struct {
		CreatedFor  int
		LoggedIn    bool
		UserIsAdmin bool
		Successes   []string
	}{requestID, true, activeUser.IsAdmin(), h.Successes})
}

// generateUniqueFilename produces a filename that is guaranteed to be unique.
// It continuously generates 16-byte script names, encoded as hex, until one
// is created that isn't already taken.
// Arguments:
// scriptDir: The directory that monitor scripts are saved to.
// ext: The filetype extension to append to the end of the filename.
// Returns:
// A filename that is guaranteed to be unique.
func generateUniqueFilename(scriptDir string, ext string) string {
	// We don't need cryptographically random names- pseudorandom will do.
	rand.Seed(int64(time.Now().Unix()))
	bytes := make([]byte, filenameLength)
	for {
		bytesRead, readErr := rand.Read(bytes)
		for readErr != nil || bytesRead != filenameLength {
			bytesRead, readErr = rand.Read(bytes)
		}
		filename := path.Join(scriptDir, hex.EncodeToString(bytes))
		f, openErr := os.Open(filename)
		if openErr != nil {
			return filename + "." + ext
		}
		f.Close()
	}
}

// Converts a filetype, the values in the upload form's filetype dropdown
// menu, into the file type's extension.
// Arguments:
// The filetype specified in the script upload form.
// Returns:
// The file type's extension if the file type is supported, and an error in
// the case that the file type is not supported.
func filetypeExtension(filetype string) (string, error) {
	switch filetype {
	case "python":
		return "py", nil
	case "ruby":
		return "rb", nil
	case "perl":
		return "pl", nil
	default:
		return "", errors.New("unknown filetype")
	}
}
