package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ReqData struct {
	Model   string    `json:"model"`
	Message []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResData struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message []Message `json:"messages"`
}

func main() {
	client := &http.Client{}

	reqBody := &ReqData{
		Model: "deepseek-v4-pro",
		Message: []Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "Hello!"},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("json序列化失败:%v\n", err)
	}

	req, err := http.NewRequest("POST", "https://api.deepseek.com/chat/completions", bytes.NewReader(body))
	if err != nil {
		fmt.Printf("创建请求失败:%v\n", err)
	}

	apiKey := os.Getenv("DEEPSEEK_API_KEY")

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求异常，异常信息:%v\n", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取参数返回失败:%v\n", err)
	}

	var resData ResData
	err = json.Unmarshal(respBody, &resData)
	fmt.Printf("返回结果-原版：%v\n", resp)

	fmt.Println()

	fmt.Printf("body：%v\n", resData)
}
