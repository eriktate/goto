#!/bin/bash

echo "Sourcing jump!"

function jmp() {
	echo "Running jump..."
	output=$(jump "$@")
	echo $output
	eval $output
	# eval $(jump "$@")
}
