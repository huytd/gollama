package main

import (
	"context"
	"fmt"
	"io"
	"os"

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

var DefaultConfig = LLMConfig{
	APIUrl:   "http://localhost:11434/v1",
	APIModel: "mistral:latest",
	APIKey:   "ollama",
}

func LoadUserConfig() LLMConfig {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return DefaultConfig
	}
	configFile := fmt.Sprintf("%s/.gollama.json", homeDir)
	file, err := os.Open(configFile)
	if err != nil {
		return DefaultConfig
	}
	defer file.Close()
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return DefaultConfig
	}
	var config LLMConfig
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return DefaultConfig
	}
	return config
}

func SaveConfig(config LLMConfig) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	configFile := fmt.Sprintf("%s/.gollama.json", homeDir)
	file, err := os.Create(configFile)
	if err != nil {
		return
	}
	defer file.Close()
	byteValue, err := json.Marshal(config)
	if err != nil {
		return
	}
	file.Write(byteValue)
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.llmConfig = LoadUserConfig()
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
	SaveConfig(a.llmConfig)
}

// Greet returns a greeting for the given name
func (a *App) GetClipboardText() string {
	err := clipboard.Init()
	if err != nil {
		return ""
	}
	return string(clipboard.Read(clipboard.FmtText))
}

func (a *App) StartLLMStream(previousMessages []string, prompt string) {
	config := openai.DefaultConfig(a.llmConfig.APIKey)
	config.BaseURL = a.llmConfig.APIUrl
	client := openai.NewClientWithConfig(config)

	messages := make([]openai.ChatCompletionMessage, len(previousMessages)+1)
	for i, message := range previousMessages {
		role := openai.ChatMessageRoleUser
		if i != 0 && i%2 == 0 {
			role = openai.ChatMessageRoleAssistant
		}
		messages[i] = openai.ChatCompletionMessage{
			Role:    role,
			Content: message,
		}
	}
	messages[len(messages)-1] = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	}

	b, _ := json.MarshalIndent(messages, "", "  ")
	fmt.Printf("MESSAGES: %+s\n", string(b))

	req := openai.ChatCompletionRequest{
		Model:    a.llmConfig.APIModel,
		Messages: messages,
		Stream:   true,
	}

	stream, err := client.CreateChatCompletionStream(a.ctx, req)
	if err != nil {
		fmt.Println("Ceate chat completion error", err)
		runtime.EventsEmit(a.ctx, "answer-update", "Something went wrong, please try again later")
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
