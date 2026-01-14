package internal

import (
    "context"
    "fmt"
    "os"
    "strings"

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
    // NOTE: Fastest path fix — aggregator currently stores items with:
    // pk = TENANT#<tenantId>, sk = METRIC#<metric>, usage=<int>
    // Month is not yet part of the key, so month is ignored here.
    pk := "TENANT#" + tenantID

    out, err := r.client.Query(ctx, &dynamodb.QueryInput{
        TableName:              aws.String(r.table),
        KeyConditionExpression: aws.String("pk = :pk"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":pk": &types.AttributeValueMemberS{Value: pk},
        },
    })
    if err != nil {
        return nil, err
    }

    totals := make(map[string]int64)
    for _, item := range out.Items {
        // Expect sk = METRIC#<metric>
        sk, ok := item["sk"].(*types.AttributeValueMemberS)
        if !ok {
            continue
        }
        metric := strings.TrimPrefix(sk.Value, "METRIC#")

        usageAttr, ok := item["usage"].(*types.AttributeValueMemberN)
        if !ok {
            continue
        }
        value := usageAttr.Value

        var v int64
        fmt.Sscan(value, &v)
        totals[metric] += v
    }

    return totals, nil
}
