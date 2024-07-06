package llm

import (
	"context"
	"os"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/weaviate"
)

const (
	WeaviateIndexName             = "Qa"
	WeaviatePropertyTextName      = "text"
	WeaviatePropertyNameSpaceName = "namespace"
	NameSpaceHTML                 = "html"
	NameSpaceCSV                  = "csv"
	NameSpaceText                 = "text"
)

type LLM struct {
	llm   llms.Model
	store vectorstores.VectorStore
}

func New() (*LLM, error) {
	ctx := context.Background()
	apiKey := os.Getenv("GOOGLEAI_API_KEY")
	opts := googleai.WithAPIKey(apiKey)
	llm, err := googleai.New(ctx, opts)
	if err != nil {
		return nil, err
	}

	e, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, err
	}
	store, err := weaviate.New(
		weaviate.WithScheme("http"),         // docker-composeの設定に合わせる
		weaviate.WithHost("localhost:8080"), // docker-composeの設定に合わせる
		weaviate.WithEmbedder(e),
		weaviate.WithIndexName(WeaviateIndexName),
		weaviate.WithTextKey(WeaviatePropertyTextName),
		weaviate.WithNameSpaceKey(WeaviatePropertyNameSpaceName),
	)
	if err != nil {
		return nil, err
	}
	return &LLM{
		llm:   llm,
		store: store,
	}, nil
}

func (l *LLM) AddDocument(ctx context.Context, namespace string, content string) ([]string, error) {
	return l.store.AddDocuments(ctx, []schema.Document{
		{
			PageContent: content,
		},
	}, vectorstores.WithNameSpace(namespace))
}

func (l *LLM) Answer(ctx context.Context, namespace string, question string) (string, error) {
	prompt := prompts.NewPromptTemplate(
		`## Introduction 
			あなたはカスタマーサポートです。丁寧な回答を心がけてください。
			以下のContextを使用して、日本語で質問に答えてください。Contextから答えがわからない場合は、「わかりません」と回答してください。

			## 質問
			{{.question}}

			## Context
			{{.context}}

			日本語での回答:`,
		[]string{"context", "question"},
	)

	combineChain := chains.NewStuffDocuments(chains.NewLLMChain(l.llm, prompt))
	result, err := chains.Run(
		ctx,
		chains.NewRetrievalQA(
			combineChain,
			vectorstores.ToRetriever(
				l.store,
				5,
				vectorstores.WithNameSpace(string(namespace)),
			),
		),
		question,
		chains.WithModel("gemini-pro"),
	)
	if err != nil {
		return "", err
	}
	return result, nil
}
