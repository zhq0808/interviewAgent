package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Thinking Thinking  `json:"thinking"`
}

type Thinking struct {
	Type string `json:"type"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func main() {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		panic("DEEPSEEK_API_KEY is empty")
	}

	reqBody := ChatRequest{
		Model: "deepseek-v4-pro",
		Messages: []Message{
			{
				Role:    "user",
				Content: "你好，先做一下自我介绍",
			},
		},
		Thinking: Thinking{
			Type: "enabled",
		},
	}

	fmt.Printf("%+v\n", reqBody)

	bodyBytes, err := json.Marshal(reqBody)
	fmt.Println(string(bodyBytes))
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.deepseek.com/chat/completions",
		bytes.NewReader(bodyBytes),
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("status:", resp.StatusCode)
	fmt.Println("body:", string(respBytes))
}
