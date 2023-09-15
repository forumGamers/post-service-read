package broker

import (
	"context"
	"encoding/json"
	"log"
	"os"

	doc "github.com/forumGamers/post-service-read/documents"
	h "github.com/forumGamers/post-service-read/helper"
	"github.com/rabbitmq/amqp091-go"
)

const (
	POSTEXCHANGE    = "Post-Exchange"
	NEWPOSTQUEUE    = "New-Post-Queue"
	DELETEPOSTQUEUE = "Delete-Post-Queue"
	LIKEEXCHANGE    = "Like-Exchange"
	NEWLIKEQUEUE    = "New-Like-Queue"
	DELETELIKEQUEUE = "Delete-Like-Queue"
)

type Consumer interface {
	ConsumePostCreate(postService doc.PostService)
	ConsumePostDelete(postService doc.PostService)
	ConsumeLikeCreate(likeService doc.LikeService)
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

	conn, err := amqp091.DialConfig(rabbitMqServerUrl, amqp091.Config{
		Heartbeat: 10,
	})
	h.PanicIfError(err)

	ch, err := conn.Channel()
	h.PanicIfError(err)

	notifyClose := conn.NotifyClose(make(chan *amqp091.Error))
	go func() {
		retries := 0
		for {
			select {
			case err := <-notifyClose:
				if err != nil && retries < 10 {
					newConn, newErr := amqp091.DialConfig(rabbitMqServerUrl, amqp091.Config{
						Heartbeat: 10,
					})
					if newErr != nil {
						log.Printf("Gagal melakukan koneksi ulang: %s", newErr)
						continue
					}

					newCh, newErr := newConn.Channel()
					if newErr != nil {
						newConn.Close()
						log.Printf("Gagal membuat channel baru: %s", newErr)
						continue
					}

					Broker = &ConsumerImpl{
						Channel: newCh,
					}
					notifyClose = conn.NotifyClose(make(chan *amqp091.Error))
				}
				break
			}
		}
	}()

	Broker = &ConsumerImpl{
		Channel: ch,
	}
	log.Println("connection to broker success")
}

func (b *ConsumerImpl) ConsumePostCreate(postService doc.PostService) {
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

	for msg := range msgs {
		var post doc.PostDocument

		if err := json.Unmarshal(msg.Body, &post); err != nil {
			log.Println(err.Error())
			continue
		}
		go func(msg amqp091.Delivery) {
			postService.Insert(context.Background(), post)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumePostDelete(postService doc.PostService) {
	msgs, err := b.Channel.Consume(
		DELETEPOSTQUEUE,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	h.PanicIfError(err)

	for msg := range msgs {
		var post doc.PostDocument

		if err := json.Unmarshal(msg.Body, &post); err != nil {
			log.Println(err.Error())
			continue
		}
		go func(msg amqp091.Delivery) {
			postService.DeleteOneById(context.Background(), post.Id)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumeLikeCreate(likeService doc.LikeService) {
	msgs, err := b.Channel.Consume(
		NEWLIKEQUEUE,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	h.PanicIfError(err)

	for msg := range msgs {
		var like doc.LikeDocument

		if err := json.Unmarshal(msg.Body, &like); err != nil {
			log.Println(err.Error())
			continue
		}
		go func(msg amqp091.Delivery) {
			likeService.Insert(context.Background(), like)
		}(msg)
	}
}

func ConsumeMessage(post doc.PostService, like doc.LikeService) {
	go Broker.ConsumePostCreate(post)
	go Broker.ConsumePostDelete(post)

	go Broker.ConsumeLikeCreate(like)
}