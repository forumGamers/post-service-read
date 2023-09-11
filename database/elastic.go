package database

import (
	"context"
	"fmt"
	"os"

	h "github.com/forumGamers/post-service-read/helper"
	"github.com/olivere/elastic/v7"
)

var DB *elastic.Client

const (
	POSTINDEX = "post-service.post"
)

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

func CreateIndexes() {
	for _, index := range []string{POSTINDEX} {
		if exists, _ := DB.IndexExists(index).Do(context.Background()); !exists {
			var schema string
			switch index {
			case POSTINDEX:
				schema = PostMapping
			default:
				schema = ""
			}
			if created, err := DB.CreateIndex(index).BodyString(schema).Do(context.Background()); err != nil || !created.Acknowledged {
				h.PanicIfError(h.InternalServer)
			}
		}
	}
}
