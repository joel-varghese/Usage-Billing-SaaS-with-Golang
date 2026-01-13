# AWS Deployment Quick Start

## Quick Setup Steps

### 1. Create AWS Resources

**Kinesis Stream:**
- Name: `usage-events` (or your preferred name)
- Capacity: On-demand or Provisioned

**DynamoDB Table:**
- Name: `usage-aggregator` (or your preferred name)
- Partition Key: `pk` (String)
- Sort Key: `sk` (String)

### 2. Create .env File

On your EC2 instance, create a `.env` file:

```bash
AWS_REGION=us-east-1
AWS_ENDPOINT=
KINESIS_STREAM=usage-events
DYNAMODB_TABLE=usage-aggregator
```

**Important:** 
- Leave `AWS_ENDPOINT` empty (or don't set it)
- Do NOT set `AWS_ACCESS_KEY_ID` or `AWS_SECRET_ACCESS_KEY` - the instance role will be used automatically

### 3. Deploy

```bash
# Build and start all services
docker-compose -f docker-compose.aws.yml up -d --build

# View logs
docker-compose -f docker-compose.aws.yml logs -f
```

### 4. Test

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

### 5. Verify

- Check Kinesis stream in AWS Console for incoming records
- Check DynamoDB table for aggregated data (pk: `TENANT#tenant-123`, sk: `METRIC#api_calls`)

## IAM Role Permissions Required

Your EC2 instance role needs:
- `kinesis:PutRecord`
- `kinesis:DescribeStream`
- `kinesis:GetShardIterator`
- `kinesis:GetRecords`
- `dynamodb:UpdateItem`
- `dynamodb:DescribeTable`

## Troubleshooting

**Access Denied?** → Check IAM role permissions
**Stream/Table not found?** → Verify names in `.env` match AWS exactly
**Can't connect?** → Check security group outbound rules

For detailed information, see [aws-deployment.md](./aws-deployment.md)
