package server

import (
	"encoding/json"
	"net/http"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application"
	codegen "github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/infrastructure/http/codegen"
)

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var body codegen.LoginUserJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, application.ErrInvalidParams)
		return
	}

	resp, err := h.authService.Login(r.Context(), body.Username, body.Password)
	if err != nil {
		writeError(w, err)
		return
	}

	rating := float32(resp.User.Rating) // todo: move to float64
	writeJSON(w, http.StatusOK, map[string]any{
		"token": resp.Token,
		"user": codegen.User{
			Id:     &resp.User.ID,
			Name:   &resp.User.Name,
			Rating: &rating,
		},
	})
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var body codegen.RegisterUserJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, application.ErrInvalidParams)
		return
	}

	resp, err := h.authService.Register(r.Context(), body.Username, body.Password)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":    resp.ID,
		"token": resp.Token,
	})
}
