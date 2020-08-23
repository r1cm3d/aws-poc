resource "aws_dynamodb_table" "chargeback-table" {
  name         = "ChargebackError"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "Dispute_ID"
  range_key    = "TimeStamp"

  attribute {
    name = "Dispute_ID"
    type = "N"
  }

  attribute {
    name = "TimeStamp"
    type = "S"
  }

  tags = {
    Environment = var.account
  }
}

resource "aws_dynamodb_table" "claim-table" {
  name         = "Claim"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "Dispute_ID"
  range_key    = "TimeStamp"

  attribute {
    name = "Dispute_ID"
    type = "N"
  }

  attribute {
    name = "TimeStamp"
    type = "S"
  }

  tags = {
    Environment = var.account
  }
}