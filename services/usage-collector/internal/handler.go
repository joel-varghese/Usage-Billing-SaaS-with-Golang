package internal

import (
    "encoding/json"
    "net/http"

    "usage-billing-platform/pkg/models"
)

type Handler struct {
    service UsageService
}

func NewHandler() *Handler {
    return &Handler{
        service: NewService(),
    }
}

func (h *Handler) PostUsage(w http.ResponseWriter, r *http.Request) {
    var event models.UsageEvent

    if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
        http.Error(w, "invalid payload", http.StatusBadRequest)
        return
    }

    if err := h.service.RecordUsage(r.Context(), event); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusAccepted)
}
