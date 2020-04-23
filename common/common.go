package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/streadway/amqp"
)

type RedisDatastore struct {
	*pool.Pool
}

type QueueConnection struct {
	Channel *amqp.Channel
	Queue   *amqp.Queue
}

type Queue struct {
	*amqp.Queue
}

type QueueItem struct {
	ItemName string
	Status   string
}

type Plumbus struct {
	Name string
}

func RedisHandler(c *RedisDatastore,
	f func(c *RedisDatastore, w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { f(c, w, r) })
}

func ItemHandler(q *QueueConnection, c *RedisDatastore,
	f func(q *QueueConnection, c *RedisDatastore, w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { f(q, c, w, r) })
}

func NewRedisDatastore(address string) (*RedisDatastore, error) {

	connectionPool, err := pool.New("tcp", address, 10)
	if err != nil {
		return nil, err
	}
	return &RedisDatastore{
		Pool: connectionPool,
	}, nil
}

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

func (r *RedisDatastore) ListKeys() (*[]Plumbus, error) {
	keysJSON, err := r.Cmd("KEYS", "*").Array()
	fmt.Println(keysJSON)
	var ps []Plumbus

	if err != nil {
		log.Print(err)

		return nil, err
	}

	for i, _ := range keysJSON {
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

func (r *RedisDatastore) Close() {

	r.Empty()
}
