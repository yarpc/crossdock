#!/bin/bash

# Make sure that Crossdock exits with a non-0
# when 1 or more tests fail

dir=$(dirname "$0")
file="$dir/../docker-compose-fail.yml"

echo "The following should FAIL - "

out=$( \
	( docker-compose -f "$file" kill \
	&& docker-compose -f "$file" rm -f \
	&& docker-compose -f "$file" build crossdock \
	&& docker-compose -f "$file" run crossdock \
	) 2>&1)

if [ $? -eq 0 ]; then
    echo "Expected non-0 exit code"
	echo "Output:"
	echo "$out"
    exit 1
fi
