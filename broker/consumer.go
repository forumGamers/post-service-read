package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	doc "github.com/forumGamers/post-service-read/documents"
	h "github.com/forumGamers/post-service-read/helper"
	"github.com/rabbitmq/amqp091-go"
)

const (
	POSTEXCHANGE    = "Post-Exchange"
	NEWPOSTQUEUE    = "New-Post-Queue"
	DELETEPOSTQUEUE = "Delete-Post-Queue"
)

type Consumer interface {
	ConsumePostCreate()
}

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

func (b *ConsumerImpl) ConsumePostCreate() {
	msgs, err := b.Channel.Consume(
		NEWPOSTQUEUE,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	h.PanicIfError(err)
	postService := doc.NewPost()

	for msg := range msgs {
		var post doc.PostDocument

		if err := json.Unmarshal(msg.Body, &post); err != nil {
			fmt.Println(err.Error())
			continue
		}
		postService.Insert(context.Background(), post)
	}
}
