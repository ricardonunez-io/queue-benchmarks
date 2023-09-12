#!/bin/bash

if [ -z "$1" ]; then
    echo "Usage: ./tests.sh <number>"
    exit 1
fi

/usr/bin/time -f "%U, %S, %P, %M, %W" -o $1-goasync.csv ./go/async/main $1
/usr/bin/time -f "%U, %S, %P, %M, %W" -o $1-gosync.csv ./go/sync/main $1
/usr/bin/time -f "%U, %S, %P, %M, %W" -o $1-jsasync.csv bun run ./javascript/async.js $1
/usr/bin/time -f "%U, %S, %P, %M, %W" -o $1-jssync.csv bun run ./javascript/sync.js $1
