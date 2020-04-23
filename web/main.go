package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"gitlab.com/sgryczan/go-worker-api/common"
)

// Items contains common.Plumbuses (plumbi?)
var Items = map[string]common.Plumbus{}
var DB_endpoint = os.Getenv("REDIS_ENDPOINT")
var Queue_endpoint = os.Getenv("RABBITMQ_ENDPOINT")

//var db, err =  common.Newcommon.RedisDatastore("redis:6379")

func main() {

	if DB_endpoint == "" {
		DB_endpoint = "redis"
	}

	if Queue_endpoint == "" {
		Queue_endpoint = "rabbitmq"
	}

	db, err := common.NewRedisDatastore(DB_endpoint + ":6379")
	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	rmq, err := amqp.Dial("amqp://guest:guest@" + Queue_endpoint + ":5672/")
	if err != nil {
		log.Print("Cant connect to rabbitmq")
	}
	defer rmq.Close()

	amqpChannel, err := rmq.Channel()
	if err != nil {
		log.Print("Cannot establish ampq channel")
	}
	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("items", true, false, false, false, nil)
	if err != nil {
		log.Print("could not declare `items` queue")
	}

	queConn := &common.QueueConnection{
		Channel: amqpChannel,
		Queue:   &queue,
	}

	rand.Seed(time.Now().UnixNano())

	r := mux.NewRouter()
	fmt.Println("Started v0.0.0")
	r.HandleFunc("/", homeHandler)
	r.Handle("/id/{item}", common.ItemHandler(queConn, db, AddHandler)).Methods("POST")
	//r.Handle("/id/{item}", common.RedisHandler(db, deleteHandler)).Methods("DELETE")
	r.Handle("/list", common.RedisHandler(db, listHandler)).Methods("GET")

	sh := http.StripPrefix("/api",
		http.FileServer(http.Dir("./swaggerui/")))
	r.PathPrefix("/api/").Handler(sh)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api/", 302)
	//w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "This is home")
}

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
	Items[item] = i

}

func listHandler(db *common.RedisDatastore, w http.ResponseWriter, r *http.Request) {
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

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	item := vars["item"]

	delete(Items, item)
	w.WriteHeader(http.StatusOK)
	log.Printf("Deleted item: %v", vars["item"])
}
