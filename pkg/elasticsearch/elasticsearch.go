package elasticsearch

import (
	"log"

	v8 "github.com/elastic/go-elasticsearch/v8"
)

func NewElasticClient() *v8.Client {
	client, err := v8.NewDefaultClient()
	if err != nil {
		log.Fatal("Failed to open connection to elasticsearch", err)
		return nil
	}
	return client
}
