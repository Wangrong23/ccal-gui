#!/bin/bash

go build -o ccal-gui -ldflags="-s -w" Form1.go Form1Impl.go jieqi24.go about.go dimu.go zeji.go listDay.go rules.go jinji.go other.go main.go
GOOS=windows GOARCH=amd64 go build -o ccal-gui-win.exe -ldflags="-s -w" Form1.go Form1Impl.go jieqi24.go about.go dimu.go zeji.go listDay.go rules.go jinji.go other.go main.go