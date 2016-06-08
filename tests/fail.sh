#!/bin/bash

# Make sure that Crossdock exits with a non-0
# when 1 or more tests fail

dir=$(dirname "$0")
file="$dir/../docker-compose-fail.yml"

docker-compose -f "$file" kill
docker-compose -f "$file" rm -f
docker-compose -f "$file" build crossdock

echo "The following should FAIL - "

docker-compose -f "$file" run crossdock

if [ $? == 0 ]; then
    echo "Expected non-0 exit code"
    exit 1
fi
