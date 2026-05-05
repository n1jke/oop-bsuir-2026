package server

import (
	"encoding/json"
	"net/http"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application"
	codegen "github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/infrastructure/http/codegen"
)

func (h *Handler) GetMyLibrary(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		writeError(w, application.ErrInvalidCredentials)
		return
	}

	books, err := h.libraryService.GetUserBooks(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	resp := make([]codegen.Book, 0, len(books))
	for _, b := range books {
		resp = append(resp, mapBookResponse(&b))
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"books": resp,
	})
}

func (h *Handler) AddToLibrary(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		writeError(w, application.ErrInvalidCredentials)
		return
	}

	var body codegen.AddToLibraryJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, application.ErrInvalidParams)
		return
	}

	err := h.libraryService.AddBook(r.Context(), userID, body.BookId)
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
