package main

import (
	b "github.com/forumGamers/post-service-read/broker"
	"github.com/forumGamers/post-service-read/controllers"
	db "github.com/forumGamers/post-service-read/database"
	"github.com/forumGamers/post-service-read/documents"
	h "github.com/forumGamers/post-service-read/helper"
	"github.com/forumGamers/post-service-read/pkg/like"
	"github.com/forumGamers/post-service-read/pkg/post"
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
	commentService := documents.NewComment()
	replyService := documents.NewReply()

	b.ConsumeMessage(
		postService,
		likeService,
		commentService,
		replyService,
	)

	r.NewRoutes(
		controllers.NewPostController(postService),
		controllers.NewCommentController(commentService),
	)
}
