package utils

type Response struct {
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

func BuildResponseSuccess(message string, data any) Response {
	res := Response{
		Message: message,
		Data:    data,
	}
	return res
}

func BuildResponseFailed(message string, err string, data any) Response {
	res := Response{
		Message: message,
		Error:   err,
		Data:    data,
	}
	return res
}
