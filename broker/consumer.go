package broker

import (
	"log"
	"os"

	h "github.com/forumGamers/post-service-read/helper"
	"github.com/forumGamers/post-service-read/pkg/comment"
	"github.com/forumGamers/post-service-read/pkg/like"
	"github.com/forumGamers/post-service-read/pkg/post"
	"github.com/forumGamers/post-service-read/pkg/reply"
	"github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	ConsumePostCreate(postService post.PostService)
	ConsumePostDelete(postService post.PostService)
	ConsumeBulkPost(postService post.PostService)

	ConsumeLikeCreate(likeService like.LikeService)
	ConsumeLikeDelete(likeService like.LikeService)
	ConsumeBulkLike(likeService like.LikeService)

	ConsumeCommentCreate(commentService comment.CommentService)
	ConsumeCommentDelete(commentService comment.CommentService)
	ConsumeBulkComment(commentService comment.CommentService)

	ConsumeReplyCreate(replyService reply.ReplyService)
	ConsumeReplyDelete(replyService reply.ReplyService)
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

func ConsumeMessage(post post.PostService, like like.LikeService, comment comment.CommentService, replyService reply.ReplyService) {
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
