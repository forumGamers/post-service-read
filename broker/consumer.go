package broker

import (
	"log"
	"os"

	doc "github.com/forumGamers/post-service-read/documents"
	h "github.com/forumGamers/post-service-read/helper"
	"github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	ConsumePostCreate(postService doc.PostService)
	ConsumePostDelete(postService doc.PostService)
	ConsumeBulkPost(postService doc.PostService)
	ConsumeLikeCreate(likeService doc.LikeService)
	ConsumeLikeDelete(likeService doc.LikeService)
	ConsumeBulkLike(likeService doc.LikeService)
	ConsumeCommentCreate(commentService doc.CommentService)
	ConsumeCommentDelete(commentService doc.CommentService)
	ConsumeBulkComment(commentService doc.CommentService)
	ConsumeReplyCreate(replyService doc.ReplyService)
	ConsumeReplyDelete(replyService doc.ReplyService)
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

func ConsumeMessage(post doc.PostService, like doc.LikeService, comment doc.CommentService, replyService doc.ReplyService) {
	//kirim otentikasi header
	go Broker.ConsumePostCreate(post)
	go Broker.ConsumePostDelete(post)
	go Broker.ConsumeBulkPost(post)

	go Broker.ConsumeLikeCreate(like)
	go Broker.ConsumeLikeDelete(like)
	go Broker.ConsumeBulkLike(like)

	go Broker.ConsumeCommentCreate(comment)
	go Broker.ConsumeCommentDelete(comment)
	go Broker.ConsumeBulkComment(comment)

	go Broker.ConsumeReplyCreate(replyService)
	go Broker.ConsumeReplyDelete(replyService)
}
