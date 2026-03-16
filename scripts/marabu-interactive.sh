#!/bin/bash

marabu="./bin/no-bootstrap"

# Make sure we're running inside kitty
if [ -z "$KITTY_WINDOW_ID" ]; then
  echo "Error: This script should be run inside a kitty terminal. Fallback to CLI only."
  exec $marabu
  exit 1
fi

FILE="${HOME}/dev/blockchain/marabu/logs/latest.log"

> $FILE

# Launch a new kitty window for logs and capture the window ID
logsWindowID=$(kitty @ launch \
  --title "Marabu Logs" \
  --keep-focus \
  bash -c "tail -n +1 -F '$FILE' 2>/dev/null")

trap "kitty @ close-window --match id:$logsWindowID" EXIT

# Run CLI in current terminal
clear
$marabu
CLI_PID=$!
wait $CLI_PID
