#!/bin/bash
while read -r v; do
	echo "Loading: $v"
	export v
done <env/.env
