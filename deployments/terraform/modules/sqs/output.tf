output "this_sqs_queue_id" {
  description = "The URL for the created Amazon SQS queue"
  value       = aws_sqs_queue.this.id
}

output "this_sqs_queue_arn" {
  description = "The ARN of the SQS queue"
  value       = aws_sqs_queue.this.arn
}


output "this_dead_letter_sqs_queue_id" {
  description = "The URL for the created Amazon SQS queue"
  value       = aws_sqs_queue.this_dead_letter.id
}

output "this_dead_letter_sqs_queue_arn" {
  description = "The ARN of the SQS queue"
  value       = aws_sqs_queue.this_dead_letter.arn
}
