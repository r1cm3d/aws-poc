#!/bin/bash

. env/.aws.env

aws s3api create-bucket \
	--bucket "$BUCKET" \
	--region "$REGION" \
	--endpoint-url "$ENDPOINT" 2>/dev/null
