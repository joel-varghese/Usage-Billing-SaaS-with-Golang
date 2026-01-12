package internal

import (
    "context"
    "os"
    "time"
    "strconv"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

    "usage-billing-platform/pkg/config"
)

type UsageRepository struct {
    client *dynamodb.Client
    table  string
}
func NewUsageRepository() *UsageRepository {
    cfg := config.LoadAWSConfig(
        os.Getenv("AWS_ENDPOINT"),
        os.Getenv("AWS_REGION"),
    )

    return &UsageRepository{
        client: dynamodb.NewFromConfig(cfg),
        table:  os.Getenv("DYNAMODB_TABLE"),
    }
}

func (r *UsageRepository) IncrementUsage(
    pk string,
    sk string,
    quantity int64,
) error {

    _, err := r.client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
        TableName: aws.String(r.table),
        Key: map[string]types.AttributeValue{
            "pk": &types.AttributeValueMemberS{Value: pk},
            "sk": &types.AttributeValueMemberS{Value: sk},
        },
        UpdateExpression: aws.String(
            "ADD #usage :inc SET last_updated = :now",
        ),
        ExpressionAttributeNames: map[string]string{
            "#usage": "usage",
        },
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":inc": &types.AttributeValueMemberN{
                Value: strconv.FormatInt(quantity, 10),
            },
            ":now": &types.AttributeValueMemberS{
                Value: time.Now().UTC().Format(time.RFC3339),
            },
        },
    })

    return err
}
