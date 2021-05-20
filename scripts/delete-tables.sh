#!/bin/bash

. env/.aws.env

for t in $(./.list-tables.sh); do
	aws dynamodb delete-table \
		--table-name "$t" \
		--region "$REGION" \
		--endpoint-url "$ENDPOINT" >/dev/null 2>&1
done
