package handler

// Error is struct for representing error response message.
type Error struct {
	Message string `json:"message"`
}

// ErrorResponse wraps Error for response.
type ErrorResponse struct {
	Error   Error  `json:"error"`
	Message string `json:"message"`
}

func returnError(msg string) *ErrorResponse {
	return &ErrorResponse{Error: Error{Message: msg}}
}

// ResponseMessage is struct for representing error response message.
type ResponseMessage struct {
	Message string `json:"message"`
}

func returnMessage(msg string) *ResponseMessage {
	return &ResponseMessage{Message: msg}
}
