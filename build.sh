#!/bin/bash

go build -o ccal-gui-linux64 -ldflags="-s -w" Form1.go Form1Impl.go jieqi24.go about.go dimu.go zeji.go listDay.go rules.go jinji.go other.go yanggongji.go main.go
GOOS=windows GOARCH=amd64 go build -o ccal-gui-win64.exe -ldflags="-s -w" Form1.go Form1Impl.go jieqi24.go about.go dimu.go zeji.go listDay.go rules.go jinji.go other.go yanggongji.go main.go
