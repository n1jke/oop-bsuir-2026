package server

import (
	"encoding/json"
	"net/http"

	"github.com/oapi-codegen/runtime/types"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application"
	codegen "github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/infrastructure/http/codegen"
)

func mapBookResponse(dto application.BookDTO) codegen.Book {
	return codegen.Book{
		Id:          &dto.ID,
		Title:       &dto.Title,
		AuthorName:  &dto.AuthorName,
		ISBN:        &dto.ISBN,
		Description: &dto.Description,
		Topics:      &dto.Topics,
	}
}

func (h *Handler) GetBooks(w http.ResponseWriter, r *http.Request, params codegen.GetBooksParams) {
	title := ""
	if params.Title != nil {
		title = *params.Title
	}

	topic := ""
	if params.Topic != nil {
		topic = *params.Topic
	}

	resp, err := h.bookService.SearchBook(r.Context(), title, topic)
	if err != nil {
		writeError(w, err)
		return
	}

	books := make([]codegen.Book, 0, len(resp.Books))
	for _, b := range resp.Books {
		books = append(books, mapBookResponse(b))
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"books": books,
		"total": resp.Total,
	})
}

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var body codegen.CreateBookJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, application.ErrInvalidParams)
		return
	}

	topic := ""
	if body.Topics != nil && len(*body.Topics) > 0 {
		topic = (*body.Topics)[0]
	}

	description := ""
	if body.Description != nil {
		description = *body.Description
	}

	book, err := h.bookService.CreateBook(r.Context(), body.Title, "", body.Isbn, description, topic)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, mapBookResponse(book))
}

func (h *Handler) GetBookByID(w http.ResponseWriter, r *http.Request, bookId types.UUID) {
	book, err := h.bookService.GetBookByID(r.Context(), bookId)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, mapBookResponse(book))
}
