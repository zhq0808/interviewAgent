package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/zhq0808/interviewAgent/model"
)

const deepseekBaseUrl = "https://api.deepseek.com/chat/completions"

const HTTP_POST = "POST"

func sendDeepSeekRequest(client http.Client, body []byte, messageArr []model.Message) error {
	reqBody := &model.ReqData{
		Model:   "deepseek-v4-flash",
		Message: messageArr,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// 创建请求
	req, err := http.NewRequest(HTTP_POST, deepseekBaseUrl, bytes.NewReader(body))
	if err != nil {
		return err
	}

	// 获取key
	apiKey := os.Getenv("DEEPSEEK_API_KEY")

	// 拼接header
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey+"11")

	// 发起请求
	resp, err := client.Do(req)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(respBody))

	var resData model.ResData
	err = json.Unmarshal(respBody, &resData)
	if err != nil {
		return err
	}

	if resData.Error != nil {
		log.Printf("ERROR INFO %+v\n", resData.Error)
		return resData.Error
	}
	return nil
}

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

	var messageArr []model.Message

	systemMessage := model.Message{Role: "system", Content: "你是我的面试助手"}

	messageArr = append(messageArr, systemMessage)

	fatalCount := 0

	for {
		fmt.Printf("你：")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			if recordFailure(&fatalCount) {
				fmt.Printf("读取用户输入失败，请重新输入：%v\n", err)
				continue
			} else {
				fmt.Printf("系统异常，已自动退出：%v\n", err)
				break
			}
		}

		userMessage := model.Message{Role: "user", Content: line}

		messageArr = append(messageArr, userMessage)

		reqBody := &model.ReqData{
			Model:   "deepseek-v4-flash",
			Message: messageArr,
		}

		body, err := json.Marshal(reqBody)
		if err != nil {
			if recordFailure(&fatalCount) {
				fmt.Printf("json序列化失败，请重新输入：%v\n", err)
				continue
			} else {
				fmt.Printf("系统异常，已自动退出：%v\n", err)
				break
			}
		}

		req, err := http.NewRequest("POST", "https://api.deepseek.com/chat/completions", bytes.NewReader(body))
		if err != nil {
			if recordFailure(&fatalCount) {
				fmt.Printf("创建请求失败，请重新输入：%v\n", err)
				continue
			} else {
				fmt.Printf("系统异常，已自动退出：%v\n", err)
				break
			}
		}

		apiKey := os.Getenv("DEEPSEEK_API_KEY")

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+apiKey+"11")

		resp, err := client.Do(req)
		if err != nil {
			if recordFailure(&fatalCount) {
				fmt.Printf("请求异常，请重新输入，异常信息：%v\n", err)
				continue
			} else {
				fmt.Printf("系统异常，已自动退出：%v\n", err)
				break
			}
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			if recordFailure(&fatalCount) {
				fmt.Printf("读取参数返回失败:%v\n", err)
				continue
			} else {
				fmt.Printf("系统异常，已自动退出：%v\n", err)
				break
			}
		}

		fmt.Println(string(respBody))

		var resData model.ResData
		err = json.Unmarshal(respBody, &resData)
		if err != nil {
			if recordFailure(&fatalCount) {
				fmt.Printf("参数解析异常:%v\n", err)
				continue
			} else {
				fmt.Printf("系统异常，已自动退出：%v\n", err)
				break
			}
		}

		if resData.Error != nil {
			fmt.Println("系统异常请稍后重试")
			log.Printf("ERROR INFO %+v\n", resData.Error)
			break
		}

		assistantOutput := resData.Choices[0].Message.Content

		fmt.Printf("面试助手：%v\n", assistantOutput)

		assistentMessage := model.Message{Role: "assistant", Content: assistantOutput}

		messageArr = append(messageArr, assistentMessage)
	}
}
