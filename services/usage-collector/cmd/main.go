package main

import (
    "log"
    "net/http"

    "usage-billing-platform/services/usage-collector/internal"
)

func main() {
    handler := internal.NewHandler()

    mux := http.NewServeMux()
    mux.HandleFunc("/v1/usage", handler.PostUsage)

    log.Println("Usage Collector listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
