package main

import (
	"context"
	"fmt"
	"log"

	"github.com/o-ga09/langchain-go/llm"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

func main() {
	ctx := context.Background()
	weaviateClient := weaviate.New(weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	})

	if ok, err := weaviateClient.Schema().ClassExistenceChecker().WithClassName(llm.WeaviateIndexName).Do(ctx); ok {
		log.Print("Already Exists")
	} else if err != nil {
		log.Fatal(err)
	}

	if err := weaviateClient.Schema().ClassCreator().WithClass(&models.Class{
		Class:       llm.WeaviateIndexName,
		Description: "qa class",
		VectorIndexConfig: map[string]any{
			"distance": "cosine",
		},
		ModuleConfig: map[string]any{},
		Properties: []*models.Property{
			{
				Name:        llm.WeaviatePropertyTextName,
				Description: "document text",
				DataType:    []string{"text"},
			},
			{
				Name:        llm.WeaviatePropertyNameSpaceName,
				Description: "namespace",
				DataType:    []string{"text"},
			},
		},
	}).Do(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("created")
}
