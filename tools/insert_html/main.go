package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/o-ga09/langchain-go/llm"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
)

func main() {
	ctx := context.Background()
	chain, err := llm.New()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("./tools/insert_html/qa.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	loader := documentloaders.NewHTML(file)
	docs, err := loader.LoadAndSplit(
		ctx,
		textsplitter.RecursiveCharacter{
			Separators:    []string{"\n\n", "\n", " ", ""},
			ChunkSize:     3200,
			ChunkOverlap:  800,
			LenFunc:       func(s string) int { return len(s) },
			KeepSeparator: true,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range docs {
		_, err := chain.AddDocument(ctx, llm.NameSpaceHTML, v.PageContent)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("done")
}
