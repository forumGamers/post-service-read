package main

import (
	b "github.com/forumGamers/post-service-read/broker"
	db "github.com/forumGamers/post-service-read/database"
	h "github.com/forumGamers/post-service-read/helper"
	r "github.com/forumGamers/post-service-read/routes"
	"github.com/joho/godotenv"
)

func main() {
	h.PanicIfError(godotenv.Load())

	db.ElasticConnection()
	db.CreateIndexes()

	b.BrokerConnection()
	b.Broker.ConsumePostCreate()

	r.NewRoutes()
}
