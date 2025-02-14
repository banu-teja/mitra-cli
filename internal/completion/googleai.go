package completion

import (
	"context"
	"fmt"
	"log"

	"github.com/banu-teja/mitra-cli/utils"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

type GoogleAIClient struct {
	llm llms.LLM // Use the interface for flexibility
}

func NewGoogleAIClient(config utils.Config) *GoogleAIClient {
	ret := &GoogleAIClient{}
	var err error

	selectedModel := "gemini" //  Get this from user input or config.

	if modelConfig, ok := config.Models[selectedModel]; ok {
		// GoogleAI requires API key, and optionally model name.
		//  BaseURL is usually not needed for GoogleAI (it uses the default).
		ret.llm, err = googleai.New(
			context.Background(), //  Pass context here
			googleai.WithAPIKey(modelConfig.APIKey),
			googleai.WithDefaultModel(modelConfig.Name),
		)
		if err != nil {
			log.Fatalf("Failed to initialize GoogleAI LLM: %v", err)
		}

	} else {
		log.Fatalf("Model %s not found in configuration\n", selectedModel)
	}

	if err != nil {
		log.Fatal(err) //  This check is redundant after the Fatalf above.
	}

	return ret
}

func (g *GoogleAIClient) Send(ctx context.Context, prompt string) (ret string, err error) {
	// Use GenerateFromSinglePrompt for consistency.
	completion, err := llms.GenerateFromSinglePrompt(ctx, g.llm, prompt)
	if err != nil {
		// Don't Fatal here.  Return the error to the caller.
		return "", fmt.Errorf("error generating from prompt: %w", err)
	}

	return completion, nil // Return nil for the error if successful.
}
