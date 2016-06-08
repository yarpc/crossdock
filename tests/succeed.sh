#!/bin/bash

# Make sure that Crossdock exits with a 0
# when 0 total tests fail

dir=$(dirname "$0")
file="$dir/../docker-compose.yml"

docker-compose -f "$file" kill
docker-compose -f "$file" rm -f
docker-compose -f "$file" build crossdock

echo "The following should PASS - "

docker-compose -f "$file" run crossdock
