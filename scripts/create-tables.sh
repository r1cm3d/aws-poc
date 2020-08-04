#!/bin/bash

. env/.aws.env

set -x

for t in json/*_table.json; do
  echo "$REGION"
  echo "file://$t"
  echo "ENDPOINT"
	aws dynamodb create-table \
		--region "$REGION" \
		--cli-input-json "file://$t" \
		--endpoint-url "$ENDPOINT" 2>/dev/null
done
