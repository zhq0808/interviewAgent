package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/zhq0808/interviewAgent/llm"
	"github.com/zhq0808/interviewAgent/model"
)

const MAXFailure = 3

// 记录重试次数，判断当前次数是否超过3
func recordFailure(failCount *int) bool {
	*failCount = *failCount + 1
	if *failCount > MAXFailure {
		return false
	}
	return true
}

func main() {

	client := &http.Client{}
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	dsClient := llm.NewClient(client, llm.DEEPSEEK_BASE_URL, apiKey, llm.DEEPSEEK_V4_FLASH)

	var messageArr []model.Message
	systemMessage := model.Message{Role: "system", Content: "你是我的面试助手"}
	messageArr = append(messageArr, systemMessage)

	failCount := 0

	for {
		fmt.Printf("你：")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			if recordFailure(&failCount) {
				fmt.Printf("读取用户输入失败，请重新输入：%v\n", err)
				continue
			} else {
				fmt.Printf("系统异常，已自动退出：%v\n", err)
				break
			}
		}

		userMessage := model.Message{Role: "user", Content: line}

		messageArr = append(messageArr, userMessage)

		resData, err := dsClient.SendRequest(messageArr)
		if err != nil {
			if recordFailure(&failCount) {
				fmt.Printf("读取用户输入失败，请重新输入\n")
				log.Printf("异常信息: %v\n", err)
				continue
			} else {
				fmt.Printf("系统异常，已自动退出\n")
				log.Printf("异常信息: %v\n", err)
				break
			}
		}

		var assistantOutput string
		if resData != nil && len(resData.Choices) > 0 {
			assistantOutput = resData.Choices[0].Message.Content
			fmt.Printf("面试助手：%v\n", assistantOutput)
			assistentMessage := model.Message{Role: "assistant", Content: assistantOutput}
			messageArr = append(messageArr, assistentMessage)
			failCount = 0
		} else {
			fmt.Printf("请求异常，请重试")
		}
	}
}
