resource "aws_s3_bucket" "raw_events" {
  bucket = "${var.project_name}-raw-events"
}
