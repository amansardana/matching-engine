package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

var Conn *amqp.Connection

func InitConnection(address string) {
	if Conn == nil {
		conn, err := amqp.Dial(address)
		if err != nil {
			log.Fatalf("failed to open a connection: %s", err)
			panic(err)
		}
		Conn = conn
	}
}
