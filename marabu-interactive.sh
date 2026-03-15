#!/bin/bash

marabu="./marabu"

# Make sure we're running inside kitty
if [ -z "$KITTY_WINDOW_ID" ]; then
  echo "Error: This script should be run inside a kitty terminal. Fallback to CLI only."
  $marabu
  exit 1
fi

FILE="${HOME}/dev/blockchain/marabu/logs/latest.log"

> $FILE

# get id of the current window
cliWindowID=$KITTY_WINDOW_ID

# Launch a new kitty window for logs and capture the window ID
logsWindowID=$(kitty @ launch --title "Marabu Logs" bash -c "tail -n +1 -F ${FILE} 2>/dev/null")

kitty @ focus-window --match id:$cliWindowID

# Run CLI in current terminal
$marabu
CLI_PID=$!
wait $CLI_PID

# When CLI exits, close the log window
kitty @ close-window --match id:$logsWindowID
