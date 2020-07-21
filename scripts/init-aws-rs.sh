#!/bin/bash

. env/.aws.env

function run() {
	for f in "$1"*; do
		echo "executing $f"
		eval "sh $f"
	done
}

run "delete"
run "create"
