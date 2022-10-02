package elastic

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"github.com/elastic/go-elasticsearch/v7"
	"strconv"
)

type ElasticSearch struct {
	client *elasticsearch.Client
	index  string
	alias  string
}

func NewElasticSearch() *ElasticSearch {
	e := godotenv.Load()
	if e != nil {
		log.Fatalf("Error loading .env file")
	}

	username := os.Getenv("ELASTIC_USER")
	password := os.Getenv("ELASTIC_PASSWORD")
	host := os.Getenv("ELASTIC_HOST")
	port := os.Getenv("ELASTIC_PORT")

	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://" + host + ":" + port,
		},
		Username: username,
		Password: password,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully connected to Elastic Search!")

	return &ElasticSearch{
		client: client,
	}
}

func (e *ElasticSearch) CreateIndex() error {
	env := godotenv.Load()
	if env != nil {
		log.Fatalf("Error loading .env file")
	}
	index := os.Getenv("ELASTIC_INDEX")

	e.index = index
	e.alias = index + "_alias"

	res, err := e.client.Indices.Exists([]string{e.index})
	if err != nil {
		return fmt.Errorf("cannot check index existence: %w", err)
	}
	fmt.Println("statusCode - " + strconv.Itoa(res.StatusCode))
	if res.StatusCode == 200 {
		return nil
	}
	if res.StatusCode != 404 {
		return fmt.Errorf("error in index existence response: %s", res.String())
	}

	res, err = e.client.Indices.Create(e.index)
	if err != nil {
		return fmt.Errorf("cannot create index: %w", err)
	}
	fmt.Println("Index created Successfully")
	if res.IsError() {
		return fmt.Errorf("error in index creation response: %s", res.String())
	}

	res, err = e.client.Indices.PutAlias([]string{e.index}, e.alias)
	if err != nil {
		return fmt.Errorf("cannot create index alias: %w", err)
	}
	if res.IsError() {
		return fmt.Errorf("error in index alias creation response: %s", res.String())
	}

	return nil
}

// document represents a single document in Get API response body.
type document struct {
	Source interface{} `json:"_source"`
}
