package elastic

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
)

type ElasticSearch struct {
	Client *elastic.Client
	Index string
	Context  context.Context
}

func NewElasticSearch(index string, mapping string) *ElasticSearch {
	e := godotenv.Load()
	if e != nil {
		log.Fatalf("Error loading .env file")
	}

	host := os.Getenv("ELASTIC_HOST")
	port := os.Getenv("ELASTIC_PORT")

	client, err := elastic.NewClient(
		elastic.SetURL("http://" + host + ":" + port),
		)

	ctx := context.TODO()

	if err != nil {
		// Handle error
		panic(err)
	}

	// Ping the Elasticsearch server to get e.g. the version number
	_, _, err = client.Ping("http://" + host + ":" + port).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(index).BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
			log.Fatalln("index was not created in  Elastic Search!")
		}

	}

	fmt.Println("Successfully connected to Elastic Search!")

	return &ElasticSearch{
		Client: client,
		Index: index,
		Context: ctx,
	}
}
