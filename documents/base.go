package documents

import "github.com/olivere/elastic/v7"

type BaseDocument struct {
	DB *elastic.Client
}
