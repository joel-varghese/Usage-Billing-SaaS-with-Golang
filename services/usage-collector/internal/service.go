package internal

import (
    "context"
    "errors"
    "time"

    "usage-billing-platform/pkg/models"
)

type UsageService interface {
    RecordUsage(ctx context.Context, event models.UsageEvent) error
}

type service struct {
    producer Producer
}

func NewService() UsageService {
    return &service{
        producer: NewKinesisProducer(),
    }
}

func (s *service) RecordUsage(
    ctx context.Context,
    event models.UsageEvent,
) error {

    // Basic validation
    if event.TenantID == "" {
        return errors.New("tenant_id is required")
    }
    if event.Metric == "" {
        return errors.New("metric is required")
    }
    if event.Value <= 0 {
        return errors.New("value must be positive")
    }

    // Enrichment
    if event.Timestamp.IsZero() {
        event.Timestamp = time.Now().UTC()
    }

    // Publish to Kinesis
    return s.producer.Publish(ctx, event)
}
