package internal

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
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

func (c *Consumer) processShard(ctx context.Context, shardID string) {
	log.Printf("Starting consumer for shard: %s", shardID)

	iterator, err := c.kinesis.GetShardIterator(ctx, &kinesis.GetShardIteratorInput{
		StreamName:        aws.String(c.stream),
		ShardId:           aws.String(shardID),
		ShardIteratorType: "LATEST",
	})
	if err != nil {
		log.Printf("Error getting shard iterator for %s: %v", shardID, err)
		return
	}

	shardIterator := iterator.ShardIterator
	log.Printf("Consumer started for shard %s, listening for new events (LATEST iterator)...", shardID)

	for {
		out, err := c.kinesis.GetRecords(ctx, &kinesis.GetRecordsInput{
			ShardIterator: shardIterator,
			Limit:         aws.Int32(100),
		})
		if err != nil {
			log.Printf("get records error for shard %s: %v", shardID, err)
			time.Sleep(2 * time.Second)
			continue
		}

		if len(out.Records) > 0 {
			log.Printf("Shard %s: Received %d record(s) from Kinesis", shardID, len(out.Records))
		}

		for _, r := range out.Records {
			var event models.UsageEvent
			if err := json.Unmarshal(r.Data, &event); err != nil {
				log.Printf("Shard %s: invalid record: %v", shardID, err)
				continue
			}

			pk := "TENANT#" + event.TenantID
			sk := "METRIC#" + event.Metric

			log.Printf("Shard %s: Processing event: tenant=%s, metric=%s, value=%d", 
				shardID, event.TenantID, event.Metric, event.Value)
			if err := c.repo.IncrementUsage(pk, sk, event.Value); err != nil {
				log.Printf("Shard %s: dynamodb error: %v", shardID, err)
			} else {
				log.Printf("Shard %s: Successfully updated DynamoDB: pk=%s, sk=%s", shardID, pk, sk)
			}
		}

		shardIterator = out.NextShardIterator
		if shardIterator == nil {
			log.Printf("Shard %s: NextShardIterator is nil, exiting consumer", shardID)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func (c *Consumer) Start() {
	ctx := context.Background()

	log.Printf("Describing stream: %s", c.stream)
	stream, err := c.kinesis.DescribeStream(ctx, &kinesis.DescribeStreamInput{
		StreamName: aws.String(c.stream),
	})
	if err != nil {
		log.Fatalf("Error describing stream: %v", err)
	}

	shards := stream.StreamDescription.Shards
	log.Printf("Found %d shard(s) in stream", len(shards))

	if len(shards) == 0 {
		log.Fatal("No shards found in stream")
	}

	var wg sync.WaitGroup
	for _, shard := range shards {
		if shard.ShardId == nil {
			continue
		}
		shardID := *shard.ShardId
		log.Printf("Starting consumer for shard: %s", shardID)
		wg.Add(1)
		go func(sid string) {
			defer wg.Done()
			c.processShard(ctx, sid)
		}(shardID)
	}

	log.Println("All shard consumers started. Waiting for events...")
	wg.Wait()
}
