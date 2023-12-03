package response

import "time"

type Response struct {
	AccessTime time.Time   `json:"accessTime,omitempty"`
	StatusCode int         `json:"statusCode,omitempty"`
	Message    string      `json:"message,omitempty"`
	Result     interface{} `json:"result,omitempty"`
}

func SetResponse(statusCode int, message string, Result interface{}) Response {
	return Response{
		AccessTime: time.Now(),
		StatusCode: statusCode,
		Message:    message,
		Result:     Result,
	}
}
