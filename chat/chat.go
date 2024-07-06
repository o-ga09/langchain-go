package chat

import (
	"context"
	"fmt"
	"os"

	"github.com/o-ga09/langchain-go/llm"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

type RequestQAData struct {
	Question string `json:"question,omitempty"`
}

type ResponseQAData struct {
	Kind     string `json:"kind,omitempty"`
	Question string `json:"question,omitempty"`
	Result   string `json:"result,omitempty"`
}

// RAGありver.
func RunWithRAG(ctx context.Context, meta []RequestQAData) (*[]ResponseQAData, error) {
	qaBot, err := llm.New()
	if err != nil {
		fmt.Println("failed to create LLM instance", err)
		return nil, err
	}

	targetNameSapce := []string{"html", "csv", "text"}
	response := []ResponseQAData{}
	for _, v := range meta {
		for _, nameSpace := range targetNameSapce {
			result, err := qaBot.Answer(ctx, nameSpace, v.Question)
			if result == "わかりません" {
				continue
			}
			if err != nil {
				fmt.Println("faied to response LLM", err)
				return nil, err
			}
			response = append(response, ResponseQAData{
				Kind:     nameSpace,
				Question: v.Question,
				Result:   result,
			})
		}
	}
	return &response, nil
}

// RAGなしver,
func RunWithOutRAG(ctx context.Context) (*ResponseQAData, error) {
	apiKey := os.Getenv("GOOGLEAI_API_KEY")
	opts := googleai.WithAPIKey(apiKey)
	llm, err := googleai.New(ctx, opts)
	if err != nil {
		fmt.Println("failed to create LLM instance", err)
		return nil, err
	}
	prompt := "Mハシの生年月日、SNSアカウントを教えてください"
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		fmt.Println("failed to generate completion", err)
		return nil, err
	}
	response := ResponseQAData{
		Kind:     "googleai",
		Question: prompt,
		Result:   completion,
	}
	return &response, nil
}
