package main

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/clipboard"
)

type LLMChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLMRequestPayload struct {
	Model    string           `json:"model"`
	Messages []LLMChatMessage `json:"messages"`
}

type LLMResponsePayload struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// App struct
type App struct {
	ctx      context.Context
	apiUrl   string
	apiModel string
	apiKey   string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.apiUrl = "http://localhost:11434/v1"
	a.apiModel = "mistral:latest"
	a.apiKey = "ollama"
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
	config := openai.DefaultConfig(a.apiKey)
	config.BaseURL = a.apiUrl
	client := openai.NewClientWithConfig(config)

	fmt.Println("Config", config.BaseURL)

	req := openai.ChatCompletionRequest{
		Model: a.apiModel,
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
