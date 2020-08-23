module "disputes_sqs_queue" {
  source = "./modules//sqs"
  name   = "Disputes"

  tags = {
    Name    = "disputes_sqs_queue"
    project = var.project
    env     = var.account
  }
}

module "chargeback_status_sqs_queue" {
  source = "./modules//sqs"
  name   = "ChargebackStatus"

  tags = {
    Name    = "chargeback_status_sqs_queue"
    project = var.project
    env     = var.account
  }
}

module "chargeback_sqs_queue" {
  source = "./modules//sqs"
  name   = "Chargeback"

  tags = {
    Name    = "Chargeback_sqs_queue"
    project = var.project
    env     = var.account
  }
}

module "chargeback_update_sqs_queue" {
  source = "./modules//sqs"
  name   = "ChargebackUpdate"

  tags = {
    Name    = "chargeback_update_sqs_queue"
    project = var.project
    env     = var.account
  }
}

