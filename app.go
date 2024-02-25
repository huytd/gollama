package main

import (
	"context"
	"fmt"

	"encoding/json"

	openai "github.com/sashabaranov/go-openai"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/clipboard"
)

type LLMConfig struct {
	APIUrl   string
	APIModel string
	APIKey   string
}

// App struct
type App struct {
	ctx       context.Context
	llmConfig LLMConfig
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.llmConfig = LLMConfig{
		APIUrl:   "http://localhost:11434/v1",
		APIModel: "mistral:latest",
		APIKey:   "ollama",
	}
}

func (a *App) GetLLMConfig() string {
	configJSON, err := json.Marshal(a.llmConfig)
	if err != nil {
		return "{}"
	}
	return string(configJSON)

}

func (a *App) SetLLMConfig(apiUrl string, apiKey string, apiModel string) {
	a.llmConfig.APIUrl = apiUrl
	a.llmConfig.APIKey = apiKey
	a.llmConfig.APIModel = apiModel
}

// Greet returns a greeting for the given name
func (a *App) GetClipboardText() string {
	err := clipboard.Init()
	if err != nil {
		return ""
	}
	return string(clipboard.Read(clipboard.FmtText))
}

func (a *App) StartLLMStream(prompt string) {
	config := openai.DefaultConfig(a.llmConfig.APIKey)
	config.BaseURL = a.llmConfig.APIUrl
	client := openai.NewClientWithConfig(config)

	fmt.Println("Config", config.BaseURL)

	req := openai.ChatCompletionRequest{
		Model: a.llmConfig.APIModel,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Stream: true,
	}

	stream, err := client.CreateChatCompletionStream(a.ctx, req)
	if err != nil {
		fmt.Println("Ceate chat completion error", err)
		runtime.EventsEmit(a.ctx, "stream-done", "")
		return
	}
	defer stream.Close()

	answer := ""
	for {
		response, _ := stream.Recv()
		if len(response.Choices) > 0 {
			answer += response.Choices[0].Delta.Content
			runtime.EventsEmit(a.ctx, "answer-update", answer)
		} else {
			runtime.EventsEmit(a.ctx, "stream-done", "")
			return
		}
	}
}
