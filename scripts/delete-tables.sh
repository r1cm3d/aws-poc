#!/bin/bash

. env/.aws.env

for t in $(./.list-tables.sh); do
	aws dynamodb delete-table \
		--table-name "$t" \
		--region "$REGION" \
		--endpoint-url "$ENDPOINT" 2>/dev/null
done
