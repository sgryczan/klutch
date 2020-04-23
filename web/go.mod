module worker

go 1.13

require (
	github.com/gorilla/mux v1.7.3
	github.com/mediocregopher/radix.v2 v0.0.0-20181115013041-b67df6e626f9 // indirect
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	gitlab.com/sgryczan/go-worker-api/common v0.0.0
)

replace gitlab.com/sgryczan/go-worker-api/common => ../common
