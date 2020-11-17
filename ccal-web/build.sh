#!/bin/bash

go build -o ccal-web-dev -ldflags="-s -w" main.go
GOOS=windows GOARCH=amd64 go build -o ccal-web-dev.exe -ldflags="-s -w" main.go
