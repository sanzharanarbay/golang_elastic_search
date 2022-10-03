package post

import (
	"os"
)

type ElasticPostConfig struct {
	Index   string
	Mapping string
}

func NewPostElasticConfig() *ElasticPostConfig {
	index := os.Getenv("ELASTIC_INDEX")
	mapping := `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 1
	},
	"mappings":{
			"properties":{
				"title":{
					"type":"keyword"
				},
				"content":{
					"type":"text"
				},
				"category_id":{
					"type":"long"
				},
				"created_at":{
					"type":"date"
				},
				"updated_at":{
					"type":"date"
				}
			}
	}
}`
	return &ElasticPostConfig{
		Index:   index,
		Mapping: mapping,
	}
}
