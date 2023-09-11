#!/bin/bash

if [ -z "$1" ]; then
    echo "Usage: ./tests.sh <number>"
    exit 1
fi

hyperfine './go/async/main '"$1" -w 100
hyperfine './go/sync/main '"$1" -w 100
hyperfine 'bun run ./javascript/async.js '"$1" -w 100
hyperfine 'bun run ./javascript/sync.js '"$1" -w 100