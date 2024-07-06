package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/o-ga09/langchain-go/llm"
	"github.com/tmc/langchaingo/documentloaders"
)

func main() {
	ctx := context.Background()
	chain, err := llm.New()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("./tools/insert_csv/qa.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	loader := documentloaders.NewCSV(file)
	docs, err := loader.Load(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range docs {
		_, err := chain.AddDocument(ctx, llm.NameSpaceCSV, v.PageContent)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("done")
}
