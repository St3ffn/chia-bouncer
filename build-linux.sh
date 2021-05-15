#!/bin/bash

SCRIPT_PATH=$(dirname "$(realpath -s "$0")")

env GOOS=linux GOARCH=amd64 go build -o chia-bouncer "$SCRIPT_PATH/main.go"
