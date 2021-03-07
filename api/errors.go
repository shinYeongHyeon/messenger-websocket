package api

import (
	"encoding/json"
	"net/http"
)

type apiErrorWriter interface {
	Write(w http.ResponseWriter)
}

func handlePanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				if e, ok := r.(apiErrorWriter); ok {
					e.Write(w)
				} else {
					internalError.Write(w)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type errorResponse struct {
	Error string `json:"error"`
}

type simpleError struct {
	message string
	status  int
}

func (e simpleError) Write(w http.ResponseWriter) {
	w.WriteHeader(e.status)
	json.NewEncoder(w).Encode(errorResponse{e.message})
}

var internalError = simpleError{
	message: "Internal Error",
	status:  http.StatusInternalServerError,
}

var malformedInputError = simpleError{
	message: "Malformed Input Error",
	status:  http.StatusBadRequest,
}

var notFoundError = simpleError{
	message: "Not Found",
	status:  http.StatusNotFound,
}

var unauthorizedError = simpleError{
	message: "Unauthorized",
	status:  http.StatusUnauthorized,
}
