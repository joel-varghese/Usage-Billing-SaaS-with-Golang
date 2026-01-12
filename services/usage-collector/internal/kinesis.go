package internal

import (
    "context"
    "encoding/json"
    "os"

    "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
    "usage-billing-platform/pkg/config"
    "usage-billing-platform/pkg/models"
)

type Producer interface {
	Publish(ctx context.Context, event models.UsageEvent) error
}

type KinesisProducer struct {
    client *kinesis.Client
    stream string
}

func NewKinesisProducer() Producer {
    awsCfg := config.LoadAWSConfig(
        os.Getenv("AWS_ENDPOINT"),
        os.Getenv("AWS_REGION"),
    )

    return &KinesisProducer{
        client:     kinesis.NewFromConfig(awsCfg),
        stream: os.Getenv("KINESIS_STREAM"),
    }
}

func (p *KinesisProducer) Publish(ctx context.Context, event models.UsageEvent) error {
    data, err := json.Marshal(event)

    if err != nil {
        return err
    }

    _, err = p.client.PutRecord(ctx, &kinesis.PutRecordInput{
        StreamName:   aws.String(p.stream),
        PartitionKey: aws.String(event.TenantID),
        Data:        data,
    })

    return err
}