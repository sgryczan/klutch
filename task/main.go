package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"gitlab.com/sgryczan/go-worker-api/common"
	"log"
	"os"
)

type plumbus struct {
	Name string
}

// Items contains plumbuses (plumbi?)
var Items = map[string]plumbus{}
var DB_endpoint = os.Getenv("REDIS_ENDPOINT")
var Queue_endpoint = os.Getenv("RABBITMQ_ENDPOINT")

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}

}

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

	conn, err := amqp.Dial("amqp://guest:guest@" + Queue_endpoint + ":5672/")
	handleError(err, "Can't connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("items", true, false, false, false, nil)
	handleError(err, "Could not declare `items` queue")

	err = amqpChannel.Qos(1, 0, false)
	handleError(err, "Could not configure QoS")

	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Could not register consumer")

	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {
			log.Printf("Received a message: %s", d.Body)

			item := &plumbus{}

			err := json.Unmarshal(d.Body, item)

			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			log.Printf("Processed item: %v", item)

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}

		}
	}()

	// Stop for program termination
	<-stopChan

}
