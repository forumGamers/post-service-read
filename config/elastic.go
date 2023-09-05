package config

import (
	"context"
	"fmt"
	"os"

	h "github.com/forumGamers/post-service-read/helper"
	"github.com/olivere/elastic/v7"
)

var DB *elastic.Client

func getUrl() string {
	url := os.Getenv("ELASTIC_URL")
	if url == "" {
		url = "http://localhost:9200/"
	}

	return url
}

func ElasticConnection() {
	client, err := elastic.NewClient(elastic.SetURL(getUrl()), elastic.SetSniff(false))
	h.PanicIfError(err)

	fmt.Println("database connection success")

	DB = client
}

func Ping() (*elastic.PingResult, int, error) {
	return DB.Ping(getUrl()).Do(context.Background())
}
