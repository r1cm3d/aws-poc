resource "aws_dynamodb_table" "dispute" {
  name         = "Dispute"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "ID"

  attribute {
    name = "ID"
    type = "S"
  }

  tags = {
    Environment = var.account
  }
}