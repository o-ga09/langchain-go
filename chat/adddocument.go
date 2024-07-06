package chat

import (
	"context"
	"log"
	"os"

	"github.com/o-ga09/langchain-go/llm"
	"github.com/tmc/langchaingo/documentloaders"
)

type RequestDocumentData struct {
	PageContent string `json:"page_content"`
}

func AddDocument(ctx context.Context, data []*RequestDocumentData) error {
	chain, err := llm.New()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("tmp.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		os.Remove("tmp.txt")
		file.Close()
	}()

	for _, v := range data {
		_, err := file.WriteString(v.PageContent + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	loader := documentloaders.NewText(file)
	docs, err := loader.Load(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range docs {
		_, err := chain.AddDocument(ctx, llm.NameSpaceText, v.PageContent)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
