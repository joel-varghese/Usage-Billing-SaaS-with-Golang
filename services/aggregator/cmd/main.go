package main

import (
    "os"
    "log"
    "usage-billing-platform/services/aggregator/internal"
    "usage-billing-platform/pkg/config"

)

func main() {
	cfg := config.LoadAWSConfig(os.Getenv("AWS_ENDPOINT"), os.Getenv("AWS_REGION"))
	consumer := internal.NewConsumer(cfg)
    log.Println("Starting aggregator consumer...")
	consumer.Start()
}
