package internal

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"

	"usage-billing-platform/pkg/models"
)


type Consumer struct {
	kinesis *kinesis.Client
	stream  string
	repo    *UsageRepository
}

func NewConsumer(cfg aws.Config) *Consumer {
	return &Consumer{
		kinesis: kinesis.NewFromConfig(cfg),
		stream:  os.Getenv("KINESIS_STREAM"),
		repo:    NewUsageRepository(),
	}
}


func (c *Consumer) Start() {
	ctx := context.Background()

	stream, err := c.kinesis.DescribeStream(ctx, &kinesis.DescribeStreamInput{
		StreamName: aws.String(c.stream),
	})
	if err != nil {
		log.Fatal(err)
	}

	shardID := *stream.StreamDescription.Shards[0].ShardId
	log.Printf("Using shard: %s\n", shardID)

	iterator, err := c.kinesis.GetShardIterator(ctx, &kinesis.GetShardIteratorInput{
		StreamName:        aws.String(c.stream),
		ShardId:           aws.String(shardID),
		ShardIteratorType: "LATEST",
	})
	if err != nil {
		log.Fatal(err)
	}

	shardIterator := iterator.ShardIterator
	log.Println("Consumer started, listening for new events (LATEST iterator)...")

	for {
		out, err := c.kinesis.GetRecords(ctx, &kinesis.GetRecordsInput{
			ShardIterator: shardIterator,
			Limit:         aws.Int32(100),
		})
		if err != nil {
			log.Println("get records error:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if len(out.Records) > 0 {
			log.Printf("Received %d record(s) from Kinesis\n", len(out.Records))
		}

		for _, r := range out.Records {
			var event models.UsageEvent
			if err := json.Unmarshal(r.Data, &event); err != nil {
				log.Println("invalid record:", err)
				continue
			}

			pk := "TENANT#" + event.TenantID
			sk := "METRIC#" + event.Metric

			log.Printf("Processing event: tenant=%s, metric=%s, value=%d\n", event.TenantID, event.Metric, event.Value)
			if err := c.repo.IncrementUsage(pk, sk, event.Value); err != nil {
				log.Println("dynamodb error:", err)
			} else {
				log.Printf("Successfully updated DynamoDB: pk=%s, sk=%s\n", pk, sk)
			}
		}

		shardIterator = out.NextShardIterator
		time.Sleep(1 * time.Second)
	}
}

