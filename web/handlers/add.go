package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sgryczan/klutch/common"
	"github.com/streadway/amqp"
)

func AddHandler(q *common.QueueConnection, db *common.RedisDatastore, w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /id/{name} Add Item
	//
	// Adds an item to the database
	// ---
	// consumes:
	// - text/plain
	// produces:
	// - text/plain
	// parameters:
	// - name: name
	//   in: path
	//   description: Name to be added.
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: Add an item to the database
	//     type: string
	vars := mux.Vars(r)
	item := vars["item"]

	w.WriteHeader(http.StatusOK)
	i := common.Plumbus{
		Name: item,
	}
	queuedItem := common.QueueItem{
		ItemName: i.Name,
		Status:   "Pending",
	}
	fmt.Printf("%+v\n", i)

	body, err := json.Marshal(i)
	if err != nil {
		log.Print("Error marshalling json")
	}
	queue := q.Queue
	channel := q.Channel

	err = channel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
	if err != nil {
		log.Print("Error adding item to queue!")
	}

	err = db.CreateItem(&queuedItem)

	if err != nil {
		log.Print(err)
	}

	fmt.Fprintf(w, "Added item: %v\n", item)
	log.Print(fmt.Sprintf("added item: %v\n", item))

}
