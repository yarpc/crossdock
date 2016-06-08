#!/bin/bash

# Make sure that Crossdock exits with a non-0
# when 1 or more tests fail

dir=$(dirname "$0")
file="$dir/../docker-compose-fail.yml"

_=$(docker-compose -f "$file" run crossdock)

if [ $? == 0 ]; then
    echo "Expected non-0 exit code"
    exit 1
fi
