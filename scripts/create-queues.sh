#!/bin/bash

. env/.aws.env

set -x

for q in $(./.list-queues.sh); do
  echo "$q"
  echo "$REGION"
  echo "$ENDPOINT"

	aws sqs create-queue \
		--queue-name "$q" \
		--region "$REGION" \
		--endpoint-url "$ENDPOINT" 2>/dev/null
done
