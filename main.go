package main

import (
	b "github.com/forumGamers/post-service-read/broker"
	"github.com/forumGamers/post-service-read/controllers"
	db "github.com/forumGamers/post-service-read/database"
	h "github.com/forumGamers/post-service-read/helper"
	"github.com/forumGamers/post-service-read/pkg/comment"
	"github.com/forumGamers/post-service-read/pkg/like"
	"github.com/forumGamers/post-service-read/pkg/post"
	"github.com/forumGamers/post-service-read/pkg/reply"
	"github.com/forumGamers/post-service-read/pkg/share"
	r "github.com/forumGamers/post-service-read/routes"
	"github.com/joho/godotenv"
)

func main() {
	h.PanicIfError(godotenv.Load())

	db.ElasticConnection()
	db.CreateIndexes()
	db.CreateAliases()

	b.BrokerConnection()
	postService := post.NewPost()
	likeService := like.NewLike()
	commentService := comment.NewComment()
	replyService := reply.NewReply()
	shareService := share.NewShare()

	b.ConsumeMessage(
		postService,
		likeService,
		commentService,
		replyService,
	)

	r.NewRoutes(
		controllers.NewPostController(postService, likeService, commentService, replyService, shareService),
		controllers.NewCommentController(commentService),
	)
}
