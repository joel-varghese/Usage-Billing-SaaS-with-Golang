resource "aws_kinesis_stream" "usage_events" {
  name             = "usage-events"
  shard_count      = 2
  retention_period = 24
}
