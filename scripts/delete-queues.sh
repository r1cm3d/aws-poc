#!/bin/bash

. env/.aws.env

for q in $(./.list-queues.sh); do
	aws sqs delete-queue \
		--queue-url "$q" \
		--region "$REGION" \
		--endpoint-url "$ENDPOINT" 2>/dev/null
done
