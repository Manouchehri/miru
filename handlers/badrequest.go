package handlers

import (
	"net/http"
)

// BadRequest is a simple net/http HandlerFunc that will write an error
// message to users if something is wrong with a request.
// Arguments:
// res: Provided by the net/http server, used to write the response.
// req: Provided by the net/http server, contains information about the request.
func BadRequest(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte("Bad request"))
}
