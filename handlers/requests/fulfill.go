package requests

import (
	"../../auth"
	"../../config"
	"../../models"
	"../common"
	"../fail"

	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

// filenameLength is the number of random bytes to get to generate names
// for uploaded scripts with.
const filenameLength int = 16

// FulfillHandler implements net/http.ServeHTTP to handle new monitor
// script uploads from administrators.
type FulfillHandler struct {
	cfg *config.Config
	db  *sql.DB
}

// NewFulfillHandler is the constructor function for FulfillHandler.
// Arguments:
// c: A reference to the application's global configuration.
// db: A reference to a database connection.
// Returns:
// A new FulfillHandler that can be bound to a router.
func NewFulfillHandler(c *config.Config, db *sql.DB) FulfillHandler {
	return FulfillHandler{
		cfg: c,
		db:  db,
	}
}

// ServeHTTP handles file uploads containing new monitor scripts.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h FulfillHandler) ServeHTTP(
	res http.ResponseWriter, req *http.Request) {
	// Check that the request is being made by an authenticated administrator.
	fmt.Println(req.Cookies())
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		fmt.Println("Could not find cookie", err)
		fail.BadRequest(res, req, h.cfg, common.ErrNotAllowed, false, false)
		return
	}
	activeUser, err := models.FindSessionOwner(h.db, cookie.Value)
	if err != nil || !activeUser.IsAdmin() {
		fmt.Println("Could not get cookie owner", err)
		fail.BadRequest(res, req, h.cfg, common.ErrNotAllowed, err == nil, false)
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
		fail.BadRequest(res, req, h.cfg, common.ErrGenericInvalidData, true, true)
		return
	}
	// Find the request that is being fulfilled to establish relational data.
	request, findErr := models.FindRequest(h.db, requestID)
	if findErr != nil {
		fmt.Println("no such request", requestID, "ERROR", findErr)
		fail.BadRequest(res, req, h.cfg, common.ErrGenericInvalidData, true, true)
		return
	}
	file, _, openErr := req.FormFile("script")
	if openErr != nil {
		fmt.Printf("Error: %v\n", openErr)
		fail.BadRequest(res, req, h.cfg, common.ErrCreateFile, true, true)
		return
	}
	defer file.Close()
	// Find a place to save the file to on disk.
	filename := generateUniqueFilename(h.cfg.ScriptDir, ext)
	toDisk, openErr := os.Create(filename)
	if openErr != nil {
		fmt.Printf("Error: %v\n", openErr)
		fail.InternalError(res, req, h.cfg, common.ErrCreateFile, true, true)
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
		fail.InternalError(res, req, h.cfg, common.ErrDatabaseOperation, true, true)
		return
	}
	handler := NewFulfillPageHandler(h.cfg, h.db)
	handler.PushSuccessMsg(fmt.Sprintf("Successfully created a new monitor script with ID %d", monitor.ID()))
	handler.ServeHTTP(res, req)
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
