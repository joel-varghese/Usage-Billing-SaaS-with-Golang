package internal

import (
    "context"
    "encoding/json"
    "log"
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
    endpoint := os.Getenv("AWS_ENDPOINT")
    region := os.Getenv("AWS_REGION")
    stream := os.Getenv("KINESIS_STREAM")
    
    log.Printf("Initializing Kinesis producer: endpoint=%s, region=%s, stream=%s", 
        endpoint, region, stream)
    
    awsCfg := config.LoadAWSConfig(endpoint, region)

    return &KinesisProducer{
        client:     kinesis.NewFromConfig(awsCfg),
        stream: stream,
    }
}

func (p *KinesisProducer) Publish(ctx context.Context, event models.UsageEvent) error {
    log.Printf("Publishing to Kinesis stream: %s", p.stream)
    
    data, err := json.Marshal(event)
    if err != nil {
        log.Printf("Error marshaling event: %v", err)
        return err
    }

    log.Printf("Putting record to stream: %s, partition_key: %s", p.stream, event.TenantID)
    result, err := p.client.PutRecord(ctx, &kinesis.PutRecordInput{
        StreamName:   aws.String(p.stream),
        PartitionKey: aws.String(event.TenantID),
        Data:        data,
    })

    if err != nil {
        log.Printf("Error putting record to Kinesis: %v", err)
        return err
    }

    log.Printf("Successfully put record to Kinesis. SequenceNumber: %s, ShardId: %s", 
        aws.ToString(result.SequenceNumber), aws.ToString(result.ShardId))
    return nil
}