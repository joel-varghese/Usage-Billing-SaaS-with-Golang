package internal

import (
    "context"
    "fmt"
    "os"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

    "usage-billing-platform/pkg/config"
)

type UsageRepository interface {
    GetMonthlyUsage(ctx context.Context, tenantID, month string) (map[string]int64, error)
}

type DynamoRepository struct {
    client *dynamodb.Client
    table  string
}

func NewDynamoRepository() *DynamoRepository {
    cfg := config.LoadAWSConfig(
        os.Getenv("AWS_ENDPOINT"),
        os.Getenv("AWS_REGION"),
    )

    return &DynamoRepository{
        client: dynamodb.NewFromConfig(cfg),
        table:  os.Getenv("DYNAMODB_TABLE"),
    }
}

func (r *DynamoRepository) GetMonthlyUsage(
    ctx context.Context,
    tenantID string,
    month string,
) (map[string]int64, error) {

    out, err := r.client.Query(ctx, &dynamodb.QueryInput{
        TableName: aws.String(r.table),
        KeyConditionExpression: aws.String(
            "tenant_id = :tenant AND usage_month = :month",
        ),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":tenant": &types.AttributeValueMemberS{Value: tenantID},
            ":month":  &types.AttributeValueMemberS{Value: month},
        },
    })
    if err != nil {
        return nil, err
    }

    totals := make(map[string]int64)
    for _, item := range out.Items {
        metric := item["metric"].(*types.AttributeValueMemberS).Value
        value := item["total"].(*types.AttributeValueMemberN).Value

        var v int64
        fmt.Sscan(value, &v)
        totals[metric] += v
    }

    return totals, nil
}
