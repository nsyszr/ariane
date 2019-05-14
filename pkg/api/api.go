package api

import "fmt"

type Header map[string][]string

type Request struct {
	Method string      `json:"method"`
	Header Header      `json:"header,omitempty"`
	Object interface{} `json:"object"`
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Object     interface{} `json:"object,omitempty"`
	Error      *Error      `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}
