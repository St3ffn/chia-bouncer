#!/bin/bash

SCRIPT_PATH=$(dirname "$(realpath -s "$0")")

go build -o chia-bouncer "$SCRIPT_PATH/main.go"
env GOOS=linux GOARCH=amd64 go build -o chia-bouncer-linux-amd64 "$SCRIPT_PATH/main.go"
env GOOS=darwin GOARCH=amd64 go build -o chia-bouncer-darwin-amd64 "$SCRIPT_PATH/main.go"
env GOOS=windows GOARCH=amd64 go build -o chia-bouncer-windows-amd64 "$SCRIPT_PATH/main.go"
env GOOS=linux GOARCH=arm64 go build -o chia-bouncer-linux-arm64 "$SCRIPT_PATH/main.go"
