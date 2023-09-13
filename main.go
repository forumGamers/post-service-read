package main

import (
	b "github.com/forumGamers/post-service-read/broker"
	"github.com/forumGamers/post-service-read/controllers"
	db "github.com/forumGamers/post-service-read/database"
	"github.com/forumGamers/post-service-read/documents"
	h "github.com/forumGamers/post-service-read/helper"
	r "github.com/forumGamers/post-service-read/routes"
	"github.com/joho/godotenv"
)

func main() {
	h.PanicIfError(godotenv.Load())

	db.ElasticConnection()
	db.CreateIndexes()

	b.BrokerConnection()
	postService := documents.NewPost()

	go b.Broker.ConsumePostDelete(postService)
	go b.Broker.ConsumePostCreate(postService)

	postController := controllers.NewPostController(postService)

	r.NewRoutes(postController)
}
