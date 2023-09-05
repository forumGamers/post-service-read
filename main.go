package main

import (
	cfg "github.com/forumGamers/post-service-read/config"
	h "github.com/forumGamers/post-service-read/helper"
	"github.com/joho/godotenv"
)

func main() {
	h.PanicIfError(godotenv.Load())

	cfg.ElasticConnection()
	cfg.BrokerConnection()
}
