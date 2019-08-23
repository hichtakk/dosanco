package handler

// Error is struct for representing error response message.
type Error struct {
	Message string `json:"message"`
}

// ErrorResponse wraps Error for response.
type ErrorResponse struct {
	Error Error `json:"error"`
}

func returnBusinessError(msg string) ErrorResponse {
	return ErrorResponse{Error: Error{Message: msg}}
}
