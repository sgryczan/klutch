package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sgryczan/klutch/common"
	"github.com/sgryczan/klutch/web/handlers"
	"github.com/streadway/amqp"
)

// Items contains common.Plumbuses (plumbi?)
var Items = map[string]common.Plumbus{}
var DBEndpoint = os.Getenv("REDIS_ENDPOINT")
var QueueEndpoint = os.Getenv("RABBITMQ_ENDPOINT")
var version string

//var db, err =  common.Newcommon.RedisDatastore("redis:6379")

func main() {

	// Default values for testing in Docker
	// remove at some point
	if DBEndpoint == "" {
		DBEndpoint = "redis"
	}
	if QueueEndpoint == "" {
		QueueEndpoint = "rabbitmq"
	}
	//

	db, err := common.NewRedisDatastore(DBEndpoint + ":6379")
	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	rmq, err := amqp.Dial("amqp://guest:guest@" + QueueEndpoint + ":5672/")
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
	fmt.Println("Started " + common.GetVersion())

	r.HandleFunc("/", handlers.HomeHandler)
	r.HandleFunc("/about", handlers.AboutHandler)
	r.Handle("/id/{item}", common.ItemHandler(queConn, db, handlers.AddHandler)).Methods("POST")
	r.Handle("/id/{item}", common.RedisHandler(db, handlers.DeleteHandler)).Methods("DELETE")
	r.Handle("/list", common.RedisHandler(db, handlers.ListHandler)).Methods("GET")

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
