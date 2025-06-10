package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func handleWritingErr(err error) {
	if err != nil {
		slog.Error(fmt.Sprintf("Error writing to http.ResponseWriter: %v", err))
	}
}

type errorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

func HandleError(w http.ResponseWriter, handlerErr error, statusCode int, shouldLog bool) {
	if shouldLog {
		slog.Error(fmt.Sprintf("Handler error: %v", handlerErr), "status", statusCode)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Hide server errors from the client
	if statusCode >= 500 {
		handlerErr = fmt.Errorf("%s: %d", http.StatusText(statusCode), statusCode)
	}

	err := json.NewEncoder(w).Encode(errorResponse{
		Error:      handlerErr.Error(),
		StatusCode: statusCode,
	})
	handleWritingErr(err)
}
