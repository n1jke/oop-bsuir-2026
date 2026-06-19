package server

import (
	"encoding/json"
	"net/http"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application"
	codegen "github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/infrastructure/http/codegen"
)

func (h *Handler) CreateBookReview(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		writeError(w, application.ErrInvalidCredentials)
		return
	}

	var body codegen.CreateBookReviewJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, application.ErrInvalidParams)
		return
	}

	report := ""
	if body.Report != nil {
		report = *body.Report
	}

	review, err := h.reviewService.AddBookReview(r.Context(), userID, body.BookId, uint(body.Mark), report)
	if err != nil {
		writeError(w, err)
		return
	}

	mark := int(review.Mark)
	writeJSON(w, http.StatusCreated, codegen.Review{
		Id:     &review.ID,
		FromId: &review.FromID,
		BookId: &review.BookID,
		Mark:   &mark,
		Report: &review.Report,
	})
}

func (h *Handler) CreateUserReview(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		writeError(w, application.ErrInvalidCredentials)
		return
	}

	var body codegen.CreateUserReviewJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, application.ErrInvalidParams)
		return
	}

	report := ""
	if body.Report != nil {
		report = *body.Report
	}

	review, err := h.reviewService.AddUserReview(r.Context(), userID, body.UserId, uint(body.Mark), report)
	if err != nil {
		writeError(w, err)
		return
	}

	mark := int(review.Mark)
	writeJSON(w, http.StatusCreated, codegen.Review{
		Id:     &review.ID,
		FromId: &review.FromID,
		ToId:   &review.ToID,
		Mark:   &mark,
		Report: &review.Report,
	})
}
