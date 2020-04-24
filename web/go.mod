module github.com/sgryczan/klutch/web

go 1.13

require (
	github.com/gorilla/mux v1.7.3
	github.com/sgryczan/klutch/common v0.0.0
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
)

replace github.com/sgryczan/klutch/common => ../common
