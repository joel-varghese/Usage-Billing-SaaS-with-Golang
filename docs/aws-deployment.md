# AWS Deployment Guide

This guide will help you deploy the usage billing platform to AWS using your AWS Learner Lab instance.

## Prerequisites

1. **AWS Learner Lab Instance** with an assigned IAM role
2. **Kinesis Stream** created in your AWS account
3. **DynamoDB Table** created in your AWS account
4. **Docker** installed on your EC2 instance

## Step 1: Create AWS Resources

### Create Kinesis Stream

1. Go to AWS Console → Kinesis → Data Streams
2. Click "Create data stream"
3. Set stream name (e.g., `usage-events`)
4. Choose capacity mode (On-demand or Provisioned)
5. Note the stream name for later

### Create DynamoDB Table

1. Go to AWS Console → DynamoDB → Tables
2. Click "Create table"
3. Set table name (e.g., `usage-aggregator`)
4. Set partition key: `pk` (String)
5. Set sort key: `sk` (String)
6. Choose table settings (defaults are fine)
7. Note the table name for later

### Verify IAM Role Permissions

Your EC2 instance role should have permissions for:
- `kinesis:PutRecord`
- `kinesis:DescribeStream`
- `kinesis:GetShardIterator`
- `kinesis:GetRecords`
- `dynamodb:UpdateItem`
- `dynamodb:DescribeTable`

## Step 2: Configure Environment Variables

Create a `.env` file in the project root with the following variables:

```bash
# AWS Configuration
AWS_REGION=us-east-1  # Change to your region
AWS_ENDPOINT=         # Leave empty for real AWS (only set for LocalStack)

# AWS Credentials (NOT NEEDED for Learner Lab - will use instance role)
# AWS_ACCESS_KEY_ID=
# AWS_SECRET_ACCESS_KEY=
# AWS_SESSION_TOKEN=

# Service Configuration
KINESIS_STREAM=usage-events        # Your Kinesis stream name
DYNAMODB_TABLE=usage-aggregator    # Your DynamoDB table name
```

**Important**: 
- Leave `AWS_ENDPOINT` empty (or don't set it) to use real AWS services
- Do NOT set `AWS_ACCESS_KEY_ID` or `AWS_SECRET_ACCESS_KEY` - the SDK will automatically use your EC2 instance role credentials

## Step 3: Deploy to AWS EC2 Instance

### Option A: Using Docker Compose (Recommended)

1. **SSH into your EC2 instance**
   ```bash
   ssh -i your-key.pem ec2-user@your-ec2-ip
   ```

2. **Clone or upload your project**
   ```bash
   git clone <your-repo>  # or upload via SCP
   cd usage-billing-platform
   ```

3. **Create the .env file** with your AWS resource names
   ```bash
   nano .env
   # Add the configuration from Step 2
   ```

4. **Build and start services using docker-compose.aws.yml**
   ```bash
   docker-compose -f docker-compose.aws.yml up -d --build
   ```

5. **Check logs**
   ```bash
   docker-compose -f docker-compose.aws.yml logs -f
   ```

### Option B: Run Individual Docker Containers

1. **Build the images**
   ```bash
   docker build -f services/usage-collector/Dockerfile -t usage-collector .
   docker build -f services/aggregator/Dockerfile -t aggregator .
   docker build -f services/billing-api/Dockerfile -t billing-api .
   ```

2. **Run the containers**
   ```bash
   # Usage Collector
   docker run -d \
     --name usage-collector \
     --env-file .env \
     -p 8080:8080 \
     usage-collector

   # Aggregator
   docker run -d \
     --name aggregator \
     --env-file .env \
     aggregator

   # Billing API
   docker run -d \
     --name billing-api \
     --env-file .env \
     -p 8081:8080 \
     billing-api
   ```

## Step 4: Test the Deployment

1. **Test Usage Collector**
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

2. **Check Kinesis Stream**
   - Go to AWS Console → Kinesis → Your stream
   - Check "Monitoring" tab for incoming records

3. **Check DynamoDB Table**
   - Go to AWS Console → DynamoDB → Your table
   - Check items to see aggregated usage data

4. **Check Container Logs**
   ```bash
   docker logs aggregator
   docker logs usage-collector
   ```

## Step 5: Configure Security Groups

Make sure your EC2 security group allows:
- **Inbound**: Port 8080 (Usage Collector) and 8081 (Billing API) from your IP or 0.0.0.0/0 for testing
- **Outbound**: All traffic (for AWS API calls)

## Troubleshooting

### Issue: "Access Denied" errors
- **Solution**: Verify your EC2 instance role has the required permissions (see Step 1)

### Issue: "Stream not found" or "Table not found"
- **Solution**: Double-check your `KINESIS_STREAM` and `DYNAMODB_TABLE` environment variables match the exact names in AWS

### Issue: "Region not found"
- **Solution**: Ensure `AWS_REGION` matches the region where your resources are created

### Issue: Containers can't connect to AWS
- **Solution**: 
  1. Verify the EC2 instance has internet access
  2. Check security group outbound rules
  3. Verify IAM role is attached to the instance

### Viewing Logs
```bash
# All services
docker-compose -f docker-compose.aws.yml logs -f

# Specific service
docker-compose -f docker-compose.aws.yml logs -f aggregator
docker-compose -f docker-compose.aws.yml logs -f usage-collector
```

## Switching Between Local and AWS

- **For LocalStack (local testing)**: Use `docker-compose.yml` with `AWS_ENDPOINT=http://localstack:4566`
- **For AWS**: Use `docker-compose.aws.yml` with `AWS_ENDPOINT` empty or unset

## Next Steps

- Set up API Gateway in front of the services
- Configure CloudWatch for monitoring
- Set up auto-scaling for high throughput
- Implement proper error handling and retries
