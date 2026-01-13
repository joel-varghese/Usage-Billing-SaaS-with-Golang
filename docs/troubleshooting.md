# Troubleshooting Guide

## Issue: DynamoDB Table Shows No Records

### Step 1: Check if Aggregator Container is Running

```bash
docker compose -f docker-compose.aws.yml ps
```

All three containers (usage-collector, aggregator, billing-api) should be in "Up" status.

### Step 2: Check Aggregator Logs

```bash
docker compose -f docker-compose.aws.yml logs aggregator
```

Or follow logs in real-time:
```bash
docker compose -f docker-compose.aws.yml logs -f aggregator
```

Look for:
- ✅ "Consumer started, listening for new events"
- ✅ "Received X record(s) from Kinesis"
- ✅ "Successfully updated DynamoDB"
- ❌ Any error messages

### Step 3: Verify Environment Variables in Aggregator

```bash
docker compose -f docker-compose.aws.yml exec aggregator env | grep -E "AWS_|KINESIS_|DYNAMODB_"
```

Should show:
- `AWS_REGION=us-east-1`
- `AWS_ENDPOINT=` (empty)
- `KINESIS_STREAM=usage-events`
- `DYNAMODB_TABLE=usage-aggregator`

### Step 4: Test the Full Flow

**Important**: The aggregator uses `LATEST` iterator, which means it only processes records sent AFTER it starts. If you sent events before starting the aggregator, they won't be processed.

1. **Make sure aggregator is running:**
   ```bash
   docker compose -f docker-compose.aws.yml logs aggregator | tail -20
   ```

2. **Send a test event:**
   ```bash
   curl -X POST http://localhost:8080/v1/usage \
     -H "Content-Type: application/json" \
     -d '{
       "tenant_id": "tenant-123",
       "metric": "api_calls",
       "value": 10,
       "timestamp": "2024-01-15T10:00:00Z"
     }'
   ```

3. **Watch aggregator logs immediately:**
   ```bash
   docker compose -f docker-compose.aws.yml logs -f aggregator
   ```

   You should see:
   - "Received 1 record(s) from Kinesis"
   - "Processing event: tenant=tenant-123, metric=api_calls, value=10"
   - "Successfully updated DynamoDB: pk=TENANT#tenant-123, sk=METRIC#api_calls"

4. **Check DynamoDB:**
   ```bash
   aws dynamodb get-item \
     --table-name usage-aggregator \
     --key '{"pk": {"S": "TENANT#tenant-123"}, "sk": {"S": "METRIC#api_calls"}}' \
     --region us-east-1
   ```

### Step 5: Check IAM Permissions

Verify your EC2 instance role has these permissions:

```bash
# Check if you can describe the stream
aws kinesis describe-stream --stream-name usage-events --region us-east-1

# Check if you can describe the table
aws dynamodb describe-table --table-name usage-aggregator --region us-east-1

# Try to update an item (this will test write permissions)
aws dynamodb update-item \
  --table-name usage-aggregator \
  --key '{"pk": {"S": "TENANT#test"}, "sk": {"S": "METRIC#test"}}' \
  --update-expression "ADD #usage :inc SET last_updated = :now" \
  --expression-attribute-names '{"#usage": "usage"}' \
  --expression-attribute-values '{":inc": {"N": "1"}, ":now": {"S": "2024-01-15T10:00:00Z"}}' \
  --region us-east-1
```

### Step 6: Common Issues

#### Issue: "Access Denied" errors in logs
**Solution**: Your EC2 instance role needs:
- `kinesis:DescribeStream`
- `kinesis:GetShardIterator`
- `kinesis:GetRecords`
- `dynamodb:UpdateItem`
- `dynamodb:DescribeTable`

#### Issue: "Stream not found" or "Table not found"
**Solution**: Verify exact names match:
- `KINESIS_STREAM` in .env matches stream name in AWS
- `DYNAMODB_TABLE` in .env matches table name in AWS

#### Issue: Aggregator logs show "Consumer started" but no records
**Solution**: 
- The aggregator uses `LATEST` iterator - it only processes NEW records
- Send a test event AFTER the aggregator is running
- Check usage-collector logs to ensure events are being sent to Kinesis

#### Issue: Records received but DynamoDB update fails
**Solution**: Check DynamoDB permissions and table structure (must have `pk` and `sk` keys)

### Step 7: Verify Usage Collector is Working

```bash
# Check usage-collector logs
docker compose -f docker-compose.aws.yml logs usage-collector

# Send a test event and watch collector logs
docker compose -f docker-compose.aws.yml logs -f usage-collector
# Then in another terminal:
curl -X POST http://localhost:8080/v1/usage \
  -H "Content-Type: application/json" \
  -d '{"tenant_id": "test", "metric": "test", "value": 1, "timestamp": "2024-01-15T10:00:00Z"}'
```

You should see the collector successfully publishing to Kinesis.
