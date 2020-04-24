module worker

go 1.13

require (
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/sgryczan/klutch/common v0.0.0
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
)

replace github.com/sgryczan/klutch/common => ../common
