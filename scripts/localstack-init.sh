#!/bin/bash
set -e

export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
ENDPOINT=http://localhost:4566

aws --endpoint-url=$ENDPOINT kinesis create-stream \
  --stream-name usage-events \
  --shard-count 1

aws --endpoint-url=$ENDPOINT dynamodb create-table \
  --table-name usage-counters \
  --billing-mode PAY_PER_REQUEST \
  --attribute-definitions \
    AttributeName=pk,AttributeType=S \
    AttributeName=sk,AttributeType=S \
  --key-schema \
    AttributeName=pk,KeyType=HASH \
    AttributeName=sk,KeyType=RANGE
