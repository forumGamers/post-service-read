package broker

import (
	"context"
	"encoding/json"

	doc "github.com/forumGamers/post-service-read/documents"
	h "github.com/forumGamers/post-service-read/helper"
	"github.com/rabbitmq/amqp091-go"
)

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
		go func(msg amqp091.Delivery) {
			var post doc.PostDocument

			json.Unmarshal(msg.Body, &post)
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
		go func(msg amqp091.Delivery) {
			var post doc.PostDocument

			json.Unmarshal(msg.Body, &post)
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
		go func(msg amqp091.Delivery) {
			var like doc.LikeDocument

			json.Unmarshal(msg.Body, &like)
			likeService.Insert(context.Background(), like)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumeLikeDelete(likeService doc.LikeService) {
	msgs, err := b.Channel.Consume(
		DELETELIKEQUEUE,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	h.PanicIfError(err)

	for msg := range msgs {
		go func(msg amqp091.Delivery) {
			var like doc.LikeDocument

			json.Unmarshal(msg.Body, &like)
			likeService.DeleteOneById(context.Background(), like.Id)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumeCommentCreate(commentService doc.CommentService) {
	msgs, err := b.Channel.Consume(
		NEWCOMMENTQUEUE,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	h.PanicIfError(err)

	for msg := range msgs {
		go func(msg amqp091.Delivery) {
			var comment doc.CommentDocument

			json.Unmarshal(msg.Body, &comment)
			commentService.Insert(context.Background(), comment)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumeCommentDelete(commentService doc.CommentService) {
	msgs, err := b.Channel.Consume(
		DELETECOMMENTQUEUE,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	h.PanicIfError(err)

	for msg := range msgs {
		go func(msg amqp091.Delivery) {
			var comment doc.CommentDocument

			json.Unmarshal(msg.Body, &comment)
			commentService.DeleteOneById(context.Background(), comment.Id)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumeBulkComment(commentService doc.CommentService) {
	msgs, err := b.Channel.Consume(
		BULKCOMMENTQUEUE,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	h.PanicIfError(err)

	for msg := range msgs {
		go func(msg amqp091.Delivery) {
			var comments []doc.CommentDocument

			json.Unmarshal(msg.Body, &comments)
			commentService.BulkCreate(context.Background(), comments)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumeReplyCreate(replyService doc.ReplyService) {
	msgs, err := b.Channel.Consume(
		NEWREPLYQUEUE,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	h.PanicIfError(err)

	for msg := range msgs {
		go func(msg amqp091.Delivery) {
			var comment doc.ReplyCommentDocument

			json.Unmarshal(msg.Body, &comment)
			replyService.Insert(context.Background(), comment)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumeReplyDelete(replyService doc.ReplyService) {
	msgs, err := b.Channel.Consume(
		DELETEREPLYQUEUE,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	h.PanicIfError(err)

	for msg := range msgs {
		go func(msg amqp091.Delivery) {
			var comment doc.ReplyCommentDocument

			json.Unmarshal(msg.Body, &comment)
			replyService.DeleteOneById(context.Background(), comment.Id)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumeBulkPost(postService doc.PostService) {
	msgs, err := b.Channel.Consume(
		BULKPOSTQUEUE,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	h.PanicIfError(err)

	for msg := range msgs {
		go func(msg amqp091.Delivery) {
			var posts []doc.PostDocument

			json.Unmarshal(msg.Body, &posts)
			postService.BulkCreate(context.Background(), posts)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumeBulkLike(likeService doc.LikeService) {
	msgs, err := b.Channel.Consume(
		BULKLIKEQUEUE,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	h.PanicIfError(err)

	for msg := range msgs {
		go func(msg amqp091.Delivery) {
			var likes []doc.LikeDocument

			json.Unmarshal(msg.Body, &likes)
			likeService.BulkCreate(context.Background(), likes)
		}(msg)
	}
}
