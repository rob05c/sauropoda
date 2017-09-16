package hndl

import (
	"fmt"
	"net/http"
	"time"
)

// handleNow handles a request for the current server time
func handleNow(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Fprintf(w, "%s", time.Now()) // TODO change to Write and time.Format, for performance and obviousness
}
