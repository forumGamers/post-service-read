package broker

import (
	"fmt"
	"os"

	h "github.com/forumGamers/post-service-read/helper"
	"github.com/rabbitmq/amqp091-go"
)

type Consumer interface{}

type ConsumerImpl struct {
	Channel *amqp091.Channel
}

var Broker Consumer

func BrokerConnection() {
	rabbitMqServerUrl := os.Getenv("RABBITMQURL")

	if rabbitMqServerUrl == "" {
		rabbitMqServerUrl = "amqp://user:password@localhost:5672"
	}

	conn, err := amqp091.Dial(rabbitMqServerUrl)
	h.PanicIfError(err)

	ch, err := conn.Channel()
	h.PanicIfError(err)

	notifyClose := conn.NotifyClose(make(chan *amqp091.Error))
	go func() { <-notifyClose }()

	Broker = &ConsumerImpl{
		Channel: ch,
	}
	fmt.Println("connection to broker success")
}
