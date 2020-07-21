#!/bin/bash
function remove-label() {
	echo "$1" | sed 's/^.*=//g' | tr -d "'"
}

while read -r v; do
	if [[ "$v" =~ .*SQS.* ]]; then
		queues="${queues}$(remove-label "$v") "
	fi
done <env/.env

echo "$queues"
