package utils

// Response is the standard API envelope, mirroring FastAPI's APIResponse:
// {message, data, success, metas}. Error paths add error_code/details.
type Response struct {
	Message   string `json:"message"`
	Data      any    `json:"data"`
	Success   bool   `json:"success"`
	Metas     any    `json:"metas"`
	ErrorCode string `json:"error_code,omitempty"`
	Details   any    `json:"details,omitempty"`
}

func BuildResponseSuccess(message string, data any, meta *any) Response {
	res := Response{
		Message: message,
		Data:    data,
		Success: true,
	}
	if meta != nil {
		res.Metas = meta
	}
	return res
}

func BuildResponseFailed(message string, errCode string, data any) Response {
	return Response{
		Message:   message,
		Data:      data,
		Success:   false,
		ErrorCode: errCode,
	}
}
