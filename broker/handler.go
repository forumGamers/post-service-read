package broker

import (
	"context"
	"encoding/json"

	h "github.com/forumGamers/post-service-read/helper"
	"github.com/forumGamers/post-service-read/pkg/comment"
	"github.com/forumGamers/post-service-read/pkg/like"
	"github.com/forumGamers/post-service-read/pkg/post"
	"github.com/forumGamers/post-service-read/pkg/reply"
	"github.com/rabbitmq/amqp091-go"
)

func (b *ConsumerImpl) ConsumePostCreate(postService post.PostService) {
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
			var post post.PostDocument

			json.Unmarshal(msg.Body, &post)
			postService.Insert(context.Background(), post)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumePostDelete(postService post.PostService) {
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
			var post post.PostDocument

			json.Unmarshal(msg.Body, &post)
			postService.DeleteOneById(context.Background(), post.Id)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumeLikeCreate(likeService like.LikeService) {
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
			var like like.LikeDocument

			json.Unmarshal(msg.Body, &like)
			likeService.Insert(context.Background(), like)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumeLikeDelete(likeService like.LikeService) {
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
			var like like.LikeDocument

			json.Unmarshal(msg.Body, &like)
			likeService.DeleteOneById(context.Background(), like.Id)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumeCommentCreate(commentService comment.CommentService) {
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
			var comment comment.CommentDocument

			json.Unmarshal(msg.Body, &comment)
			commentService.Insert(context.Background(), comment)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumeCommentDelete(commentService comment.CommentService) {
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
			var comment comment.CommentDocument

			json.Unmarshal(msg.Body, &comment)
			commentService.DeleteOneById(context.Background(), comment.Id)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumeBulkComment(commentService comment.CommentService) {
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
			var comments []comment.CommentDocument

			json.Unmarshal(msg.Body, &comments)
			commentService.BulkCreate(context.Background(), comments)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumeReplyCreate(replyService reply.ReplyService) {
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
			var comment reply.Reply

			json.Unmarshal(msg.Body, &comment)
			replyService.Insert(context.Background(), comment)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumeReplyDelete(replyService reply.ReplyService) {
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
			var comment reply.Reply

			json.Unmarshal(msg.Body, &comment)
			replyService.DeleteOneById(context.Background(), comment.Id)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumeBulkPost(postService post.PostService) {
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
			var posts []post.PostDocument

			json.Unmarshal(msg.Body, &posts)
			postService.BulkCreate(context.Background(), posts)
		}(msg)
	}
}

func (b *ConsumerImpl) ConsumeBulkLike(likeService like.LikeService) {
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
			var likes []like.LikeDocument

			json.Unmarshal(msg.Body, &likes)
			likeService.BulkCreate(context.Background(), likes)
		}(msg)
	}
}
