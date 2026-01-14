#!/bin/bash

echo "=== AWS Deployment Debug Script ==="
echo ""

echo "1. Checking Docker containers status..."
docker compose -f docker-compose.aws.yml ps
echo ""

echo "2. Checking usage-collector environment variables..."
docker compose -f docker-compose.aws.yml exec usage-collector env | grep -E "AWS_|KINESIS_|DYNAMODB_" || echo "Container not running or exec failed"
echo ""

echo "3. Checking aggregator environment variables..."
docker compose -f docker-compose.aws.yml exec aggregator env | grep -E "AWS_|KINESIS_|DYNAMODB_" || echo "Container not running or exec failed"
echo ""

echo "4. Recent usage-collector logs (last 20 lines)..."
docker compose -f docker-compose.aws.yml logs --tail=20 usage-collector
echo ""

echo "5. Recent aggregator logs (last 20 lines)..."
docker compose -f docker-compose.aws.yml logs --tail=20 aggregator
echo ""

echo "6. Testing if usage-collector is listening on port 8080..."
curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" http://localhost:8080/v1/usage || echo "Connection failed"
echo ""

echo "7. Checking Kinesis stream records (last 5 minutes)..."
echo "Note: This requires AWS CLI and may show empty if no records were sent"
aws kinesis get-records \
  --shard-iterator $(aws kinesis get-shard-iterator \
    --stream-name usage-events \
    --shard-id shardId-000000000000 \
    --shard-iterator-type TRIM_HORIZON \
    --region us-east-1 \
    --query 'ShardIterator' \
    --output text) \
  --region us-east-1 \
  --max-items 5 2>/dev/null | jq -r '.Records | length' && echo " records found" || echo "No records or error"
echo ""

echo "=== Debug Complete ==="
echo ""
echo "Next steps:"
echo "1. If containers are not running: docker compose -f docker-compose.aws.yml up -d"
echo "2. Watch logs in real-time: docker compose -f docker-compose.aws.yml logs -f"
echo "3. Send a test event: curl -X POST http://localhost:8080/v1/usage -H 'Content-Type: application/json' -d '{\"tenant_id\":\"test\",\"metric\":\"test\",\"value\":1,\"timestamp\":\"2024-01-15T10:00:00Z\"}'"
