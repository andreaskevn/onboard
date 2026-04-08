package dto

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}
