#!/bin/bash

go build -o ccal-web -ldflags="-s -w" main.go
