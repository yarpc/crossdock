#!/bin/bash

# Make sure that Crossdock exits with a 0
# when 0 total tests fail

dir=$(dirname "$0")
file="$dir/../docker-compose.yml"

echo "The following should PASS - "

out=$( \
	( docker-compose -f "$file" kill \
	&& docker-compose -f "$file" rm -f \
	&& docker-compose -f "$file" build crossdock \
	&& docker-compose -f "$file" run crossdock \
	) 2>&1)

if [[ $? -ne 0 ]]; then
    echo "Expected 0 exit code"
	echo "Output:"
	echo "$out"
    exit 1
fi
