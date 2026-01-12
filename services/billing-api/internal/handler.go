package internal

import (
    "encoding/json"
    "net/http"
    "strings"
)

type Handler struct {
    service BillingService
}

func NewHandler() *Handler {
    return &Handler{
        service: NewService(),
    }
}

func (h *Handler) GetUsage(w http.ResponseWriter, r *http.Request) {
    // /v1/billing/{tenantId}?month=YYYY-MM
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 4 {
        http.Error(w, "invalid path", http.StatusBadRequest)
        return
    }

    tenantID := parts[3]
    month := r.URL.Query().Get("month")

    if tenantID == "" || month == "" {
        http.Error(w, "tenantId and month are required", http.StatusBadRequest)
        return
    }

    usage, err := h.service.GetMonthlyUsage(r.Context(), tenantID, month)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(usage)
}
