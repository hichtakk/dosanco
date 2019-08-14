package handler

type Error struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

func returnBusinessError(msg string) ErrorResponse {
	return ErrorResponse{Error: Error{Message: msg}}
}