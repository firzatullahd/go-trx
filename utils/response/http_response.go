package response

import "time"

type Response struct {
	AccessTime time.Time   `json:"accessTime"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Result     interface{} `json:"result"`
}

func SetResponse(statusCode int, message string, Result interface{}) Response {
	return Response{
		AccessTime: time.Now(),
		StatusCode: statusCode,
		Message:    message,
		Result:     Result,
	}
}
