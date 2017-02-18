package handlers

import (
	"net/http"
)

// InternalError is a simple net/http HandlerFunc that will write an error
// message to users if something goes wrong with the application.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func InternalError(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusInternalServerError)
	res.Write([]byte("Internal server error"))
}
