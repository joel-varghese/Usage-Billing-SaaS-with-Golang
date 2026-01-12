package internal

import (
    "time"

    "usage-billing-platform/pkg/models"
)

func monthKey(t time.Time) string {
    return t.Format("2006-01")
}

func (c *Consumer) aggregate(event models.UsageEvent) error {
    pk := event.TenantID + "#" + monthKey(event.Timestamp)
    sk := event.Metric

    return c.repo.IncrementUsage(pk, sk, event.Value)
}
