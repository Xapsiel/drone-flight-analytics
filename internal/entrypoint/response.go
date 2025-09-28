package httpv1

import "time"

type APIResponse struct {
	Data      interface{} `json:"data"`              // Данные ответа
	Error     *APIError   `json:"error,omitempty"`   // Ошибка, если есть
	Status    string      `json:"status"`            // Статус (success, error)
	Timestamp time.Time   `json:"timestamp"`         // Время ответа
	Message   string      `json:"message,omitempty"` // Дополнительное сообщение
}

type APIError struct {
	Code    int    `json:"code"`    // Код ошибки (например, 400, 500)
	Message string `json:"message"` // Текст ошибки
}

func (r *Router) NewSuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Data:      data,
		Error:     nil,
		Status:    "success",
		Timestamp: time.Now(),
		Message:   message,
	}
}

func (r *Router) NewErrorResponse(code int, message string) APIResponse {
	return APIResponse{
		Data:      nil,
		Error:     &APIError{Code: code, Message: message},
		Status:    "error",
		Timestamp: time.Now(),
		Message:   message,
	}
}
