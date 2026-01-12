package internal

import "context"

type BillingService interface {
    GetMonthlyUsage(ctx context.Context, tenantID, month string) (map[string]int64, error)
}

type service struct {
    repo UsageRepository
}

func NewService() BillingService {
    return &service{
        repo: NewDynamoRepository(),
    }
}

func (s *service) GetMonthlyUsage(
    ctx context.Context,
    tenantID string,
    month string,
) (map[string]int64, error) {

    // Business rules go here later
    return s.repo.GetMonthlyUsage(ctx, tenantID, month)
}
