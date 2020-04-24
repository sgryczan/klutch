package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sgryczan/klutch/common"
)

var version = common.Version

type aboutResponse struct {
	Version string
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /about About About
	//
	// Returns information about the application
	// ---
	// consumes:
	// - text/plain
	// produces:
	// - text/plain
	//
	// responses:
	//   '200':
	//     description: About
	//     type: string
	w.WriteHeader(http.StatusOK)
	data := &aboutResponse{
		Version: version,
	}
	res, _ := json.Marshal(data)

	fmt.Fprintf(w, "%s", res)
}
