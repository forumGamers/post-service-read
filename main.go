package main

import (
	b "github.com/forumGamers/post-service-read/broker"
	cfg "github.com/forumGamers/post-service-read/config"
	h "github.com/forumGamers/post-service-read/helper"
	r "github.com/forumGamers/post-service-read/routes"
	"github.com/joho/godotenv"
)

func main() {
	h.PanicIfError(godotenv.Load())

	cfg.ElasticConnection()
	b.BrokerConnection()

	r.NewRoutes()
}
