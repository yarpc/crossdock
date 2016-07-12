#!/bin/bash

# Make sure that Crossdock exits with a non-0
# when 1 or more tests fail

dir=$(dirname "$0")
file="$dir/../docker-compose-fail.yml"

echo "The following should FAIL - "

build=$( \
    ( docker-compose -f "$file" kill \
    && docker-compose -f "$file" rm -f \
    && docker-compose -f "$file" build crossdock \
    ) 2>&1 )

if [ $? -ne 0 ]; then
    echo "Failed to build:"
    echo "$build"
    exit 1
fi

# run takes a host name that correlates to a layer in docker-compose-fail.yml
# and a description. The first failure scenario can be ran like so:
#
#   run fail2 "invalid json"
#
function run {
    out=$( docker-compose -f "$file" run -e AXIS_CLIENT="pass,$1" crossdock )
    if [ $? -eq 0 ]; then
        failmsg="Expected non-0 exit code for '$2'"
        sep="-------------------------------------------------------------------"
        echo "$failmsg:"
        echo "$sep"
        echo "$out"
        echo "$sep"
        echo "^^^ $failmsg"
        exit 1
    fi
}

run fail0 "bad host"
run fail1 "invalid json"
run fail2 "incorrect json"
run fail3 "no results from client"
run fail4 "explicit failure"
