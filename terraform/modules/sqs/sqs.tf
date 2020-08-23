resource "aws_sqs_queue" "this" {
  depends_on = [aws_sqs_queue.this_dead_letter]

  name                        = var.name
  visibility_timeout_seconds  = var.visibility_timeout_seconds
  message_retention_seconds   = var.message_retention_seconds
  max_message_size            = var.max_message_size
  delay_seconds               = var.delay_seconds
  receive_wait_time_seconds   = var.receive_wait_time_seconds
  policy                      = var.policy
  redrive_policy              = "{\"deadLetterTargetArn\":\"${aws_sqs_queue.this_dead_letter.arn}\",\"maxReceiveCount\":${var.redrive_policy_maxreceivecount}}"
  fifo_queue                  = var.fifo_queue
  content_based_deduplication = var.content_based_deduplication

  tags = var.tags
}

resource "aws_sqs_queue" "this_dead_letter" {
  name                        = "${var.name}-dead-letter"
  visibility_timeout_seconds  = var.visibility_timeout_seconds_dd
  message_retention_seconds   = var.message_retention_seconds == 1209600 ? var.message_retention_seconds : var.message_retention_seconds + 86400
  max_message_size            = var.max_message_size_dd
  delay_seconds               = var.delay_seconds_dd
  receive_wait_time_seconds   = var.receive_wait_time_seconds_dd
  policy                      = var.policy_dd
  redrive_policy              = var.redrive_policy_dd
  fifo_queue                  = var.fifo_queue_dd
  content_based_deduplication = var.content_based_deduplication_dd

  tags = var.tags
}

