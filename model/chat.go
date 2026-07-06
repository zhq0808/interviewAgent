package model

import "fmt"

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
	Error   *Error   `json:"error"`
}

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

type Choice struct {
	Message Message `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}
