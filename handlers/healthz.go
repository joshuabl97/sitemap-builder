package handlers

import (
	"fmt"
	"net/http"
)

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	// Check the health of your application here.
	// You can perform various health checks and return an appropriate response.
	// For simplicity, this example always returns a 200 OK response.
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}
