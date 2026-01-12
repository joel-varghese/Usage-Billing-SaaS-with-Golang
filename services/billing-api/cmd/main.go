package main

import (
    "log"
    "net/http"

    "usage-billing-platform/services/billing-api/internal"
)

func main() {
    handler := internal.NewHandler()

    mux := http.NewServeMux()
    mux.HandleFunc("/v1/billing/", handler.GetUsage)

    log.Println("Billing API listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
