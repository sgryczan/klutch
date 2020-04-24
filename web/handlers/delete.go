package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sgryczan/klutch/common"
)

// DeleteHandler removes items from redis
func DeleteHandler(db *common.RedisDatastore, w http.ResponseWriter, r *http.Request) {
	// swagger:operation DELETE /id/{name} Delete Item
	//
	// Removes an item from the database
	// ---
	// consumes:
	// - text/plain
	// produces:
	// - text/plain
	// parameters:
	// - name: name
	//   in: path
	//   description: Name to be removed.
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: Successfully removed
	//     type: string
	vars := mux.Vars(r)
	item := vars["item"]

	queuedItem := common.QueueItem{
		ItemName: item,
		Status:   "Pending",
	}

	db.DeleteItem(&queuedItem)
	w.WriteHeader(http.StatusOK)
	log.Printf("Deleted item: %v", vars["item"])
}
