package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/oapi-codegen/runtime/types"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application"
	codegen "github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/infrastructure/http/codegen"
)

func mapExchangeResponse(dto application.ExchangeDTO) codegen.ExchangeRequest {
	status := codegen.ExchangeRequestStatus(dto.Status)
	return codegen.ExchangeRequest{
		Id:          &dto.ID,
		OwnedBookId: &dto.OwnedBookID,
		FromId:      &dto.FromID,
		ToId:        &dto.ToID,
		Status:      &status,
		CreatedAt:   &dto.CreatedAt,
		ExpiresAt:   &dto.ExpiresAt,
		Note:        &dto.Note,
	}
}

func (h *Handler) CreateExchange(w http.ResponseWriter, r *http.Request) {
	var body codegen.CreateExchangeJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, application.ErrInvalidParams)
		return
	}

	expiresAt := time.Time{}
	if body.ExpiresAt != nil {
		expiresAt = *body.ExpiresAt
	}

	note := ""
	if body.Note != nil {
		note = *body.Note
	}

	exchange, err := h.exchangeService.CreateExchange(r.Context(), body.OwnedBookId, body.ToUserId, expiresAt, note)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, mapExchangeResponse(exchange))
}

func (h *Handler) GetExchange(w http.ResponseWriter, r *http.Request, exchangeId types.UUID) {
	exchange, err := h.exchangeService.GetExchangeByID(r.Context(), exchangeId)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, mapExchangeResponse(exchange))
}

func (h *Handler) AcceptExchange(w http.ResponseWriter, r *http.Request, exchangeId types.UUID) {
	exchange, err := h.exchangeService.AcceptExchange(r.Context(), exchangeId)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, mapExchangeResponse(exchange))
}

func (h *Handler) RejectExchange(w http.ResponseWriter, r *http.Request, exchangeId types.UUID) {
	exchange, err := h.exchangeService.RejectExchange(r.Context(), exchangeId)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, mapExchangeResponse(exchange))
}

func (h *Handler) CancelExchange(w http.ResponseWriter, r *http.Request, exchangeId types.UUID) {
	exchange, err := h.exchangeService.CancelExchange(r.Context(), exchangeId)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, mapExchangeResponse(exchange))
}

func (h *Handler) GetMyExchanges(w http.ResponseWriter, r *http.Request, params codegen.GetMyExchangesParams) {
	userID, ok := GetUserID(r.Context())
	if !ok {
		writeError(w, application.ErrInvalidCredentials)
		return
	}

	status := ""
	if params.Status != nil {
		status = string(*params.Status)
	}

	exchanges, err := h.exchangeService.GetUserExchanges(r.Context(), userID, status)
	if err != nil {
		writeError(w, err)
		return
	}

	resp := make([]codegen.ExchangeRequest, 0, len(exchanges))
	for _, e := range exchanges {
		resp = append(resp, mapExchangeResponse(e))
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"exchanges": resp,
	})
}
