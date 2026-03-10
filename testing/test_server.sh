#!/bin/bash

clients=$1

for (( i=1; i<=clients; i++ )); do 
(echo "{\"type\":\"hello\",\"agent\":\"script\",\"version\":\"0.10.0\"}"; echo "{\"type\":\"getpeers\"}"; sleep 2) | nc localhost 18018 & done
