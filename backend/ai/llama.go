package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Regex pattern to hide the thinking/reasoning from user
var thinkRegex = regexp.MustCompile(`(?s)<think>.*?</think>`)

func getLlamaURL() string {
	url := os.Getenv("LLAMA_URL")
	if url == "" {
		// default if missing
		return "http://localhost:8080/completion"
	}
	return url
}

type completionRequest struct {
	Prompt      string  `json:"prompt"`
	NPredict    int     `json:"n_predict"`
	Temperature float64 `json:"temperature"`
	Stream      bool    `json:"stream"`
}

type completionResponse struct {
	Content string `json:"content"`
}

// (string, error) is the return type -> so either string or error
func AskLlama(prompt string) (string, error) {
	// Build the request body
	// .Marshal converts Go struct to JSON bytes -> Unmarshal to reverse
	reqBody, err := json.Marshal(completionRequest{
		Prompt:      prompt,
		NPredict:    512,
		Temperature: 0.7,
		Stream:      false,
	})

	if err != nil {
		return "", fmt.Errorf("failed to build request: %w", err)
	}

	// Send to llama-server
	resp, err := http.Post(getLlamaURL(), "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("llama-server unreachable - Ensure SSH tunnel is open: %w", err)
	}
	// defer means "run this when the function exits"
	defer resp.Body.Close()

	// Read the raw response bytes
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Parse JSON into our struct
	var result completionResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", fmt.Errorf("bad response from llama-server: %w\nraw: %s", err, string(raw))
	}

	// Strip and clean tags before returning
	return cleanTags(result.Content), nil
}

func cleanTags(response string) string {
	cleaned := thinkRegex.ReplaceAllString(response, "")

	// Remove instruction tags
	cleaned = strings.ReplaceAll(cleaned, "[/INST]", "")
	cleaned = strings.ReplaceAll(cleaned, "[INST]", "")

	// Remove sys tags
	cleaned = strings.ReplaceAll(cleaned, "<<SYS>>", "")
	cleaned = strings.ReplaceAll(cleaned, "<</SYS>>", "")

	// trim any leftover whitespace or new lines
	return strings.TrimSpace(cleaned)
}
