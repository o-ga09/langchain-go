package main

import (
	"context"
	"fmt"
	"os"

	"github.com/o-ga09/langchain-go/api"
	"github.com/o-ga09/langchain-go/llm"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

type QA_MetaData struct {
	Question  string
	NameSpace string
}

func main() {
	server := api.New()
	if err := server.Run(context.Background()); err != nil {
		fmt.Println("failed to run server", err)
	}

	// ctx := context.Background()
	// meta := []QA_MetaData{
	// 	{Question: "Mハシの生年月日、SNSアカウントを教えてください", NameSpace: "csv"},
	// }

	// if err := RunWithRAG(ctx, meta); err != nil {
	// 	fmt.Println("failed to run with RAG", err)
	// }

	// fmt.Println("=====================================")

	// if err := RunWithOutRAG(ctx); err != nil {
	// 	fmt.Println("failed to run with out RAG", err)
	// }
}

// RAGありver.
func RunWithRAG(ctx context.Context, meta []QA_MetaData) error {
	qaBot, err := llm.New()
	if err != nil {
		fmt.Println("failed to create LLM instance", err)
		return err
	}

	for _, v := range meta {
		result, err := qaBot.Answer(ctx, v.NameSpace, v.Question)
		if err != nil {
			fmt.Println("faied to response LLM", err)
			return err
		}

		fmt.Println("=====================================")
		fmt.Printf("kind:\n %s\n", v.NameSpace)
		fmt.Printf("question:\n %s\n", v.Question)
		fmt.Printf("result:\n %s\n", result)
		fmt.Println("=====================================")
	}
	return nil
}

// RAGなしver,
func RunWithOutRAG(ctx context.Context) error {
	apiKey := os.Getenv("GOOGLEAI_API_KEY")
	opts := googleai.WithAPIKey(apiKey)
	llm, err := googleai.New(ctx, opts)
	if err != nil {
		fmt.Println("failed to create LLM instance", err)
		return err
	}
	prompt := "Mハシの生年月日、SNSアカウントを教えてください"
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		fmt.Println("failed to generate completion", err)
		return err
	}
	fmt.Println("===answer===")
	fmt.Println(completion)
	fmt.Println("============")
	return nil
}
