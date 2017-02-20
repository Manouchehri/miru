package handlers

import (
	"encoding/hex"
	"errors"
	"io"
	"os"

	"../config"

	"database/sql"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"path"
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
	cfg *config.Config
}

// NewUploadScriptHandler is the constructor function for UploadScriptHandler.
// Arguments:
// c: A reference to the application's global configuration.
// db: A reference to a database connection.
func NewUploadScriptHandler(c *config.Config, db *sql.DB) UploadScriptHandler {
	return UploadScriptHandler{
		cfg: c,
		db:  db,
	}
}

// NewUploadPageHandler is the constructor function for UploadPageHandler.
// Arguments:
// cfg: A reference to the application's global configuration.
func NewUploadPageHandler(cfg *config.Config) UploadPageHandler {
	return UploadPageHandler{
		cfg: cfg,
	}
}

// ServeHTTP handles file uploads containing new monitor scripts.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h UploadScriptHandler) ServeHTTP(
	res http.ResponseWriter, req *http.Request) {
	filetype := req.FormValue("filetype")
	ext, ftErr := filetypeExtension(filetype)
	if ftErr != nil {
		BadRequest(res, req)
		return
	}
	file, _, openErr := req.FormFile("script")
	if openErr != nil {
		fmt.Printf("Error: %v\n", openErr)
		BadRequest(res, req)
		return
	}
	defer file.Close()
	filename := generateUniqueFilename(h.cfg.ScriptDir, ext)
	toDisk, openErr := os.Create(filename)
	if openErr != nil {
		fmt.Printf("Error: %v\n", openErr)
		InternalError(res, req)
		return
	}
	defer toDisk.Close()
	io.Copy(toDisk, file)
	res.Write([]byte("hello world"))
}

// ServeHTTP serves a page that administrators can use to upload new
// monitor scripts through.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h UploadPageHandler) ServeHTTP(
	res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles(path.Join(h.cfg.TemplateDir, uploadPage))
	if err != nil {
		InternalError(res, req)
		return
	}
	t.Execute(res, nil)
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
