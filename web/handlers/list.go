package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sgryczan/klutch/common"
)

func ListHandler(db *common.RedisDatastore, w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /list List Item
	//
	// Lists all keys in the database
	// ---
	// consumes:
	// - text/plain
	// produces:
	// - text/plain
	//
	// responses:
	//   '200':
	//     description: List of keys
	//     type: string
	res, err := db.ListKeys()
	if err != nil {
		log.Print(err)

	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, fmt.Sprintf("%v items: %+v\n", len(*res), res))
}
