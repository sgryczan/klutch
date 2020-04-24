package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/streadway/amqp"
)

// Version holds the version at buildtime
var Version string
var dbEndpoint = os.Getenv("REDIS_ENDPOINT")
var queueEndpoint = os.Getenv("RABBITMQ_ENDPOINT")

// DBConn exports the database
var DBConn *RedisDatastore

// QueueConn exports our queue
var QueueConn *QueueConnection

// RedisDatastore houses pools
type RedisDatastore struct {
	*pool.Pool
}

// QueueConnection represents a connection to a RabbitMQ channel
type QueueConnection struct {
	Channel    *amqp.Channel
	Queue      *amqp.Queue
	Connection *amqp.Connection
}

// Queue represents a RabbitMQ queue
type Queue struct {
	*amqp.Queue
}

// QueueItem represents a queued action
type QueueItem struct {
	ItemName string
	Status   string
}

type Plumbus struct {
	Name string
}

// Cluster represents a K8S Cluster stub
type Cluster struct {
	Name string
	ID   int
}

// GetVersion returns the version
func GetVersion() string {
	return Version
}

// RedisHandler transforms incoming requests into Redis actions
func RedisHandler(c *RedisDatastore,
	f func(c *RedisDatastore, w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { f(c, w, r) })
}

// ItemHandler transforms incoming requests into a workable format
// Generally used to pass a Queue, DB and Function to a mux handler
func ItemHandler(q *QueueConnection, c *RedisDatastore,
	f func(q *QueueConnection, c *RedisDatastore, w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { f(q, c, w, r) })
}

// NewRedisDatastore returns a RedisDataStore
func NewRedisDatastore(address string) (*RedisDatastore, error) {

	connectionPool, err := pool.New("tcp", address, 10)
	if err != nil {
		return nil, err
	}
	return &RedisDatastore{
		Pool: connectionPool,
	}, nil
}

func NewQueue(endpoint string) (*QueueConnection, error) {

	connection, err := amqp.Dial(endpoint)
	if err != nil {
		log.Print("Cant connect to rabbitmq")
	}

	amqpChannel, err := connection.Channel()
	if err != nil {
		log.Print("Cannot establish ampq channel")
	}

	queue, err := amqpChannel.QueueDeclare("items", true, false, false, false, nil)
	if err != nil {
		log.Print("could not declare `items` queue")
	}

	queConn := &QueueConnection{
		Channel:    amqpChannel,
		Queue:      &queue,
		Connection: connection,
	}

	return queConn, err
}

func (q *QueueConnection) Close() {
	channel := q.Channel
	connection := q.Connection

	channel.Close()
	connection.Close()
}

// CreateItem creates an object in the database
func (r *RedisDatastore) CreateItem(item *QueueItem) error {

	itemJSON, err := json.Marshal(*item)
	if err != nil {
		return err
	}

	if r.Cmd("SET", item.ItemName, string(itemJSON)).Err != nil {
		return errors.New("Failed to execute Redis SET command")
	}

	return nil
}

// DeleteItem deletes an object in the database
func (r *RedisDatastore) DeleteItem(item *QueueItem) error {

	if r.Cmd("DEL", item.ItemName).Err != nil {
		return errors.New("Failed to execute Redis DEL command")
	}

	return nil
}

// GetItem returns an object in the database
func (r *RedisDatastore) GetItem(item string) (*Plumbus, error) {

	exists, err := r.Cmd("EXISTS", "item:"+item).Int()

	if err != nil {
		return nil, err
	} else if exists == 0 {
		return nil, nil
	}

	var p Plumbus

	itemJSON, err := r.Cmd("GET", "item:"+item).Str()

	if err != nil {
		log.Print(err)

		return nil, err
	}

	if err := json.Unmarshal([]byte(itemJSON), &p); err != nil {
		log.Print(err)
		return nil, err
	}

	return &p, nil
}

// ListKeys returns all keys in the database (#WONTSCALE)
func (r *RedisDatastore) ListKeys() (*[]Plumbus, error) {
	keysJSON, err := r.Cmd("KEYS", "*").Array()
	fmt.Println(keysJSON)
	var ps []Plumbus

	if err != nil {
		log.Print(err)

		return nil, err
	}

	for i := range keysJSON {
		key, _ := keysJSON[i].Str()
		p := Plumbus{
			Name: key,
		}
		ps = append(ps, p)
	}
	/* if err := json.Unmarshal(keysJSON, &p); err != nil {
		log.Print(err)
		return nil, err
	} */

	return &ps, nil
}

// Close the connection to Redis
func (r *RedisDatastore) Close() {

	r.Empty()
}
