package assistant

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

const endpoint = "https://api.groq.com/openai/v1/chat/completions"

type Assistant struct {
	apiKey string
	model  string
	client *http.Client
}

type chatRequest struct {
	Model               string        `json:"model"`
	Messages            []chatMessage `json:"messages"`
	Temperature         float64       `json:"temperature"`
	MaxCompletionTokens int           `json:"max_completion_tokens"`
	ReasoningEffort     string        `json:"reasoning_effort"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

func New() (*Assistant, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, errors.New("GROQ_API_KEY is not set")
	}

	return &Assistant{
		apiKey: apiKey,
		model:  "openai/gpt-oss-120b",
		client: &http.Client{Timeout: 60 * time.Second},
	}, nil
}

func (a *Assistant) Reply(ctx context.Context, message string) (string, error) {
	payload := chatRequest{
		Model:               a.model,
		Temperature:         1,
		MaxCompletionTokens: 2048,
		ReasoningEffort:     "medium",
		Messages: []chatMessage{
			{
				Role:    "system",
				Content: "You are Noah, a helpful and concise personal AI assistant.",
			},
			{
				Role:    "user",
				Content: message,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, endpoint, bytes.NewReader(body),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("Groq API returned HTTP %d", resp.StatusCode)
	}

	var result chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", errors.New("Groq returned no response choices")
	}

	return result.Choices[0].Message.Content, nil
}
