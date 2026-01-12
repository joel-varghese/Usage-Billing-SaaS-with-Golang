package models

import "time"

type UsageEvent struct {
	TenantID  string    `json:"tenant_id"`
	Metric    string    `json:"metric"`
	Value     int64     `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}