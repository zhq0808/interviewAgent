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

type RequestData struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func main() {
	//1.读环境变量里面的key（通过os读的环境变量）
	apiKey := os.Getenv("DEEPSEEK_API_KEY")

	//2.初始化一个client？自定义设置关键参数
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	//3.构造请求参数的结构体（此时是golang对象）
	req := &RequestData{
		Model: "deepseek-v4-pro",
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是一个帮助助手.",
			}, {
				Role:    "user",
				Content: "你好!",
			},
		},
	}
	//4.将golang对象转化为json格式（deepseek入参需要json格式）
	bodyReqBytes, err := json.Marshal(req)
	if err != nil {
		panic("golang对象转化为json失败")
	}

	url := "https://api.deepseek.com/chat/completions"

	//5.创建一个http请求对象，注意入参要求body io.Reader，需要转化为带read方法
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyReqBytes))

	//6.拼接header
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+apiKey)

	//7.发起http请求到deepseek
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("请求异常")
		return
	}
	//7.1结束时关闭连接（小demo暂时不涉及）
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Println(string(bodyBytes))
}
