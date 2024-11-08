package completion

import (
	"context"
	"log"

	"github.com/banu-teja/mitra-cli/utils"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
)

type Client struct {
	llm *anthropic.LLM
}

func NewClient(config utils.Config) *Client {
	ret := &Client{}
	var err error

	selectedModel := "anthropic" // This could come from user input

	if modelConfig, ok := config.Models[selectedModel]; ok {
		ret.llm, err = anthropic.New(
			anthropic.WithBaseURL(modelConfig.BaseURL),
			anthropic.WithModel(modelConfig.Name),
			anthropic.WithToken(modelConfig.APIKey),
		)
	} else {
		log.Fatalf("Model %s not found in configuration\n", selectedModel)
	}

	if err != nil {
		log.Fatal(err)
	}

	return ret
}

func (an *Client) Send(ctx context.Context, prompt string) (ret string, err error) {
	completion, err := llms.GenerateFromSinglePrompt(ctx, an.llm, prompt)
	if err != nil {
		log.Fatal(err)
	}

	return completion, err

}
