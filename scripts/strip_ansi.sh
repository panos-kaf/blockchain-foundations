#!/bin/bash

if [ $# -ne 1 ]; then
    echo "Usage: $0 <file>"
    exit 1
fi

input="$1"

# Remove ANSI escape sequences
sed -E 's/\x1B\[[0-9;]*[A-Za-z]//g' "$input"
