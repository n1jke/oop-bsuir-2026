package server

import (
	"net/http"

	codegen "github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/infrastructure/http/codegen"
)

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAll(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}

	resp := make([]codegen.User, 0, len(users))
	for _, u := range users {
		rating := float32(u.Rating)
		resp = append(resp, codegen.User{
			Id:     &u.ID,
			Name:   &u.Name,
			Rating: &rating,
		})
	}

	writeJSON(w, http.StatusOK, resp)
}
