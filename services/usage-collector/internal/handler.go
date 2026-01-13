package internal

import (
    "encoding/json"
    "log"
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

    log.Printf("Received POST /v1/usage request")

    if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
        log.Printf("Error decoding request: %v", err)
        http.Error(w, "invalid payload", http.StatusBadRequest)
        return
    }

    log.Printf("Processing usage event: tenant=%s, metric=%s, value=%d", 
        event.TenantID, event.Metric, event.Value)

    if err := h.service.RecordUsage(r.Context(), event); err != nil {
        log.Printf("Error recording usage: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    log.Printf("Successfully published usage event to Kinesis")
    w.WriteHeader(http.StatusAccepted)
    w.Write([]byte("Accepted"))
}
