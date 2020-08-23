provider "aws" {
  region                      = var.region
  access_key                  = var.access_key
  s3_force_path_style         = var.s3_force_path_style
  secret_key                  = var.secret_key
  skip_credentials_validation = var.skip_credentials_validation
  skip_metadata_api_check     = var.skip_metadata_api_check
  skip_requesting_account_id  = var.skip_requesting_account_id

  endpoints {
    iam      = var.localstack_url
    s3       = var.localstack_url
    dynamodb = var.localstack_url
    sqs      = var.localstack_url
  }
}

variable "project" {
  default = "aws-poc"
}

resource "aws_iam_user" "aws-poc" {
  name = "aws-poc"
  tags = {
    project = var.project
    env     = var.account
  }
}

resource "aws_iam_policy" "aws_poc_sqs_policy" {
  name = "aws_poc_sqs_policy"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "sqs",
            "Effect": "Allow",
            "Action": [
                "sqs:GetQueueUrl",
                "sqs:ChangeMessageVisibility",
                "sqs:SendMessageBatch",
                "sqs:ReceiveMessage",
                "sqs:SendMessage",
                "sqs:GetQueueAttributes",
                "sqs:ListQueueTags",
                "sqs:ListDeadLetterSourceQueues",
                "sqs:ChangeMessageVisibilityBatch",
                "sqs:DeleteMessage"
            ],
            "Resource": [
                "${module.disputes_sqs_queue.this_sqs_queue_arn}",
                "${module.chargeback_status_sqs_queue.this_sqs_queue_arn}",
                "${module.chargeback_update_sqs_queue.this_sqs_queue_arn}",
                "${module.chargeback_sqs_queue.this_sqs_queue_arn}"
            ]
        }
    ]
}
EOF
}
