#!/bin/bash

go build -o ccal-gui-duit-v0.6.0 -ldflags="-s -w" main.go run.go
