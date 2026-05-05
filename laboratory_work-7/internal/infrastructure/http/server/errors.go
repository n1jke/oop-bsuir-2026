package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application"
	codegen "github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/infrastructure/http/codegen"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	status := errorStatus(err)
	msg := err.Error()

	writeJSON(w, status, codegen.Error{
		Error:   &msg,
		Message: &msg,
	})
}

func errorStatus(err error) int {
	switch {
	case errors.Is(err, application.ErrInvalidCredentials):
		return http.StatusUnauthorized
	case errors.Is(err, application.ErrUserNotFound):
		return http.StatusNotFound
	case errors.Is(err, application.ErrBookNotFound):
		return http.StatusNotFound
	case errors.Is(err, application.ErrUserAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, application.ErrAlreadyExist):
		return http.StatusConflict
	case errors.Is(err, application.ErrInvalidParams):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
