# Usage-Based Billing & Metering Platform

A multi-tenant, cloud-native usage metering and billing system built with Golang and AWS.

## Features
- High-throughput usage ingestion
- Fault-tolerant aggregation
- Monthly billing summaries
- Multi-tenant isolation

## Stack
- Golang
- AWS (API Gateway, Kinesis, DynamoDB, ECS)
- Terraform

## Getting Started

### Local Development
```bash
docker-compose up -d --build
```

### AWS Deployment
- **Quick Start**: See [AWS Quick Start Guide](./docs/aws-quickstart.md)
- **Detailed Guide**: See [AWS Deployment Guide](./docs/aws-deployment.md)

The platform supports both LocalStack (for local testing) and real AWS services. When deploying to AWS, the services automatically use EC2 instance role credentials - no access keys needed!
