package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/zhq0808/interviewAgent/model"
)

type Client struct {
	Client  *http.Client
	BaseUrl string
	APIKey  string
	Model   string
}

func NewClient(client *http.Client, baseUrl, apiKey, model string) *Client {
	return &Client{
		BaseUrl: baseUrl,
		APIKey:  apiKey,
		Client:  client,
		Model:   model,
	}
}

func (c *Client) SendRequest(messageArr []model.Message) (*model.ResData, error) {
	reqBody := &model.ReqData{
		Model:   c.Model,
		Message: messageArr,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	// 创建请求
	req, err := http.NewRequest(HTTP_POST, c.BaseUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// 拼接header
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.APIKey+"11")

	// 发起请求
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(respBody))

	var resData model.ResData
	err = json.Unmarshal(respBody, &resData)
	if err != nil {
		return nil, err
	}

	if resData.Error != nil {
		log.Printf("ERROR INFO %+v\n", resData.Error)
		return nil, resData.Error
	}
	return &resData, nil
}
