package archivers

import (
	"../../auth"
	"../../config"
	"../../models"
	"../common"
	"../fail"

	"database/sql"
	"fmt"
	"net/http"
	"time"
)

// LogoutHandler implements net/http.ServeHTTP to handle user logouts.
type LogoutHandler struct {
	cfg *config.Config
	db  *sql.DB
}

// NewLogoutHandler is the constructor function for a new LogoutHandler.
// Arguments:
// db: A database connection.
// Returns:
// A new LoginHandler that can be bound to a router.
func NewLogoutHandler(cfg *config.Config, db *sql.DB) LogoutHandler {
	return LogoutHandler{
		cfg: cfg,
		db:  db,
	}
}

// ServeHTTP handles requests to log a user out.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func (h LogoutHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(auth.SessionCookieName)
	if err != nil {
		fail.BadRequest(res, req, h.cfg, common.ErrNotAllowed, false, false)
		return
	}
	session, err := models.FindSession(h.db, cookie.Value)
	if err != nil {
		fail.BadRequest(res, req, h.cfg, common.ErrNotAllowed, false, false)
		return
	}
	err = session.Delete(h.db)
	if err != nil {
		fmt.Println("Failed to delete session", err)
		fail.InternalError(res, req, h.cfg, common.ErrDatabaseOperation, false, false)
		return
	}
	http.SetCookie(res, &http.Cookie{
		Name:    auth.SessionCookieName,
		Value:   "deleted",
		Path:    "/",
		Expires: time.Unix(0, 0),
	})
	http.Redirect(res, req, "/", http.StatusSeeOther)
}
